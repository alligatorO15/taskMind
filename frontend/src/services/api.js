// Axios-инстанс для API-запросов с JWT и обработкой 401
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Переменная для хранения функции обновления токена (из authSlice)
let refreshTokenFn = null;
export const setRefreshTokenFn = (fn) => {
  refreshTokenFn = fn;
};

// Переменная для хранения флага повторной попытки
let isRefreshing = false;
let failedQueue = [];

// Обработка очереди запросов при обновлении токена
const processQueue = (error, token = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

// Интерцептор запросов: добавляем JWT из localStorage
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Интерцептор ответов: при 401 пробуем обновить токен
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Если 401 и это не повторная попытка
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // Ждём завершения обновления токена
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            return api(originalRequest);
          })
          .catch((err) => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        if (refreshTokenFn) {
          await refreshTokenFn();
          const newToken = localStorage.getItem('accessToken');
          processQueue(null, newToken);
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
          return api(originalRequest);
        }
      } catch (refreshError) {
        processQueue(refreshError, null);
        // Выход из системы при неудачном обновлении
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

export default api;
