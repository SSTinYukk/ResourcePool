import axios from './axios'

// 管理员相关API
export const adminApi = {
  // 用户管理
  getUsers(params) {
    return axios.get('/admin/users', { params })
  },
  deleteUser(id) {
    return axios.delete(`/admin/users/${id}`)
  },
  updateUserRole(id, role) {
    return axios.put(`/admin/users/${id}/role`, { role })
  },
  
  // 资源管理
  getResources(params) {
    return axios.get('/admin/resources', { params })
  },
  reviewResource(id, data) {
    return axios.put(`/admin/resources/${id}/review`, data)
  },
  deleteResource(id) {
    return axios.delete(`/admin/resources/${id}`)
  },
  
  // 论坛管理
  getTopics(params) {
    return axios.get('/admin/forum/topics', { params })
  },
  deleteTopic(id) {
    return axios.delete(`/admin/forum/topics/${id}`)
  },
  
  // 积分管理
  getPointsRecords(params) {
    return axios.get('/admin/points/records', { params })
  },
  adjustPoints(data) {
    return axios.post('/admin/points/adjust', data)
  },
  
  // 统计信息
  getStats() {
    return axios.get('/admin/stats')
  },
  getUserStats() {
    return axios.get('/admin/stats/users')
  },
  getResourceStats() {
    return axios.get('/admin/stats/resources')
  },
  getForumStats() {
    return axios.get('/admin/stats/forum')
  }
}