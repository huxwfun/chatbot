import React, { useState } from 'react'
import TextField from '@mui/material/TextField';
import SendIcon from '@mui/icons-material/Send';
import Button from '@mui/material/Button';


export const TextInput = ({ onSubmit, disabled }: any) => {
    const [text, setText] = useState('')
    return (
        <>
            <form noValidate autoComplete="off" style={{
                display: "flex",
                justifyContent: "center",
                width: "95%",
                margin: `0 auto`
            }} onSubmit={(e: any) => {
                e.preventDefault()
                setText('')
                onSubmit(text)
            }}>
                <TextField
                    disabled={disabled}
                    id="standard-text"
                    label="input message"
                    variant="standard"
                    fullWidth
                    value={text}
                    onChange={(e: any) => {
                        setText(e.target.value)
                    }}
                />
                <Button variant="contained" color="primary" disabled={disabled} onClick={(e) => {
                    e.preventDefault()
                    setText('')
                    onSubmit(text)
                }}>
                    <SendIcon />
                </Button>
            </form>
        </>
    )
}