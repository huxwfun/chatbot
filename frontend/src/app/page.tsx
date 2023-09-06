"use client"

import React, { useState, useEffect } from 'react';
import Box from '@mui/material/Box';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import Messages from '@/components/Messages';
import { Data, USER_ANONYMOUS } from '@/components/Messages'
import ListItemButton from '@mui/material/ListItemButton';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import Avatar from '@mui/material/Avatar';
import ListItemText from '@mui/material/ListItemText';
import Link from 'next/link';
import Divider from '@mui/material/Divider';
import Modal from '@mui/material/Modal';
import Typography from '@mui/material/Typography';

export default function HomePage() {
  const [data, setData] = useState<Data>({ users: [], me: USER_ANONYMOUS, messages: [], chats: [] })
  const [isLoading, setLoading] = useState(true)
  const [activeChatId, setActiveChatId] = useState<string | null>(null)
  const [error, setError] = useState(false)
  const [logs, setLogs] = useState('')
  useEffect(() => {
    fetch('http://localhost:8080/data')
      .then((res) => res.json())
      .then((d) => {
        const chats = d.chats || []
        setActiveChatId(chats[0]?.id || null)
        setData({
          users: d.users,
          me: d.me,
          messages: d.messages || [],
          chats: chats
        })
        setLoading(false)
      }).catch(err => {
        setError(true)
        setLoading(false)
      })
  }, [])
  console.log(logs)
  if (error) {
    return (<Box sx={{ display: 'flex' }}>
      <Box sx={{ width: '100%', textAlign: 'center', paddingTop: 10 }}>
        Fetch data from server failed. Check server status pls, then refresh page to continue.
      </Box>
    </Box>)
  }
  return (
    <Box sx={{ display: 'flex' }}>
      <Box sx={{ width: '100%' }}>
        {isLoading ? 'loading...' :
          (< Messages chatId={activeChatId || ''} data={data} />)
        }
      </Box>
      <Drawer
        sx={{
          width: 320,
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: 320,
            boxSizing: 'border-box',
            top: ['48px', '56px', '64px'],
            height: 'auto',
            bottom: 0,
          },
        }}
        variant="permanent"
        anchor="right"
      >
        {!isLoading && (
          <List>
            <ListItem disablePadding>
              <ListItemButton component={Link} href='/'>
                <ListItemAvatar>
                  <Avatar>
                    <img src={data.me?.avatar} height={40} width={40} />
                  </Avatar>
                </ListItemAvatar>
                <ListItemText primary={data.me.name} />
              </ListItemButton>
            </ListItem>
            <ListItem>
              <ListItemText primary={'recent chats:'} />
            </ListItem>
            {
              data.chats.map(({ id, botId }) => {
                const bot = data.users.find(u => u.id == botId)
                if (bot)
                  return (
                    <ListItem key={id} disablePadding>
                      <ListItemButton disabled={id == activeChatId} component={Link} href='/' onClick={() => setActiveChatId(id)}>
                        <ListItemAvatar>
                          <Avatar style={{ width: 30, height: 30 }}>
                            <img src={bot.avatar} height={30} width={30} />
                          </Avatar>
                        </ListItemAvatar>
                        <ListItemAvatar >
                          <Avatar style={{ width: 30, height: 30 }}>
                            <img src={data.me.avatar} height={30} width={30} />
                          </Avatar>
                        </ListItemAvatar>
                      </ListItemButton>
                    </ListItem>
                  )
              })
            }
            <Divider />
            <ListItem disablePadding>
              <ListItemButton component={Link} href='/' onClick={() => {
                fetch(`http://localhost:8080/review?authentication=${data.me.id}`, { method: 'POST' })
              }}>
                <ListItemText primary={'Trigger sample review flow!'} />
              </ListItemButton>
            </ListItem>
            <ListItem disablePadding>
              <ListItemButton component={Link} href='/' onClick={async () => {
                const output: any = await fetch(`http://localhost:8080/instruction?authentication=${data.me.id}`).then(res => res.text())
                setLogs(output)
              }}>
                <ListItemText primary={'See instruction!'} />
              </ListItemButton>
            </ListItem>
            <ListItem disablePadding>
              <ListItemButton component={Link} href='/' onClick={async () => {
                const output: any = await fetch(`http://localhost:8080/logs?authentication=${data.me.id}`).then(res => res.text())
                setLogs(output)
              }}>
                <ListItemText primary={'See all domain event logs'} />
              </ListItemButton>
            </ListItem>
          </List>
        )}
      </Drawer>
      <Modal
        open={logs.length > 0}
        onClose={() => setLogs('')}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
      >
        <Box sx={{
          position: 'absolute' as 'absolute',
          top: '50%',
          left: '50%',
          width: '80%',
          transform: 'translate(-50%, -50%)',
          bgcolor: 'background.paper',
          border: '2px solid #000',
          boxShadow: 24,
          p: 4,
          maxHeight: '60%',
          overflowY: 'scroll',
        }}>
          {logs.split('\n').map((l, i) => {
            if (l.indexOf('---') > -1) {
              return <Typography key={i} variant="h6" gutterBottom>{l}</Typography>
            } else {
              return (<Typography key={i} gutterBottom={l.startsWith('from state')}>{l ? l : ' '}</Typography>)
            }
          })}
        </Box>
      </Modal>
    </Box>
  );
}
