// API服务集合
import axios from './axios'

// 资源相关API
export const resourceApi = {
  // 获取资源列表
  getResources(params) {
    return axios.get('/resources', { params })
  },
  // 获取资源详情
  getResourceById(id) {
    return axios.get(`/resources/${id}`)
  },
  // 获取资源分类
  getCategories() {
    return axios.get('/resources/categories')
  },
  // 搜索资源
  searchResources(query) {
    return axios.get('/resources/search', { params: { query } })
  },
  // 创建资源
  createResource(data) {
    return axios.post('/resources', data)
  },
  // 更新资源
  updateResource(id, data) {
    return axios.put(`/resources/${id}`, data)
  },
  // 删除资源
  deleteResource(id) {
    return axios.delete(`/resources/${id}`)
  },
  // 获取用户资源
  getUserResources() {
    return axios.get('/user/resources')
  },
  // 上传文件
  uploadFile(formData) {
    return axios.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },
  // 下载文件
  downloadFile(id) {
    return axios.get(`/download/${id}`, { responseType: 'blob' })
  }
}

// 论坛相关API
export const forumApi = {
  // 获取论坛分类
  getCategories() {
    return axios.get('/forum/categories')
  },
  // 获取主题列表
  getTopics(params) {
    return axios.get('/forum/topics', { params })
  },
  // 获取主题详情
  getTopicById(id) {
    return axios.get(`/forum/topics/${id}`)
  },
  // 创建主题
  createTopic(data) {
    return axios.post('/forum/topics', data)
  },
  // 更新主题
  updateTopic(id, data) {
    return axios.put(`/forum/topics/${id}`, data)
  },
  // 删除主题
  deleteTopic(id) {
    return axios.delete(`/forum/topics/${id}`)
  },
  // 创建回复
  createReply(topicId, data) {
    return axios.post(`/forum/topics/${topicId}/replies`, data)
  },
  // 更新回复
  updateReply(id, data) {
    return axios.put(`/forum/replies/${id}`, data)
  },
  // 删除回复
  deleteReply(id) {
    return axios.delete(`/forum/replies/${id}`)
  }
}



// 积分相关API
export const pointsApi = {
  // 获取用户积分
  getUserPoints() {
    return axios.get('/user/points')
  },
  // 获取积分历史
  getPointsHistory() {
    return axios.get('/user/points/history')
  }
}

// 管理员相关API
export const adminApi = {
  // 获取待审核资源
  getPendingResources() {
    return axios.get('/admin/resources/pending')
  },
  // 审核资源
  reviewResource(id, data) {
    return axios.put(`/admin/resources/${id}/review`, data)
  },
  // 添加积分
  addPoints(data) {
    return axios.post('/admin/points/add', data)
  },
  // 获取用户统计信息
  getUserStats() {
    return axios.get('/admin/stats')
  }
}