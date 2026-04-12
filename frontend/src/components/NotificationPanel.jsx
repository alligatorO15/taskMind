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
import AlarmIcon from '@mui/icons-material/Alarm';
import WarningAmberIcon from '@mui/icons-material/WarningAmber';
import InfoIcon from '@mui/icons-material/Info';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router';
import {
  fetchNotifications,
  fetchUnreadCount,
  markAllAsRead,
  markAsRead,
} from '../store/notificationSlice';

const NOTIF_ICONS = {
  reminder: <AlarmIcon fontSize="small" color="info" />,
  overdue: <WarningAmberIcon fontSize="small" color="warning" />,
  system: <InfoIcon fontSize="small" color="action" />,
};

export default function NotificationPanel({ anchorEl, onOpen, onClose }) {
  const dispatch = useDispatch();
  const navigate = useNavigate();
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

  const handleClick = (n) => {
    if (!n.read) {
      dispatch(markAsRead(n.id));
    }
    onClose();
    if (n.task_id) {
      navigate('/tasks');
    }
  };

  const formatDate = (dateStr) => {
    if (!dateStr) return null;
    return new Date(dateStr).toLocaleString('ru');
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
                onClick={() => handleClick(n)}
                sx={{
                  bgcolor: n.read ? 'transparent' : 'action.hover',
                  py: 1.5,
                  alignItems: 'flex-start',
                }}
              >
                <ListItemIcon sx={{ mt: 0.5 }}>
                  {NOTIF_ICONS[n.type] || NOTIF_ICONS.system}
                </ListItemIcon>
                <ListItemText
                  primary={n.title || 'Уведомление'}
                  secondary={
                    <>
                      {n.message && (
                        <Typography variant="body2" component="span" display="block" color="text.secondary">
                          {n.message}
                        </Typography>
                      )}
                      <Typography variant="caption" component="span" color="text.disabled">
                        {formatDate(n.created_at)}
                      </Typography>
                    </>
                  }
                  primaryTypographyProps={{ fontWeight: n.read ? 400 : 600 }}
                />
                {n.task_id && (
                  <Typography variant="caption" color="primary" sx={{ ml: 1, mt: 0.5, whiteSpace: 'nowrap' }}>
                    Открыть задачу
                  </Typography>
                )}
              </MenuItem>
            ))
          )}
        </Box>
      </Menu>
    </>
  );
}
