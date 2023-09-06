import React from 'react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import DashboardIcon from '@mui/icons-material/Dashboard';
import ThemeRegistry from '@/components/ThemeRegistry/ThemeRegistry';

export const metadata = {
  title: 'Next.js App Router + Material UI v5',
  description: 'Next.js App Router + Material UI v5',
};

const DRAWER_WIDTH = 240;
const USER_ANONYMOUS = {
  id: '',
  name: 'anonymous',
  avatar: ''
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body style={{ overflow: "hidden" }}>
        <ThemeRegistry>
          <AppBar position="fixed" sx={{ zIndex: 2000 }}>
            <Toolbar sx={{ backgroundColor: 'background.paper' }}>
              <DashboardIcon sx={{ color: '#444', mr: 2, transform: 'translateY(-2px)' }} />
              <Typography variant="h6" noWrap component="div" color="black">
                Chatbot example
              </Typography>
            </Toolbar>
          </AppBar>
          {/* <Drawer
            sx={{
              width: DRAWER_WIDTH,
              flexShrink: 0,
              '& .MuiDrawer-paper': {
                width: DRAWER_WIDTH,
                boxSizing: 'border-box',
                top: ['48px', '56px', '64px'],
                height: 'auto',
                bottom: 0,
              },
            }}
            variant="permanent"
            anchor="left"
          >
            <Fab variant="extended" sx={{ mt: 5, ml: 1, mr: 1 }}>
              <NavigationIcon color="primary" sx={{ mr: 1 }} />
              Show instructions!
            </Fab>
            <Fab variant="extended"
             sx={{ mt: 5, ml: 1, mr: 1 }}
             onClick={() => {
              // fetch(`http://localhost:8080/review?authentication=${data.me.id}`, { method: 'POST' })
            }}>
              <NavigationIcon color="secondary" sx={{ mr: 1 }} />
              Start review!
            </Fab>
            <Fab variant="extended" sx={{ mt: 5, ml: 1, mr: 1 }}>
              <NavigationIcon sx={{ mr: 1 }} />
              Show event logs
            </Fab>
            <Divider sx={{ mt: 'auto' }} />
          </Drawer> */}
          <Box
            component="main"
            sx={{
              flexGrow: 1,
              bgcolor: 'background.default',
              mt: ['48px', '56px', '64px'],
              pl: 3,
              pr: 3,
            }}
          >
            {children}
          </Box>
        </ThemeRegistry>
      </body>
    </html>
  );
}
