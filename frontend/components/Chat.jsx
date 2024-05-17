import React, { useState, useEffect, useRef } from 'react';
import './Chat.css';
import { useRouter } from 'next/navigation';

const Chat = ({ chatID, role, login }) => {
    const [messages, setMessages] = useState([]);
    const [inputMessage, setInputMessage] = useState('');
    const ws = useRef(null);
    const router = useRouter();

    useEffect(() => {
        ws.current = new WebSocket('ws://localhost:8080/ws');

        ws.current.onopen = () => {
            console.log('Connected to server');
            // Adjusted to handle both roles more dynamically
            if (!login) ws.current.send(JSON.stringify({ action: role === 'client' ? 'create_chat' : 'join_chat', chat_id: chatID }));
            else ws.current.send(JSON.stringify({ action: 'create_chat', chat_id: chatID, login: login }));
        };

        ws.current.onmessage = (event) => {
            const messageData = JSON.parse(event.data);
            if (messageData.action === 'chat_history') {
                setMessages((prevMessages) => [...prevMessages, ...messageData.history.map(msg => ({ message: msg.message, role: msg.role }))]);
            } else if (messageData.message) {
                setMessages((prevMessages) => [...prevMessages, { message: messageData.message, role: messageData.role }]);
            }
        };

        ws.current.onclose = () => {
            console.log('Disconnected from server');
        };

        ws.current.onerror = (error) => {
            console.error(`WebSocket error: ${error.message}`);
        };

        return () => {
            ws.current.close();

        };
    }, [chatID, role]);

    const sendMessage = () => {
        if (inputMessage.trim() === '') return;

        ws.current.send(JSON.stringify({ action: 'send_message', chat_id: chatID, message: inputMessage, role: role }));
        setInputMessage('');
    };

    const closeChat = () => {

        ws.current.send(JSON.stringify({ action: 'close_chat', chat_id: chatID/*, user: userID*/ }));
        router.push('/admin');
    }

    return (
        <div className="parent-container">
            {role !== 'client' && (
                <button className = "close-chat"onClick={closeChat}>Close Chat</button>
            )}
            <div className="chat-container">
                <div className="messages-container">
                    <ul>
                        {messages.map((msg, index) => (
                            <li className={`message ${msg.role === "admin" ? "admin_message" : ""}`} key={index}>
                                <strong>{msg.role.toUpperCase()}: </strong>{msg.message}
                            </li>
                        ))}
                    </ul>
                </div>
                <div className="input-container">
                    <input
                        type="text"
                        value={inputMessage}
                        onChange={(e) => setInputMessage(e.target.value)}
                        onKeyPress={(e) => {
                            if (e.key === 'Enter') {
                                sendMessage();
                            }
                        }}
                    />
                    <button onClick={sendMessage}>Send</button>
                </div>
            </div>
        </div>
    );
};

export default Chat;