import React, { useState, useEffect, useRef } from 'react';
import './Chat.css';

const Chat = ({ chatID, role }) => {
    const [messages, setMessages] = useState([]);
    const [inputMessage, setInputMessage] = useState('');
    const ws = useRef(null);

    useEffect(() => {
        ws.current = new WebSocket('ws://localhost:8080/ws');

        ws.current.onopen = () => {
            console.log('Connected to server');
            // Adjusted to handle both roles more dynamically
            ws.current.send(JSON.stringify({ action: role === 'client' ? 'create_chat' : 'join_chat', chat_id: chatID }));
        };

        ws.current.onmessage = (event) => {
            const message = JSON.parse(event.data);
            if (message.action === 'chat_history') {
                // Assuming 'history' is an array of message objects
                setMessages((prevMessages) => [...prevMessages, ...message.history.map(msg => msg.message)]);
            } else if (message.message) {
                // Handling new messages
                setMessages((prevMessages) => [...prevMessages, message.message]);
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

        ws.current.send(JSON.stringify({ action: 'send_message', chat_id: chatID, message: inputMessage }));
        setInputMessage('');
    };

    return (
    <div className="parent-container">
       <div className="chat-container">
            <div className="messages-container">
                <ul>
                    {messages.map((msg, index) => (
                        <li className="message" key={index}>{msg}</li>
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