"use client"
import React, {useState, useEffect, useRef} from 'react';
import { v4 as uuidv4 } from 'uuid';
import {useRouter} from "next/navigation";

const ClientPage = () => {
    const [login, setLogin] = useState('');
    const [chats, setChats] = useState([]);
    const [newChatName, setNewChatName] = useState('');
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const ws = useRef(null);
    const router = useRouter();


    useEffect(() => {
        ws.current = new WebSocket('ws://localhost:8080/ws');
        ws.current.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.action === 'list_user_chats') {
                setChats(data.chats ? data.chats : []);
                setIsLoggedIn(true);
            }
        };
        return () => ws.current.close();
    }, []);

    const handleLoginSubmit = (e) => {
        e.preventDefault();
        ws.current.send(JSON.stringify({ action: 'list_user_chats', login: login }));
    };

    const handleCreateChat = (e) => {
        e.preventDefault();
        router.push(`/chat/${uuidv4()}?role=client&login=${login}`)
        setNewChatName('');
    };

    if (!isLoggedIn) {
        return (
            <form onSubmit={handleLoginSubmit}>
                <input
                    type="text"
                    value={login}
                    onChange={(e) => setLogin(e.target.value)}
                    placeholder="Enter your login"
                    required
                />
                <button type="submit">Login</button>
            </form>
        );
    }

    return (
        <div>
            {chats.length === 0 && <h1>No chats available</h1>}
            {chats.length > 0 && <h1>Chats:</h1>}
            <ul>
                {chats.map((chat, index) => (
                    <li key={index}>
                        <a href={`/chat/${chat}?role=client`}>{chat}</a>
                    </li>
                ))}
            </ul>
            <form onSubmit={handleCreateChat}>
                <button type="submit">Create Chat</button>
            </form>
        </div>
    );
};

export default ClientPage;