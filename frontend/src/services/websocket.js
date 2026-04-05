// WebSocket-сервис с поддержкой переподключения и уведомлений
const WS_BASE_URL = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/v1/ws`;

class WebSocketService {
  constructor() {
    this.ws = null;
    this.token = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 3000;
    this.reconnectTimer = null;
    this.notificationCallbacks = [];
  }

  // Подключение к WebSocket с токеном
  connect(token) {
    if (!token) return;
    this.token = token;
    this.reconnectAttempts = 0;
    this._connect();
  }

  _connect() {
    const url = `${WS_BASE_URL}?token=${this.token}`;
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
    };

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'notification') {
          this.notificationCallbacks.forEach((cb) => cb(data.payload));
        }
      } catch (e) {
        console.error('Ошибка парсинга WebSocket сообщения:', e);
      }
    };

    this.ws.onclose = () => {
      if (this.token && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectTimer = setTimeout(() => {
          this.reconnectAttempts++;
          this._connect();
        }, this.reconnectDelay);
      }
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket ошибка:', error);
    };
  }

  // Отключение от WebSocket
  disconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    this.reconnectAttempts = this.maxReconnectAttempts;
    this.token = null;
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  // Подписка на уведомления
  onNotification(callback) {
    this.notificationCallbacks.push(callback);
    return () => {
      this.notificationCallbacks = this.notificationCallbacks.filter((cb) => cb !== callback);
    };
  }
}

export const wsService = new WebSocketService();
export default wsService;
