FROM node:18-alpine

WORKDIR /app

# Копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm ci

# Копируем исходники
COPY . .

# Компилируем TypeScript в JavaScript
RUN npm run build

EXPOSE 8001

# Запускаем скомпилированный JS
CMD ["node", "dist/index.js"]