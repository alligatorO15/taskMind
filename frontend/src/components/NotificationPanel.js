// Панель уведомлений: выпадающий список, кнопка «отметить все прочитанными»
import React, { useEffect } from 'react';
import {
  IconButton,
  Badge,
  Menu,
  MenuItem,
  ListItemText,
  ListItemIcon,
  Typography,
  Button,
  Box,
  Divider,
} from '@mui/material';
import NotificationsIcon from '@mui/icons-material/Notifications';
import DoneAllIcon from '@mui/icons-material/DoneAll';
import { useDispatch, useSelector } from 'react-redux';
import {
  fetchNotifications,
  fetchUnreadCount,
  markAllAsRead,
  markAsRead,
} from '../store/notificationSlice';

export default function NotificationPanel({ anchorEl, onOpen, onClose }) {
  const dispatch = useDispatch();
  const { notifications, unreadCount, loading } = useSelector((state) => state.notification);
  const open = Boolean(anchorEl);

  useEffect(() => {
    if (open) {
      dispatch(fetchNotifications());
      dispatch(fetchUnreadCount());
    }
  }, [open, dispatch]);

  const handleMarkAllRead = () => {
    dispatch(markAllAsRead());
    onClose();
  };

  const handleMarkRead = (id) => {
    dispatch(markAsRead(id));
  };

  return (
    <>
      <IconButton color="inherit" onClick={onOpen}>
        <Badge badgeContent={unreadCount} color="error">
          <NotificationsIcon />
        </Badge>
      </IconButton>
      <Menu
        anchorEl={anchorEl}
        open={open}
        onClose={onClose}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
        transformOrigin={{ vertical: 'top', horizontal: 'right' }}
        slotProps={{
          paper: {
            sx: {
              mt: 1.5,
              minWidth: 360,
              maxHeight: 400,
              boxShadow: '0 8px 32px rgba(0,0,0,0.12)',
              borderRadius: 2,
            },
          },
        }}
      >
        <Box sx={{ px: 2, py: 1.5, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="subtitle1" fontWeight={600}>
            Уведомления
          </Typography>
          {unreadCount > 0 && (
            <Button size="small" startIcon={<DoneAllIcon />} onClick={handleMarkAllRead}>
              Прочитать все
            </Button>
          )}
        </Box>
        <Divider />
        <Box sx={{ maxHeight: 320, overflow: 'auto' }}>
          {loading ? (
            <Box sx={{ p: 3, textAlign: 'center' }}>
              <Typography color="text.secondary">Загрузка...</Typography>
            </Box>
          ) : notifications.length === 0 ? (
            <Box sx={{ p: 3, textAlign: 'center' }}>
              <Typography color="text.secondary">Нет уведомлений</Typography>
            </Box>
          ) : (
            notifications.map((n) => (
              <MenuItem
                key={n.id}
                onClick={() => handleMarkRead(n.id)}
                sx={{
                  bgcolor: n.read ? 'transparent' : 'action.hover',
                  py: 1.5,
                }}
              >
                <ListItemIcon>
                  {!n.read && <Box sx={{ width: 8, height: 8, borderRadius: '50%', bgcolor: 'primary.main' }} />}
                </ListItemIcon>
                <ListItemText
                  primary={n.title || n.message || 'Уведомление'}
                  secondary={n.createdAt ? new Date(n.createdAt).toLocaleString('ru') : null}
                  primaryTypographyProps={{ fontWeight: n.read ? 400 : 600 }}
                />
              </MenuItem>
            ))
          )}
        </Box>
      </Menu>
    </>
  );
}
