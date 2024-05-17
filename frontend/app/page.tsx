"use client";
import React, { useState } from 'react';
import {useRouter} from "next/navigation";
import './page.css'
const Home = () => {
    const router = useRouter();

    const showAdminPanel = () => {
        router.push('/admin');
    };

    return (
        <div className="home-container">
            <a href={"/client"}>Client Panel</a>
            <button onClick={showAdminPanel}>Admin Panel</button>
        </div>
    );
};

export default Home;
