// Слайс уведомлений: список, непрочитанные, WebSocket
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import api from '../services/api';

// Загрузка уведомлений
export const fetchNotifications = createAsyncThunk(
  'notification/fetchNotifications',
  async (params = {}, { rejectWithValue }) => {
    try {
      const { data } = await api.get('/notifications', { params });
      return data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.error || 'Ошибка загрузки уведомлений');
    }
  }
);

// Отметить как прочитанное
export const markAsRead = createAsyncThunk(
  'notification/markAsRead',
  async (id, { rejectWithValue }) => {
    try {
      await api.put(`/notifications/${id}/read`);
      return id;
    } catch (err) {
      return rejectWithValue(err.response?.data?.error || 'Ошибка отметки');
    }
  }
);

// Отметить все как прочитанные
export const markAllAsRead = createAsyncThunk(
  'notification/markAllAsRead',
  async (_, { rejectWithValue }) => {
    try {
      await api.put('/notifications/read-all');
    } catch (err) {
      return rejectWithValue(err.response?.data?.error || 'Ошибка отметки');
    }
  }
);

// Загрузка счётчика непрочитанных
export const fetchUnreadCount = createAsyncThunk(
  'notification/fetchUnreadCount',
  async (_, { rejectWithValue }) => {
    try {
      const { data } = await api.get('/notifications/unread-count');
      return data.count ?? data;
    } catch (err) {
      return rejectWithValue(err.response?.data?.error || 'Ошибка загрузки');
    }
  }
);

const notificationSlice = createSlice({
  name: 'notification',
  initialState: {
    notifications: [],
    unreadCount: 0,
    loading: false,
  },
  reducers: {
    // Добавление уведомления из WebSocket
    addNotification: (state, action) => {
      state.notifications.unshift(action.payload);
      state.unreadCount += 1;
    },
  },
  extraReducers: (builder) => {
    builder
      // fetchNotifications
      .addCase(fetchNotifications.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchNotifications.fulfilled, (state, action) => {
        state.loading = false;
        state.notifications = Array.isArray(action.payload)
          ? action.payload
          : action.payload.notifications || [];
      })
      .addCase(fetchNotifications.rejected, (state) => {
        state.loading = false;
      })
      // markAsRead
      .addCase(markAsRead.fulfilled, (state, action) => {
        const n = state.notifications.find((x) => x.id === action.payload);
        if (n) n.read = true;
        state.unreadCount = Math.max(0, state.unreadCount - 1);
      })
      // markAllAsRead
      .addCase(markAllAsRead.fulfilled, (state) => {
        state.notifications.forEach((n) => (n.read = true));
        state.unreadCount = 0;
      })
      // fetchUnreadCount
      .addCase(fetchUnreadCount.fulfilled, (state, action) => {
        state.unreadCount = action.payload;
      });
  },
});

export const { addNotification } = notificationSlice.actions;
export default notificationSlice.reducer;
