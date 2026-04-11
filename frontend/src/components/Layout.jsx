// Основной layout: AppBar, боковая панель, уведомления, меню пользователя
import React, { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router';
import {
  Box,
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Menu,
  MenuItem,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import TaskAltIcon from '@mui/icons-material/TaskAlt';
import FolderIcon from '@mui/icons-material/Folder';
import DashboardIcon from '@mui/icons-material/Dashboard';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import DarkModeIcon from '@mui/icons-material/DarkMode';
import LightModeIcon from '@mui/icons-material/LightMode';
import { useDispatch } from 'react-redux';
import { logout } from '../store/authSlice';
import { addNotification } from '../store/notificationSlice';
import wsService from '../services/websocket';
import NotificationPanel from './NotificationPanel';

const DRAWER_WIDTH = 260;

export default function Layout({ darkMode, onToggleDarkMode }) {
  const theme = useTheme();
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useDispatch();
  const [mobileOpen, setMobileOpen] = useState(false);
  const [anchorUser, setAnchorUser] = useState(null);
  const [anchorNotif, setAnchorNotif] = useState(null);

  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const menuItems = [
    { path: '/dashboard', label: 'Дашборд', icon: <DashboardIcon /> },
    { path: '/tasks', label: 'Задачи', icon: <TaskAltIcon /> },
    { path: '/projects', label: 'Проекты', icon: <FolderIcon /> },
  ];

  const handleDrawerToggle = () => setMobileOpen(!mobileOpen);
  const handleUserMenu = (e) => setAnchorUser(e.currentTarget);
  const handleUserClose = () => setAnchorUser(null);
  const handleNotifOpen = (e) => setAnchorNotif(e.currentTarget);
  const handleNotifClose = () => setAnchorNotif(null);

  const handleLogout = () => {
    wsService.disconnect();
    dispatch(logout());
    handleUserClose();
    navigate('/login');
  };

  // Подключение WebSocket при наличии токена
  useEffect(() => {
    const token = localStorage.getItem('accessToken');
    if (token) {
      wsService.connect(token);
      const unsub = wsService.onNotification((notification) => {
        dispatch(addNotification(notification));
      });
      return () => {
        unsub();
        wsService.disconnect();
      };
    }
  }, [dispatch]);

  const drawer = (
    <Box sx={{ pt: 2, pb: 2 }}>
      <Typography variant="h6" sx={{ px: 2, mb: 2, fontWeight: 700, color: 'primary.main' }}>
        TaskMind
      </Typography>
      <List>
        {menuItems.map((item) => (
          <ListItem key={item.path} disablePadding sx={{ px: 1 }}>
            <ListItemButton
              selected={location.pathname === item.path}
              onClick={() => {
                navigate(item.path);
                if (isMobile) setMobileOpen(false);
              }}
              sx={{ borderRadius: 2, mx: 0.5 }}
            >
              <ListItemIcon sx={{ minWidth: 40, color: 'inherit' }}>{item.icon}</ListItemIcon>
              <ListItemText primary={item.label} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      <AppBar
        position="fixed"
        sx={{
          width: { md: `calc(100% - ${DRAWER_WIDTH}px)` },
          ml: { md: `${DRAWER_WIDTH}px` },
          background: theme.palette.mode === 'dark'
            ? 'linear-gradient(135deg, #161b22 0%, #0d1117 100%)'
            : 'linear-gradient(135deg, #1976d2 0%, #1565c0 100%)',
          boxShadow: '0 4px 20px rgba(0,0,0,0.15)',
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { md: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" sx={{ flexGrow: 1, fontWeight: 600 }}>
            TaskMind
          </Typography>

          <IconButton color="inherit" onClick={onToggleDarkMode} sx={{ mr: 0.5 }}>
            {darkMode ? <LightModeIcon /> : <DarkModeIcon />}
          </IconButton>

          <NotificationPanel
            anchorEl={anchorNotif}
            onOpen={handleNotifOpen}
            onClose={handleNotifClose}
          />

          <IconButton color="inherit" onClick={handleUserMenu}>
            <AccountCircleIcon />
          </IconButton>
          <Menu
            anchorEl={anchorUser}
            open={Boolean(anchorUser)}
            onClose={handleUserClose}
            anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            transformOrigin={{ vertical: 'top', horizontal: 'right' }}
            slotProps={{ paper: { sx: { mt: 1.5, minWidth: 180 } } }}
          >
            <MenuItem onClick={() => { navigate('/dashboard'); handleUserClose(); }}>
              Профиль
            </MenuItem>
            <MenuItem onClick={handleLogout}>Выйти</MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      <Box
        component="nav"
        sx={{ width: { md: DRAWER_WIDTH }, flexShrink: { md: 0 } }}
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          slotProps={{ root: { keepMounted: true } }}
          sx={{
            display: { xs: 'block', md: 'none' },
            '& .MuiDrawer-paper': {
              width: DRAWER_WIDTH,
              boxSizing: 'border-box',
              border: 'none',
              boxShadow: theme.shadows[8],
            },
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', md: 'block' },
            '& .MuiDrawer-paper': {
              width: DRAWER_WIDTH,
              boxSizing: 'border-box',
              mt: '64px',
              border: 'none',
              background: theme.palette.background.paper,
              boxShadow: '4px 0 20px rgba(0,0,0,0.05)',
            },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>

      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { md: `calc(100% - ${DRAWER_WIDTH}px)` },
          mt: '64px',
          minHeight: 'calc(100vh - 64px)',
          background: theme.palette.background.default,
        }}
      >
        <Outlet />
      </Box>
    </Box>
  );
}
