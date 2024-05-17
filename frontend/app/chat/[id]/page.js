"use client";
import React, { useEffect } from 'react';
import {useRouter, useSearchParams} from 'next/navigation';
import Chat from '../../../components/Chat';

const ChatPage = ({params}) => {
    const router = useRouter();
    const searchParams = useSearchParams()
    const id = params.id;
    const role = searchParams.get("role");
    const login = searchParams.get("login");

    useEffect(() => {
        if (!id || !role) {
            // If ID or role is missing, redirect to the home page
            router.push('/');
        }
    }, [id, role, router]);

    if (!id || !role) {
        return null; // Avoid rendering if id or role is missing
    }

    return (
        <div>
            <Chat chatID={id} role={role} login={login}/>
        </div>
    );
};

export default ChatPage;
