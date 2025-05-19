import axios from './axios'

// 认证相关的API
export const authApi = {
  // 用户登录
  login(credentials) {
    return axios.post('/login', credentials)
  },
  
  // 用户注册
  register(userData) {
    return axios.post('/register', userData)
  },
  
  // 获取当前用户信息
  getCurrentUser() {
    return axios.get('/auth/me')
  },
  
  // 更新用户信息
  updateProfile(userData) {
    return axios.put('/auth/profile', userData)
  },
  
  // 修改密码
  changePassword(passwordData) {
    return axios.post('/auth/change-password', passwordData)
  },
  
  // 退出登录
  logout() {
    return axios.post('/auth/logout')
  },
  
  // 重置密码请求
  requestPasswordReset(email) {
    return axios.post('/auth/forgot-password', { email })
  },
  
  // 重置密码
  resetPassword(resetData) {
    return axios.post('/auth/reset-password', resetData)
  },
  
  // 刷新Token
  refreshToken() {
    return axios.post('/api/user/refresh-token')
  }
}