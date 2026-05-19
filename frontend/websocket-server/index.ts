import axios from 'axios';
import express from 'express';
import http from 'http';
import { WebSocketServer, WebSocket } from 'ws';

const port: number = 8001;
const hostname = '0.0.0.0';

const transportLevelPort = parseInt(process.env.TRANSPORT_PORT || '8080', 10);
const transportLevelHostname = process.env.TRANSPORT_HOST || '172.29.58.128'; // IP транспорта

const app = express();
const server = http.createServer(app);
const wss = new WebSocketServer({ server });


app.use(express.json());

app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization');
  
  if (req.method === 'OPTIONS') {
    return res.sendStatus(200);
  }
  next();
});

// Хранилище пользователей
const users: Map<string, WebSocket> = new Map();

// Отправка списка активных пользователей всем подключённым клиентам
function broadcastUserList() {
  const userList = Array.from(users.keys());
  const userListMessage = {
    type: 'user_list',
    users: userList
  };
  
  users.forEach((client) => {
    if (client.readyState === WebSocket.OPEN) {
      client.send(JSON.stringify(userListMessage));
    }
  });
}

app.get('/health', (req, res) => {
  res.json({ status: 'ok', message: 'WebSocket server is running' });
});

// Эндпоинт для суммаризации (фронтенд вызывает его по кнопке "Пересказать")
app.post('/summarize', async (req, res) => {
  console.log('[summarize] Received request:', req.body);
  
  const { messages, username } = req.body;
  
  if (!messages || !messages.length) {
    return res.status(400).json({ error: 'No messages provided' });
  }
  
  const payload = messages.join('. ');
  
  const transportData = {
    sender: username,
    timestamp: new Date().toISOString(),
    payload: payload
  };
  
  console.log(`[summarize] Sending to transport: ${transportLevelHostname}:${transportLevelPort}`);
  
  try {
    const response = await axios.post(
      `http://${transportLevelHostname}:${transportLevelPort}/send`,
      transportData
    );
    res.status(200).json({ status: 'ok' });
  } catch (error: any) {
    console.error('[summarize] Error:', error.message);
    res.status(500).json({ error: 'Failed to send to transport level' });
  }
});

// Эндпоинт для получения данных от транспортного уровня
    app.post('/receive', (req, res) => {
    const { username, data, send_time, error } = req.body;
    
    console.log(`[receive] Received from transport:`);
    console.log(`  username: ${username}`);
    console.log(`  data: ${data.substring(0, 100)}...`);
    
    const message = {
        id: Date.now().toString(),
        username: 'Пересказ',  // ← имя отправителя = Пересказ
        data: data,
        send_time: send_time || new Date().toISOString(),
        error: error === "" ? false : true
    };
    
    // Отправляем только тому пользователю, который указан в username
    const targetUser = users.get(username);
    if (targetUser && targetUser.readyState === WebSocket.OPEN) {
        targetUser.send(JSON.stringify(message));
        console.log(`[receive] Sent summary to user: ${username}`);
    } else {
        console.log(`[receive] User ${username} not found or not connected`);
    }
    
    res.sendStatus(200);
    });

wss.on('connection', (ws: WebSocket, req: any) => {
  // Получаем имя пользователя из URL
  const url = new URL(req.url, `http://${req.headers.host}`);
  const username = url.searchParams.get('username');
  
  if (username) {
    console.log(`[open] Connected, username: ${username}`);
    users.set(username, ws);
    broadcastUserList();
    console.log('Users online:', Array.from(users.keys()));
    
    // Обработка входящих сообщений
    ws.on('message', (data: Buffer) => {
      const messageStr = data.toString();
      console.log(`[message] Received from ${username}: ${messageStr}`);
      
      // Рассылаем сообщение ВСЕМ пользователям КРОМЕ отправителя
      const message = JSON.parse(messageStr);
      users.forEach((client, user) => {
        if (user !== username && client.readyState === WebSocket.OPEN) {
          client.send(JSON.stringify(message));
        }
      });
    });
    
    // Обработка отключения
    ws.on('close', () => {
      console.log(`[close] User disconnected: ${username}`);
      users.delete(username);
      broadcastUserList();
    });
  } else {
    console.log('[open] Connected without username');
    ws.close();
  }
});

server.listen(port, hostname, () => {
  console.log(`Server started at http://${hostname}:${port}`);
  console.log(`WebSocket server is running on ws://${hostname}:${port}`);
});