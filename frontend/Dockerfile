# Этап 1: Сборка React-приложения (Vite)
FROM node:22-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей для кеширования
COPY package.json package-lock.json* ./
RUN npm install --silent

# Копируем исходный код и собираем production-сборку
COPY . .
RUN npm run build

# Этап 2: Раздача статики через Nginx
FROM nginx:alpine

# Копируем собранное приложение
COPY --from=builder /app/dist /usr/share/nginx/html

# Копируем конфигурацию Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

# сам процесс - рабочий, без демонизации. Чтобы docker не убивал контейнер
CMD ["nginx", "-g", "daemon off;"]
