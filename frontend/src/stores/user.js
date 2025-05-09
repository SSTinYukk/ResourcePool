import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../api/auth'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  
  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  
  // 动作
  const login = async (credentials) => {
    try {
      const response = await authApi.login(credentials)
      token.value = response.data.token
      user.value = response.data.user
      
      // 保存到本地存储
      localStorage.setItem('token', token.value)
      localStorage.setItem('user', JSON.stringify(user.value))
      
      return response
    } catch (error) {
      console.error('登录失败:', error)
      throw error
    }
  }
  
  const register = async (userData) => {
    try {
      const response = await authApi.register(userData)
      return response
    } catch (error) {
      console.error('注册失败:', error)
      throw error
    }
  }
  
  const logout = () => {
    // 清除状态
    token.value = ''
    user.value = null
    
    // 清除本地存储
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }
  
  const updateProfile = async (userData) => {
    try {
      const response = await authApi.updateProfile(userData)
      user.value = response.data
      localStorage.setItem('user', JSON.stringify(user.value))
      return response
    } catch (error) {
      console.error('更新个人资料失败:', error)
      throw error
    }
  }
  
  return {
    token,
    user,
    isLoggedIn,
    isAdmin,
    login,
    register,
    logout,
    updateProfile
  }
})