package main

import (
	"bufio"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type ChatRoom struct {
	ID      string
	Clients map[*websocket.Conn]bool
	Users   []string // List of usernames in the chat
	Admin   string   // Username of the chat room admin
	Mutex   sync.Mutex
}

type ChatManager struct {
	Rooms map[string]*ChatRoom
	Mutex sync.Mutex
}

var manager = ChatManager{
	Rooms: make(map[string]*ChatRoom),
}

var userChats = make(map[string][]string)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		var msg map[string]string
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Read:", err)
			break
		}

		processMessage(ws, msg)
	}
}

func processMessage(ws *websocket.Conn, msg map[string]string) {
	action := msg["action"]

	logrus.Info("Processing action: ", action)
	switch action {
	case "create_chat":
		createOrJoinChat(ws, msg)
	case "list_chats":
		listChats(ws)
	case "list_user_chats":
		logrus.Info(userChats, msg)
		listUserChats(ws, msg)
	case "join_chat":
		joinChat(ws, msg)
	case "send_message":
		sendMessage(ws, msg)
	case "close_chat":
		closeChat(ws, msg)
	}
}

func listUserChats(ws *websocket.Conn, msg map[string]string) {
	login := msg["login"]
	manager.Mutex.Lock()
	ws.WriteJSON(map[string]interface{}{
		"action": "list_user_chats",
		"chats":  userChats[login],
	})
	manager.Mutex.Unlock()
}

func createOrJoinChat(ws *websocket.Conn, msg map[string]string) {
	chatID := msg["chat_id"]
	login := msg["login"]

	manager.Mutex.Lock()
	room, exists := manager.Rooms[chatID]
	if !exists {
		room = &ChatRoom{
			ID:      chatID,
			Clients: make(map[*websocket.Conn]bool),
		}
		manager.Rooms[chatID] = room
		userChats[login] = append(userChats[login], chatID)
		notifyAdmins("new_chat", chatID)
	}
	manager.Mutex.Unlock()

	room.Mutex.Lock()
	room.Clients[ws] = true
	room.Mutex.Unlock()

	chatHistory, err := getChatHistory(chatID)
	if err != nil {
		log.Println("Error retrieving chat history:", err)
		return
	}

	room.Mutex.Lock()
	for client := range room.Clients {
		if err := client.WriteJSON(map[string]interface{}{"action": "chat_history", "history": chatHistory}); err != nil {
			log.Println("Error sending chat history to client:", err)
			client.Close()
			delete(room.Clients, client)
		}
	}
	room.Mutex.Unlock()
}

func listChats(ws *websocket.Conn) {
	manager.Mutex.Lock()
	chatIDs := []string{}
	for chatID := range manager.Rooms {
		chatIDs = append(chatIDs, chatID)
	}
	manager.Mutex.Unlock()
	ws.WriteJSON(map[string]interface{}{
		"action": "list_chats",
		"chats":  chatIDs,
	})
}

func joinChat(ws *websocket.Conn, msg map[string]string) {
	chatID := msg["chat_id"]
	manager.Mutex.Lock()
	if room, exists := manager.Rooms[chatID]; exists {
		room.Mutex.Lock()
		room.Clients[ws] = true

		chatHistory, err := getChatHistory(chatID)
		if err != nil {
			log.Println("Error retrieving chat history:", err)
		} else {
			if err := ws.WriteJSON(map[string]interface{}{"action": "chat_history", "history": chatHistory}); err != nil {
				log.Println("Error sending chat history:", err)
			}
		}

		room.Mutex.Unlock()
	}
	manager.Mutex.Unlock()
}

func getChatHistory(chatID string) ([]map[string]string, error) {
	var history []map[string]string
	file, err := os.Open(chatID + ".txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var msg map[string]string
		if err := json.Unmarshal([]byte(scanner.Text()), &msg); err != nil {
			log.Println("Error decoding message:", err)
			continue
		}
		history = append(history, msg)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

func sendMessage(ws *websocket.Conn, msg map[string]string) {
	chatID := msg["chat_id"]
	message := msg["message"]
	role := msg["role"]

	manager.Mutex.Lock()
	if room, exists := manager.Rooms[chatID]; exists {
		room.Mutex.Lock()
		for client := range room.Clients {
			if err := client.WriteJSON(map[string]string{"message": message, "role": role}); err != nil {
				log.Println("Write:", err)
				client.Close()
				delete(room.Clients, client)
			}
		}
		room.Mutex.Unlock()

		saveChatData(chatID, msg)
	}
	manager.Mutex.Unlock()
}

func closeChat(ws *websocket.Conn, msg map[string]string) {
	chatID := msg["chat_id"]

	// find in userChats corresponding chatID and delete it from userChats login is not defined.
	removeFromUserChats(chatID)

	manager.Mutex.Lock()
	if room, exists := manager.Rooms[chatID]; exists {
		room.Mutex.Lock()
		for client := range room.Clients {
			client.Close()
		}

		delete(manager.Rooms, chatID)
		room.Mutex.Unlock()
	}
	manager.Mutex.Unlock()
}

func removeFromUserChats(chatID string) {
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()

	for login, chats := range userChats {
		for i, id := range chats {
			if id == chatID {
				userChats[login] = append(chats[:i], chats[i+1:]...)
				break
			}
		}
	}
}

func notifyAdmins(event string, chatID string) {
	logrus.Info("Notifying admins about event: ", event, " in chat: ", chatID)
}

func saveChatData(chatID string, data map[string]string) {
	file, err := os.OpenFile(chatID+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		log.Println("Error writing to file:", err)
	}
}
