import axios from 'axios';
import { useUserStore } from '@/stores/user';

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
});

// 添加请求拦截器
api.interceptors.request.use(config => {
  const userStore = useUserStore();
  if (userStore.token) {
    config.headers.Authorization = `Bearer ${userStore.token}`;
  }
  return config;
}, error => {
  return Promise.reject(error);
});

export const forumApi = {
  // 获取论坛分类
  getCategories() {
    return api.get('/forum/categories');
  },
  
  // 获取话题列表
  getTopics(params) {
    return api.get('/forum/topics', { params });
  },
  
  // 获取话题详情
  getTopicById(id) {
    return api.get(`/forum/topics/${id}`);
  },
  
  // 创建话题
  createTopic(data) {
    return api.post('/forum/topics', data);
  },
  
  // 更新话题
  updateTopic(id, data) {
    return api.put(`/forum/topics/${id}`, data);
  },
  
  // 删除话题
  deleteTopic(id) {
    return api.delete(`/forum/topics/${id}`);
  },
  
  // 创建回复
  createReply(topicId, data) {
    return api.post(`/forum/topics/${topicId}/replies`, data);
  },
  
  // 更新回复
  updateReply(id, data) {
    return api.put(`/forum/replies/${id}`, data);
  },
  
  // 删除回复
  deleteReply(id) {
    return api.delete(`/forum/replies/${id}`);
  }
};