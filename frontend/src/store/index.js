// Redux store: объединение всех слайсов
import { configureStore } from '@reduxjs/toolkit';
import authReducer from './authSlice';
import taskReducer from './taskSlice';
import projectReducer from './projectSlice';
import notificationReducer from './notificationSlice';
import { setupRefreshToken } from './authSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    task: taskReducer,
    project: projectReducer,
    notification: notificationReducer,
  },
});

// Регистрируем функцию обновления токена для api.js
setupRefreshToken(store.dispatch);
