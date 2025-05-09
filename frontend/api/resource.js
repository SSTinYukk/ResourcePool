import axios from './axios';

// 上传资源
const uploadResource = async (file, formData) => {
  try {
    const response = await axios.post('/resources/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    return response.data;
  } catch (error) {
    throw error.response.data;
  }
};

// 点赞/取消点赞资源
const likeResource = async (resourceId) => {
  try {
    const response = await axios.post(`/resources/${resourceId}/like`);
    return response.data;
  } catch (error) {
    throw error.response.data;
  }
};

// 收藏/取消收藏资源
const favoriteResource = async (resourceId) => {
  try {
    const response = await axios.post(`/resources/${resourceId}/favorite`);
    return response.data;
  } catch (error) {
    throw error.response.data;
  }
};

// 搜索资源
const searchResources = async (query, page = 1, pageSize = 10, categoryId = '') => {
  try {
    const response = await axios.get('/resources/search', {
      params: {
        q: query,
        page,
        pageSize,
        category: categoryId
      }
    });
    return response.data;
  } catch (error) {
    throw error.response.data;
  }
};

export default {
  uploadResource,
  likeResource,
  favoriteResource,
  searchResources
};