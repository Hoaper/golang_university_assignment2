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
	Mutex   sync.Mutex
}

type ChatManager struct {
	Rooms map[string]*ChatRoom
	Mutex sync.Mutex
}

var manager = ChatManager{
	Rooms: make(map[string]*ChatRoom),
}

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
	case "join_chat":
		joinChat(ws, msg)
	case "send_message":
		sendMessage(ws, msg)
	case "close_chat":
		closeChat(ws, msg)
	}
}

func createOrJoinChat(ws *websocket.Conn, msg map[string]string) {
	chatID := msg["chat_id"]
	manager.Mutex.Lock()
	room, exists := manager.Rooms[chatID]
	if !exists {
		// Create chat if it doesn't exist
		room = &ChatRoom{
			ID:      chatID,
			Clients: make(map[*websocket.Conn]bool),
		}
		manager.Rooms[chatID] = room
		notifyAdmins("new_chat", chatID)
	}
	manager.Mutex.Unlock()

	room.Mutex.Lock()
	room.Clients[ws] = true
	room.Mutex.Unlock()

	// Retrieve chat history
	chatHistory, err := getChatHistory(chatID)
	if err != nil {
		log.Println("Error retrieving chat history:", err)
		return
	}

	// Send chat history to all clients in the room
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

		// Retrieve chat history
		chatHistory, err := getChatHistory(chatID)
		if err != nil {
			log.Println("Error retrieving chat history:", err)
		} else {
			// Send chat history to client
			if err := ws.WriteJSON(map[string]interface{}{"action": "chat_history", "history": chatHistory}); err != nil {
				log.Println("Error sending chat history:", err)
			}
		}

		room.Mutex.Unlock()
	}
	manager.Mutex.Unlock()
}

// Assume this function retrieves chat history for a given chatID
func getChatHistory(chatID string) ([]map[string]string, error) {
	var history []map[string]string
	// Logic to read chat history from a file or database
	// For example, reading from a file (simplified):
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
	manager.Mutex.Lock()
	if room, exists := manager.Rooms[chatID]; exists {
		room.Mutex.Lock()
		for client := range room.Clients {
			if err := client.WriteJSON(map[string]string{"message": message}); err != nil {
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
