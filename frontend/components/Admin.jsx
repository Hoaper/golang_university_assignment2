import React, { useState, useEffect, useRef } from 'react';
import { useRouter } from 'next/navigation';

const Admin = () => {
    const [openChats, setOpenChats] = useState([]);
    const ws = useRef(null);
    const router = useRouter();

    useEffect(() => {
        ws.current = new WebSocket('ws://localhost:8080/ws');

        ws.current.onopen = () => {
            console.log('Connected to server');
            listChats();
        };

        ws.current.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.action === 'list_chats') {
                setOpenChats(data.chats);
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
    }, []);

    const listChats = () => {
        const msg = {
            action: 'list_chats',
        };
        ws.current.send(JSON.stringify(msg));
    };

    const joinChat = (chatID) => {
        router.push(`/chat/${chatID}?role=admin`);
    };

    return (
        <div>
            <h1>Admin Panel</h1>
            <ul>
                {openChats.map((chatID) => (
                    <li key={chatID}>
                        {chatID}
                        <button onClick={() => joinChat(chatID)}>Join Chat</button>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default Admin;
