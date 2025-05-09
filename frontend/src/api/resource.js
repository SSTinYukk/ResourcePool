import axios from './axios'

// 资源相关API
export const resourceApi = {
  // 获取资源列表
  getResources(params) {
    return axios.get('/api/resources', { 
      params: {
        page: params.page || 1,
        pageSize: params.pageSize || 12,
        sort: params.sort || 'newest'
      }
    })
  },
  
  // 获取分类列表
  getCategories() {
    return axios.get('/resources/categories')
  },
  
  // 上传文件
  uploadFile(formData) {
    return axios.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      timeout: 30000
    })
  },
  
  // 获取上传文件列表
  getUploadedFiles() {
    return axios.get('/upload/files')
  },
  
  // 删除上传文件
  deleteFile(fileId) {
    return axios.delete(`/upload/files/${fileId}`)
  }
}