"use client";
import React, { useState } from 'react';
import { v4 as uuidv4 } from 'uuid';
import {useRouter} from "next/navigation";
import './page.css'
const Home = () => {
    const router = useRouter();

    const showAdminPanel = () => {
        router.push('/admin');
    };

    return (
        <div className="home-container">
            <a href={`/chat/${uuidv4()}?role=client`}>Start Chat as Client</a>
            <button onClick={showAdminPanel}>Admin Panel</button>
        </div>
    );
};

export default Home;
