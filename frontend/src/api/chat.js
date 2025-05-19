import axios from './axios'
import { useUserStore } from '@/stores/user';

// 聊天相关API
export const chatApi = {
  // 获取会话列表
  getSessions() {
    return axios.get('/chat/sessions')
  },
  
  // 创建新会话
  createSession() {
    return axios.post('/chat/sessions', {
      title: '新对话'
    })
  },
  
  // 获取会话消息
  getSessionMessages(sessionId) {
    const userStore = useUserStore();

    return axios.get(`/chat/sessions/${sessionId}/messages`, {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }

    })

  },
  
  // 删除会话
  deleteSession(sessionId) {
    return axios.delete(`/chat/sessions/${sessionId}`)
  },
  
  // 发送消息
  sendMessage(data) {
    const userStore = useUserStore();
    const payload = {
      session_id: Number(data.session_id),
      content: data.content
    };
    return axios.post('/chat/messages', payload, {
      headers: {
        Authorization: `Bearer ${userStore.token}`,
        'Content-Type': 'application/json'
      },
      timeout: 30000
    })
  }
}