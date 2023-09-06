"use client"
import React, { useEffect, useState } from "react";
import { Button, Paper, Snackbar } from "@mui/material";
import { TextInput } from "./TextInput";
import { MessageLeft, MessageRight } from "./Message";
import { v4 as uuid } from 'uuid'

const styles = {
    container: {
        height: "calc(100vh - 64px)",
        width: "100%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        flexDirection: "column",
        position: "relative",
    },
    messagesBody: {
        width: "100%",
        margin: "0 0 12px 0",
        height: "calc(100vh - 124px)",
        overflow: "scroll",
    }
}

type User = {
    id: string
    name: string,
    avatar: string,
}

export const USER_ANONYMOUS = {
    id: '',
    name: 'anonymous',
    avatar: ''
}

type Message = {
    id: string
    chatId: string,
    body: string,
    from: string,
    timeCreated: string,
}

type DisplayMessage = {
    id: string
    body: string,
    name: string,
    avatar: string,
    timeCreated: string,
}

type Chat = {
    id: string,
    botId: string,
    customerId: string
}
function toDisplayMessage(m: Message, users: User[]): DisplayMessage {
    const user = users.find(u => u.id === m.from) || USER_ANONYMOUS
    return {
        id: m.id,
        body: m.body,
        name: user.name,
        avatar: user.avatar,
        timeCreated: m.timeCreated,
    }
}

export type Data = {
    users: User[],
    me: User,
    messages: Message[],
    chats: Chat[]
}
var conn: WebSocket
export default function App({ data: { me, users, messages: original }, chatId }: { data: Data, chatId: string }) {
    const [messages, setMessages] = useState(original.filter(m => m.chatId == chatId).map(m => toDisplayMessage(m, users)))
    const [disconnected, setDisconnected] = useState<boolean>(false)

    useEffect(() => {
        conn = new WebSocket(`ws://localhost:8080/ws?authentication=${me.id}`)
        // const conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onopen = e => {
            setDisconnected(false)
        };
        conn.onclose = e => {
            setDisconnected(true)
        };
        conn.onerror = e => {
            console.error(e)
        }
        conn.onmessage = function (evt) {
            const msg = JSON.parse(evt.data) as Message
            let message = toDisplayMessage(msg, users)
            setMessages(m => [...m, message])
        }
        if (conn.readyState == conn.CLOSED) {
            setDisconnected(true)
        }
        return () => {
            console.log('ws closed')
            conn.close()
        }
    }, [me])
    return (
        <div id="111" style={styles.container as any}>
            <Paper id="style-1" sx={styles.messagesBody}>
                {
                    messages.map(m => (
                        m.name == me.name ?
                            <MessageRight
                                key={m.id}
                                message={m.body}
                                timestamp=""
                                photoURL={m.avatar}
                                displayName=""
                                avatarDisp={true}
                            /> :
                            <MessageLeft
                                key={m.id}
                                message={m.body}
                                timestamp=""
                                photoURL={m.avatar}
                                displayName=""
                                avatarDisp={true}
                            />
                    ))
                }
            </Paper>
            <TextInput disabled={disconnected} onSubmit={(text: string) => {
                const message: Omit<Message, 'timeCreated'> = {
                    id: uuid(),
                    chatId: chatId,
                    body: text.trim(),
                    from: me.id,
                }
                conn.send(JSON.stringify(message))
            }} />
            <Snackbar
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
                open={disconnected && window != null}
                onClose={() => window.location.reload()}
                message="Disconnected, refresh page to retry"
                action={
                    <>
                        <Button color="secondary" size="small" onClick={() => window.location.reload()}>
                            REFRESH
                        </Button>
                    </>
                }
            />
        </div>
    );
}