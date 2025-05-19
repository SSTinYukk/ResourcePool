import axios from 'axios'

const likeService = {
  // 点赞主题
  likeTopic: async (topicId) => {
    try {
      const response = await axios.post(`/api/forum/topics/${topicId}/like`, {}, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      return response.data
    } catch (error) {
      throw error
    }
  },

  // 点踩主题
  dislikeTopic: async (topicId) => {
    try {
      const response = await axios.post(`/api/forum/topics/${topicId}/dislike`, {}, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      return response.data
    } catch (error) {
      throw error
    }
  }
}

export { likeService }