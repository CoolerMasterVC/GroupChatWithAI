// Для Docker: REACT_APP_WS_URL будет ws://websocket-server:8001
// Для локальной разработки: ws://localhost:8001
export const WS_URL = process.env.REACT_APP_WS_URL || 'ws://172.29.71.190:8001';
export const API_URL = process.env.REACT_APP_API_URL || 'http://172.29.71.190:8001';