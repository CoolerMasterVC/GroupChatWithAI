import React, { useState, useEffect, useRef } from 'react';
import { useUser } from '../hooks/useUser';
import { Button, TextField, Box, Typography, IconButton } from '@mui/material';
import { Send } from '@mui/icons-material';
import { WS_URL, API_URL, Message } from '../consts';

const AppContent: React.FC = () => {
  const { login, isAuthenticated, loginUser, logoutUser } = useUser();
  const [username, setUsername] = useState('');
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState('');
  const [selectedMessageIds, setSelectedMessageIds] = useState<string[]>([]); // ← массив вместо строки
  const wsRef = useRef<WebSocket | null>(null);
  const [onlineUsers, setOnlineUsers] = useState<string[]>([]);

  // Подключение к WebSocket при авторизации
  useEffect(() => {
    if (isAuthenticated && login) {
      const websocket = new WebSocket(`${WS_URL}/?username=${encodeURIComponent(login)}`);
      wsRef.current = websocket;
      
      websocket.onopen = () => {
        console.log('WebSocket connected');
      };
      
      websocket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        
        // Проверяем тип сообщения
        if (data.type === 'user_list') {
            setOnlineUsers(data.users.filter((user: string) => user !== login));
        } else {
            // Обычное сообщение чата
            const message: Message = data;
            setMessages(prev => [...prev, message]);
        }
        };
      
      websocket.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
      
      setWs(websocket);
      
      return () => {
        if (websocket.readyState === WebSocket.OPEN) {
          websocket.close();
        }
      };
    }
  }, [isAuthenticated, login]);

  // Отправка сообщения
  const sendMessage = () => {
    const currentWs = wsRef.current;
    if (currentWs && currentWs.readyState === WebSocket.OPEN && inputMessage.trim() && login && selectedMessageIds.length === 0) {
      const message: Message = {
        id: Date.now().toString(),
        username: login,
        data: inputMessage,
        send_time: new Date().toISOString(),
        error: false,
      };
      currentWs.send(JSON.stringify(message));
      setMessages(prev => [...prev, message]);
      setInputMessage('');
    }
  };

  // Обработка выделения сообщения (несколько сообщений)
  const handleMessageSelect = (messageId: string, e: React.MouseEvent) => {
    e.stopPropagation();
    setSelectedMessageIds(prev => {
      if (prev.includes(messageId)) {
        // Если уже выделено — убираем
        return prev.filter(id => id !== messageId);
      } else {
        // Если не выделено — добавляем
        return [...prev, messageId];
      }
    });
  };

  // Пересказ выделенных сообщений
    const handleSummarize = async () => {
    if (selectedMessageIds.length === 0) return;
    
    const selectedMessages = messages.filter(msg => selectedMessageIds.includes(msg.id));
    const selectedTexts = selectedMessages.map(msg => msg.data);
    
    try {
        // Используем API_URL вместо относительного пути
        const response = await fetch(`${API_URL}/summarize`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            messages: selectedTexts,
            username: login
        }),
        });
        
        if (response.ok) {
        console.log('Summarize request sent successfully');
        setSelectedMessageIds([]);
        } else {
        console.error('Failed to send summarize request:', response.status);
        }
    } catch (error) {
        console.error('Error sending summarize request:', error);
    }
    };

  // Выход из чата
  const handleLogout = () => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.close();
    }
    logoutUser();
    setMessages([]);
    setUsername('');
    setSelectedMessageIds([]);
  };

  const handleLogin = () => {
    if (username.trim()) {
      loginUser(username);
    }
  };

  const otherUsers = [...new Set(messages.filter(m => m.username !== login && m.username !== 'Пересказ').map(m => m.username))];

  const isSelected = (messageId: string) => selectedMessageIds.includes(messageId);
  const hasSelectedMessages = selectedMessageIds.length > 0;

  return (
    <Box sx={{ 
      width: '100vw',
      height: '100vh',
      background: 'linear-gradient(180deg, #1F50E8 0%, #587CED 100%)',
      fontFamily: 'Lato, sans-serif',
      position: 'relative',
      overflow: 'hidden'
    }}>
      {/* Верхняя панель */}
      <Box sx={{ 
        width: '100%',
        height: { xs: '56px', sm: '64px' },
        bgcolor: '#FFFFFF',
        position: 'fixed',
        top: 0,
        left: 0,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        boxShadow: '0px 1px 0px rgba(0,0,0,0.05)',
        borderRadius: '0px',
        zIndex: 10,
        px: { xs: 2, sm: 3 }
      }}>
        {isAuthenticated && (
          <>
            <Typography sx={{ 
              fontSize: { xs: '12px', sm: '14px', md: '16px' },
              color: '#0D59D4',
              fontFamily: 'Lato, sans-serif',
              fontWeight: 400,
              textAlign: 'center'
            }}>
              Сообщения от {onlineUsers.length > 0 ? onlineUsers.join(', ') : 'никого'}
            </Typography>
            
            <Button 
              onClick={handleLogout}
              sx={{ 
                position: 'absolute',
                right: { xs: '16px', sm: '24px', md: 'calc(50% - 402px + 16px)' },
                width: { xs: '60px', sm: '70px', md: '77px' },
                height: { xs: '32px', sm: '36px', md: '40px' },
                bgcolor: '#EDF2FE',
                color: '#0D59D4',
                fontSize: { xs: '12px', sm: '14px', md: '16px' },
                fontFamily: 'Lato, sans-serif',
                textTransform: 'none',
                borderRadius: '10px',
                '&:hover': { bgcolor: '#EDF2FE' }
              }}
            >
              Выход
            </Button>
          </>
        )}
      </Box>

      {/* Окно входа */}
      {!isAuthenticated && (
        <Box sx={{ 
          position: 'fixed',
          width: { xs: '320px', sm: '350px', md: '390px' },
          height: { xs: '300px', sm: '320px', md: '340px' },
          left: '50%',
          top: '50%',
          transform: 'translate(-50%, -50%)',
          bgcolor: 'white',
          borderRadius: '10px',
          boxShadow: '0px 4px 20px rgba(0,0,0,0.1)',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          zIndex: 5
        }}>
          <Typography sx={{ 
            fontSize: { xs: '24px', sm: '28px', md: '32px' },
            color: '#0065B1',
            fontFamily: 'Lato, sans-serif',
            fontWeight: 400,
            mt: { xs: '30px', sm: '35px', md: '43px' },
            mb: { xs: '25px', sm: '35px', md: '41px' }
          }}>
            Вход
          </Typography>
          
          <TextField
            placeholder="Имя / Логин"
            variant="outlined"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            onKeyPress={(e) => {
              if (e.key === 'Enter') {
                handleLogin();
              }
            }}
            sx={{ 
              width: { xs: '260px', sm: '290px', md: '329px' },
              '& .MuiOutlinedInput-root': {
                height: { xs: '48px', sm: '52px', md: '57px' },
                bgcolor: '#F5F7FA',
                borderRadius: '10px',
                '& fieldset': { borderColor: 'transparent' },
                '&:hover fieldset': { borderColor: '#0D4CD3' },
                '&.Mui-focused fieldset': { borderColor: '#0D4CD3' }
              },
              '& .MuiInputBase-input': {
                fontSize: { xs: '14px', sm: '15px', md: '16px' },
                fontFamily: 'Lato, sans-serif',
                color: '#0B1F33',
                '&::placeholder': {
                  color: '#A7AEB6',
                  fontSize: { xs: '14px', sm: '15px', md: '16px' },
                  fontFamily: 'Lato, sans-serif'
                }
              }
            }}
          />
          
          <Button
            variant="contained"
            onClick={handleLogin}
            disabled={!username.trim()}
            sx={{ 
              width: { xs: '260px', sm: '290px', md: '329px' },
              height: { xs: '48px', sm: '52px', md: '57px' },
              bgcolor: '#0D4CD3',
              mt: { xs: '25px', sm: '35px', md: '43px' },
              fontSize: { xs: '14px', sm: '15px', md: '16px' },
              fontFamily: 'Lato, sans-serif',
              textTransform: 'none',
              borderRadius: '10px',
              color: '#FFFFFF',
              '&:hover': { bgcolor: '#0D4CD3' },
              '&.Mui-disabled': {
                bgcolor: '#0D4CD3',
                color: '#FFFFFF',
                opacity: 1
              }
            }}
          >
            Войти
          </Button>
        </Box>
      )}

      {/* Область сообщений */}
      {isAuthenticated && (
        <Box sx={{ 
          position: 'fixed',
          left: 0,
          right: 0,
          top: { xs: '72px', sm: '80px', md: '84px' },
          bottom: { xs: '82px', sm: '92px', md: '101px' },
          overflow: 'auto',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          zIndex: 5
        }}>
          <Box sx={{ 
            width: { xs: 'calc(100% - 32px)', sm: 'calc(100% - 48px)', md: '804px' },
            maxWidth: '804px',
            display: 'flex',
            flexDirection: 'column',
            gap: '12px'
          }}>
            {messages.map((msg) => {
              const selected = isSelected(msg.id);
              return (
                <Box
                  key={msg.id}
                  sx={{
                    display: 'flex',
                    justifyContent: msg.username === login ? 'flex-end' : 'flex-start',
                  }}
                >
                  <Box 
                    sx={{ 
                      maxWidth: { xs: '85%', sm: '80%', md: '70%' },
                      minWidth: { xs: '100px', sm: '120px' },
                      opacity: selected ? 0.7 : 1,
                      transition: 'opacity 0.2s ease'
                    }}
                    onMouseDown={(e) => handleMessageSelect(msg.id, e)}
                  >
                    <Typography sx={{ 
                      fontSize: '18px',
                      color: '#FFFFFF',
                      fontFamily: 'Lato, sans-serif',
                      mb: '4px',
                      ml: '12px',
                      opacity: 0.8
                    }}>
                      {msg.username}
                    </Typography>
                    <Box sx={{ 
                      bgcolor: selected ? 'rgba(255, 255, 255, 0.3)' : '#FFFFFF',
                      borderRadius: '10px',
                      py: { xs: '8px', sm: '10px', md: '12px' },
                      px: { xs: '12px', sm: '16px', md: '20px' },
                      boxShadow: '0px 1px 2px rgba(0,0,0,0.05)',
                      transition: 'background-color 0.2s ease'
                    }}>
                      <Typography sx={{ 
                        fontSize: { xs: '14px', sm: '16px', md: '18px' },
                        color: selected ? '#FFFFFF' : '#0B1F33',
                        fontFamily: 'Lato, sans-serif',
                        wordBreak: 'break-word'
                      }}>
                        {msg.error ? '⚠️ Ошибка доставки' : msg.data}
                      </Typography>
                    </Box>
                  </Box>
                </Box>
              );
            })}
          </Box>
        </Box>
      )}

      {/* Поле ввода сообщения */}
      <Box sx={{ 
        position: 'fixed',
        left: 0,
        right: 0,
        bottom: { xs: '16px', sm: '18px', md: '21px' },
        display: 'flex',
        justifyContent: 'center',
        zIndex: 10
      }}>
        <Box sx={{ 
          width: { xs: 'calc(100% - 32px)', sm: 'calc(100% - 48px)', md: '804px' },
          maxWidth: '804px',
          bgcolor: '#FFFFFF',
          borderRadius: '10px',
          boxShadow: '0px 2px 8px rgba(0,0,0,0.05)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: 0,
          height: { xs: '50px', sm: '55px', md: '60px' },
          overflow: 'hidden'
        }}>
          <TextField
            placeholder="Сообщение..."
            variant="standard"
            value={inputMessage}
            onChange={(e) => setInputMessage(e.target.value)}
            onKeyPress={(e) => {
              if (e.key === 'Enter' && isAuthenticated && !hasSelectedMessages) {
                sendMessage();
              }
            }}
            disabled={!isAuthenticated || hasSelectedMessages}
            sx={{ 
              flex: 1,
              px: { xs: '16px', sm: '20px', md: '24px' },
              '& .MuiInputBase-root': {
                fontSize: { xs: '14px', sm: '16px', md: '18px' },
                fontFamily: 'Lato, sans-serif',
                color: '#0B1F33'
              },
              '& .MuiInputBase-input': {
                '&::placeholder': {
                  color: '#66727F',
                  fontSize: { xs: '14px', sm: '16px', md: '18px' },
                  fontFamily: 'Lato, sans-serif'
                }
              },
              '& .MuiInput-underline:before': {
                borderBottom: 'none'
              },
              '& .MuiInput-underline:hover:before': {
                borderBottom: 'none'
              },
              '& .MuiInput-underline:after': {
                borderBottom: 'none'
              }
            }}
            slotProps={{
              input: {
                disableUnderline: true
              }
            }}
          />
          
          {hasSelectedMessages ? (
            <Button
              onClick={handleSummarize}
              sx={{ 
                width: { xs: '100px', sm: '110px', md: '116px' },
                height: '100%',
                bgcolor: '#0D4CD3',
                fontSize: { xs: '14px', sm: '16px', md: '18px' },
                fontFamily: 'Lato, sans-serif',
                color: '#FFFFFF',
                textTransform: 'none',
                borderRadius: 0,
                '&:hover': { bgcolor: '#0B3DA8' },
                flexShrink: 0
              }}
            >
              Пересказать
            </Button>
          ) : (
            <Box sx={{ 
              px: { xs: '16px', sm: '20px', md: '24px' },
              display: 'flex',
              alignItems: 'center'
            }}>
              <IconButton 
                onClick={sendMessage}
                disabled={!isAuthenticated || !inputMessage.trim()}
                sx={{ 
                  p: 0,
                  borderRadius: '10px',
                  '&:hover': { bgcolor: 'transparent' }
                }}
              >
                <Send sx={{ 
                  color: '#0D4CD3',
                  fontSize: { xs: '20px', sm: '22px', md: '24px' }
                }} />
              </IconButton>
            </Box>
          )}
        </Box>
      </Box>
    </Box>
  );
};

export default AppContent;