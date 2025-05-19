<template>
  <div class="bg-white rounded-lg shadow-lg p-8 max-w-5xl mx-auto">
    <!-- 主题标题和操作 -->
    <div class="mb-8">
      <div class="flex justify-between items-start">
        <div>
          <h1 class="text-2xl font-bold text-gray-800 mb-2">{{ topic.title }}</h1>
          <div class="flex items-center text-sm text-gray-500">
            <span>{{ formatDate(topic.createTime) }}</span>
            <span class="mx-2">•</span>
            <span>{{ topic.views }} 次查看</span>
            <span class="mx-2">•</span>
            <span>{{ topic.replies.length }} 条回复</span>
            <Tag :value="getCategoryName(topic.category)" class="ml-2" severity="info" />
          </div>
        </div>
        <div class="flex space-x-2">
          <Button 
            icon="pi pi-bookmark" 
            :class="['p-button-text p-button-rounded', favoriteData.isFavorited ? 'p-button-primary text-yellow-600 ' : '']" 
            @click="toggleFavorite"
            v-tooltip.top="'收藏'"
          />
          <Button icon="pi pi-share-alt" class="p-button-text p-button-rounded" v-tooltip.top="'分享'" />
          <Button icon="pi pi-ellipsis-v" class="p-button-text p-button-rounded" @click="toggleMenu($event)" />
          <Menu ref="menu" :model="menuItems" :popup="true" />
        </div>
      </div>
    </div>

    <!-- 主题内容 -->
    <div class="border-t border-b border-gray-200 py-8 px-6 bg-gray-50 rounded-lg my-6">
      <div class="flex mb-4">
        <div class="mr-4">
          <Avatar :image="topic.authorAvatar" :label="!topic.authorAvatar && topic.author ? topic.author.charAt(0).toUpperCase() : undefined" shape="circle" size="large" />
          <div class="text-center mt-2 text-sm font-medium">{{ topic.author }}</div>
        </div>
        <div class="flex-1">
          <div class="prose max-w-none" v-html="topic.content"></div>
        </div>
      </div>
      <div class="flex justify-end mt-4 space-x-2">
        <div class="like-button-container">
          <Button 
            icon="pi pi-thumbs-up" 
            :class="['p-button-text p-button-sm like-button', topic.userLiked ? 'p-button-primary text-blue-600 liked' : '']"
            @click="likeTopic"
            v-tooltip.top="'点赞'"
          />
          <span :class="['text-sm text-gray-500 like-count', topic.userLiked ? 'text-blue-600' : '']">{{ topic.likes }}</span>
        </div>
        <div class="like-button-container">
          <Button 
            icon="pi pi-thumbs-down" 
            :class="['p-button-text p-button-sm dislike-button', topic.userDisliked ? 'text-red-500 disliked' : '']"
            @click="dislikeTopic"
            v-tooltip.top="'点踩'"
          />
          <span :class="['text-sm text-gray-500 dislike-count', topic.userDisliked ? 'text-red-500' : '']">{{ topic.dislikes }}</span>
        </div>
      </div>
    </div>

    <!-- 回复列表 -->
    <div class="mt-10 bg-white rounded-lg shadow-sm p-6">
      <h2 class="text-xl font-bold text-gray-800 mb-4">{{ topic.replies.length }} 条回复</h2>
      
      <div v-for="(reply, index) in topic.replies" :key="index" class="border-b border-gray-200 py-4 last:border-b-0">
        <div class="flex">
          <div class="mr-4">
            <Avatar :image="reply?.authorAvatar" :label="!reply?.authorAvatar && reply?.author ? reply.author.charAt(0).toUpperCase() : undefined" shape="circle" />
            <div class="text-center mt-2 text-sm font-medium">{{ reply?.author || '匿名用户' }}</div>
          </div>
          <div class="flex-1">
            <div class="flex justify-between items-start mb-2">
              <div class="text-sm text-gray-500">
                <span>#{{ index + 1 }}</span>
                <span class="mx-2">•</span>
                <span>{{ formatDate(reply.createTime) }}</span>
              </div>
              <Button icon="pi pi-ellipsis-h" class="p-button-text p-button-rounded p-button-sm" />
            </div>
            <div class="prose max-w-none" v-html="reply.content"></div>
            <div class="flex justify-end mt-4">
              <Button icon="pi pi-reply" class="p-button-text p-button-sm" label="回复" @click="replyTo(reply)" />
            </div>
          </div>
        </div>
      </div>

      <!-- 无回复时显示 -->
      <div v-if="topic.replies.length === 0" class="text-center py-8">
        <i class="pi pi-comments text-5xl text-gray-300 mb-4"></i>
        <p class="text-gray-500">暂无回复，成为第一个回复者吧！</p>
      </div>
    </div>

    <!-- 回复编辑器 -->
    <div class="mt-10 bg-white rounded-lg shadow-sm p-6">
      <h3 class="text-lg font-bold text-gray-800 mb-4">发表回复</h3>
      <div v-if="replyingTo" class="bg-blue-50 p-3 rounded-md mb-4 flex justify-between items-center">
        <div>回复 <span class="font-medium">{{ replyingTo.author }}</span>: {{ replyingTo.content.substring(0, 50) }}{{ replyingTo.content.length > 50 ? '...' : '' }}</div>
        <Button icon="pi pi-times" class="p-button-text p-button-rounded p-button-sm" @click="cancelReply" />
      </div>
      <InputText v-model="newReply" class="w-full h-40 p-4 rounded-lg border border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-colors duration-300 text-left pt-2" placeholder="输入回复内容..." @keydown.enter="submitReply" />
      <div class="flex justify-end mt-4">
        <Button label="发表回复" icon="pi pi-send" @click="submitReply" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Tag from 'primevue/tag'
import Menu from 'primevue/menu'
import InputText from 'primevue/inputtext'
import { useToast } from 'primevue/usetoast'

const route = useRoute()
const router = useRouter()
const menu = ref(null)
const newReply = ref('')
const replyingTo = ref(null)
const toast = useToast()
const loading = ref(true)
const error = ref(null)

// 从API获取分类数据
const categories = ref([])

// 主题数据
const topic = ref({
    id: 0,
    title: '',
    author: '',
    authorAvatar: null,
    category: 0,
    createTime: new Date(),
    views: 0,
    likes: 0,
    dislikes: 0,
    content: '',
    replies: [],
    userLiked: false,
    userDisliked: false,
    userFavorited: false
  })

// 菜单项
const menuItems = [
  {
    label: '编辑主题',
    icon: 'pi pi-pencil',
    command: () => router.push(`/forum/edit/${topic.value.id}`)
  },
  {
    label: '删除主题',
    icon: 'pi pi-trash',
    command: () => confirmDelete()
  }
]

// 获取分类名称
function getCategoryName(categoryId) {
  const category = categories.value.find(c => c.id === categoryId)
  return category ? category.name : '未分类'
}

// 格式化日期
function formatDate(date) {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

// 显示菜单
function toggleMenu(event) {
  menu.value.toggle(event)
}

// 确认删除
function confirmDelete() {
  // 实际应用中应使用确认对话框
  if (confirm('确定要删除这个主题吗？')) {
    // 调用API删除主题
    router.push('/forum')
  }
}

// 回复某人
function replyTo(reply) {
  replyingTo.value = reply
  // 滚动到编辑器
  document.querySelector('.p-editor-container').scrollIntoView({ behavior: 'smooth' })
}

// 取消回复
function cancelReply() {
  replyingTo.value = null
}


// 收藏状态数据
const favoriteData = ref({ isFavorited: false })

// 加载主题详情
async function loadTopic(id) {
      // 验证ID有效性
      if (!id || id === 'undefined') {
        id = route.params.id
        if (!id || id === 'undefined') {
          toast.add({
            severity: 'error',
            summary: '参数错误',
            detail: '主题ID不能为空',
            life: 3000
          })
          router.push('/forum')
          return
        }
      }

      try {
        loading.value = true
        error.value = null
        
        console.log(`正在加载主题ID: ${id}`)
        
        // 获取主题点赞信息和收藏状态
        const [likesResponse, favoriteResponse, response] = await Promise.all([
          axios.get(`/api/forum/topics/${id}/likes`, {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
          }),
          axios.get(`/api/forum/topics/${id}/favorite-status`, {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
          }),
          axios.get(`/api/forum/topics/${id}`)
        ])
        
        // 处理收藏状态API响应
        if (favoriteResponse?.data) {
          if (favoriteResponse.data.isFavorited !== undefined) {
            favoriteData.value.isFavorited = favoriteResponse.data.isFavorited
          } else if (favoriteResponse.data.is_favorite !== undefined) {
            favoriteData.value.isFavorited = favoriteResponse.data.is_favorite
          } else if (typeof favoriteResponse.data === 'boolean') {
            favoriteData.value.isFavorited = favoriteResponse.data
          }
          console.log('初始化收藏状态:', favoriteData.value.isFavorited)
          console.log('喜欢数和点踩数:', likesResponse.data)
        
        }

    // 验证API返回数据结构
    if (!response.data || !response.data.topic || !response.data.topic.id) {
      console.error('API返回数据格式不正确:', response.data)
      throw new Error('API返回数据格式不正确')
    }
    
    // 合并点赞信息
    const likesData = likesResponse?.data || {}
    const likes = likesData.likes || 0
    const dislikes = likesData.dislikes || 0
    const userLiked = likesData.userLiked || false
    const userDisliked = likesData.userDisliked || false
    
    // 处理收藏状态
    console.log('收藏状态API响应:', favoriteResponse?.data)
    
    // 处理收藏状态API响应
    if (favoriteResponse?.data) {
      if (favoriteResponse.data.isFavorited !== undefined) {
        favoriteData.value.isFavorited = favoriteResponse.data.isFavorited
      } else if (favoriteResponse.data.is_favorite !== undefined) {
        favoriteData.value.isFavorited = favoriteResponse.data.is_favorite
      } else if (typeof favoriteResponse.data === 'boolean') {
        favoriteData.value.isFavorited = favoriteResponse.data
      }
    }
    console.log('处理后的收藏状态:', favoriteData.value.isFavorited)
    
    // 映射后端字段到前端格式
    const backendData = response.data.topic
    const mappedTopic = {
      id: backendData.id,
      title: backendData.title,
      author: backendData.user?.username || '未知用户',
      authorAvatar: backendData.user?.avatar || null,
      category: backendData.category_id,
      createTime: backendData.created_at ? new Date(backendData.created_at) : new Date(),
      views: backendData.view_count || 0,
      likes: likes,
      dislikes: dislikes,
      userLiked: userLiked,
      userDisliked: userDisliked,
      content: backendData.content,
      replies: []
    }
    topic.value = mappedTopic
    console.log('加载的主题:', topic.value)
    console.log(mappedTopic.userLiked)
    
    // 确保回复数据是数组
    if (!Array.isArray(response.data.replies)) {
      console.log('API返回的replies不是数组，已初始化为空数组')
      response.data.replies = []
    }
    
    // 转换日期格式
    if (typeof response.data.createTime === 'string') {
      response.data.createTime = new Date(response.data.createTime)
    }
    
    // 确保回复中的日期格式正确并映射用户信息
    response.data.replies.forEach(reply => {
      if (reply.createTime && typeof reply.createTime === 'string') {
        reply.createTime = new Date(reply.createTime)
      }
      // 映射用户信息
      reply.author = reply.user?.username || '匿名用户'
      reply.authorAvatar = reply.user?.avatar || null
    })
    
    console.log('主题数据加载成功:', mappedTopic)
    topic.value = mappedTopic
    topic.value.replies = response.data.replies || []
    loading.value = false
  } catch (err) {
    error.value = err
    loading.value = false
    
    console.error('加载主题详情失败:', err)
    console.error('错误详情:', {
      message: err.message,
      response: err.response,
      request: err.request,
      config: err.config
    })
    
    const errorDetail = err.response?.status === 404 
      ? '主题不存在或已被删除'
      : err.response?.data?.message || err.message
    
    toast.add({
      severity: 'error',
      summary: '加载失败',
      detail: `无法加载主题详情: ${errorDetail}`,
      life: 3000
    })
    
    // 仅当不是参数错误时才跳转
    if (!id || !isNaN(Number(id))) {
      router.push('/forum')
    }
  }
}
// 收藏主题
async function toggleFavorite() {
  try {
    let response;
    const wasFavorited = favoriteData.value.isFavorited;
    console.log('当前收藏状态:', wasFavorited);
    // 立即更新UI状态
    favoriteData.value.isFavorited = !wasFavorited;
    
    if (wasFavorited) {
      response = await axios.delete(`/api/forum/topics/${topic.value.id}/favorite`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
    } else {
      response = await axios.post(`/api/forum/topics/${topic.value.id}/favorite`, {}, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
    }
    
    // 确保最终状态与API响应一致
    if (response.data.isFavorited !== undefined) {
      favoriteData.value.isFavorited = response.data.isFavorited;
    } else if (response.data.is_favorite !== undefined) {
      favoriteData.value.isFavorited = response.data.is_favorite;
    } else if (typeof response.data === 'boolean') {
      favoriteData.value.isFavorited = response.data;
    }
    console.log('更新后的收藏状态:', favoriteData.value.isFavorited);
    
    toast.add({
      severity: 'success',
      summary: '成功',
      detail: topic.value.userFavorited ? '已添加到收藏' : '已从收藏中移除',
      life: 3000
    })
  } catch (err) {
    const errorDetail = err.response?.data?.message || err.message
    
    toast.add({
      severity: 'error',
      summary: '操作失败',
      detail: errorDetail.includes('Unauthorized') 
        ? '请先登录再进行收藏操作'
        : errorDetail || '收藏操作失败',
      life: 3000
    })
    
    if (err.response?.status === 401) {
      router.push('/login')
    }
  }
}

// 点赞主题
async function likeTopic() {
  try {
    const wasLiked = topic.value.userLiked;
    
    // 立即更新UI状态
    topic.value.userLiked = !wasLiked;
    topic.value.likes = wasLiked ? topic.value.likes - 1 : topic.value.likes + 1;
    
    // 如果之前点踩过，需要取消点踩状态
    if (!wasLiked && topic.value.userDisliked) {
      topic.value.dislikes -= 1;
      topic.value.userDisliked = false;
    }
    
    const method = wasLiked ? 'DELETE' : 'POST';
    const response = await axios({
      method,
      url: `/api/forum/topics/${topic.value.id}/like`,
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    // 确保最终状态与API响应一致
    if (response.data.liked !== undefined) {
      topic.value.userLiked = response.data.liked;
      topic.value.likes = response.data.likes || topic.value.likes;
      
      // 如果点赞成功，确保取消点踩状态
      if (response.data.liked && topic.value.userDisliked) {
        topic.value.userDisliked = false;
        topic.value.dislikes -= 1;
      }
    }
    
    toast.add({
      severity: 'success',
      summary: '成功',
      detail: topic.value.userLiked ? '点赞成功' : '已取消点赞',
      life: 3000
    });
  } catch (err) {
    // 恢复原始状态
    topic.value.userLiked = !topic.value.userLiked;
    topic.value.likes = topic.value.userLiked ? topic.value.likes + 1 : topic.value.likes - 1;
    
    const errorDetail = err.response?.data?.message || err.message;
    
    toast.add({
      severity: 'error',
      summary: '操作失败',
      detail: errorDetail.includes('Unauthorized') 
        ? '请先登录再进行点赞操作'
        : errorDetail || '点赞操作失败',
      life: 3000
    });
    
    if (err.response?.status === 401) {
      router.push('/login');
    }
  }
}

// 点踩主题
async function dislikeTopic() {
  try {
    const wasDisliked = topic.value.userDisliked;
    
    // 立即更新UI状态
    topic.value.userDisliked = !wasDisliked;
    topic.value.dislikes = wasDisliked ? topic.value.dislikes - 1 : topic.value.dislikes + 1;
    
    // 如果之前点赞过，需要取消点赞状态
    if (!wasDisliked && topic.value.userLiked) {
      topic.value.likes -= 1;
      topic.value.userLiked = false;
    }
    
    const method = wasDisliked ? 'DELETE' : 'POST';
    const response = await axios({
      method,
      url: `/api/forum/topics/${topic.value.id}/dislike`,
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    // 确保最终状态与API响应一致
    if (response.data.disliked !== undefined) {
      topic.value.userDisliked = response.data.disliked;
      topic.value.dislikes = response.data.dislikes || topic.value.dislikes;
      
      // 如果点踩成功，确保取消点赞状态
      if (response.data.disliked && topic.value.userLiked) {
        topic.value.userLiked = false;
        topic.value.likes -= 1;
      }
    }
    
    toast.add({
      severity: 'success',
      summary: '成功',
      detail: topic.value.userDisliked ? '点踩成功' : '已取消点踩',
      life: 3000
    });
  } catch (err) {
    // 恢复原始状态
    topic.value.userDisliked = !topic.value.userDisliked;
    topic.value.dislikes = topic.value.userDisliked ? topic.value.dislikes + 1 : topic.value.dislikes - 1;
    
    const errorDetail = err.response?.data?.message || err.message;
    
    toast.add({
      severity: 'error',
      summary: '操作失败',
      detail: errorDetail.includes('Unauthorized') 
        ? '请先登录再进行点踩操作'
        : errorDetail || '点踩操作失败',
      life: 3000
    });
    
    if (err.response?.status === 401) {
      router.push('/login');
    }
  }
}

// 处理点赞错误
function handleLikeError(error, action) {
  console.error(`${action}失败:`, error)
  const errorDetail = error.response?.data?.message || error.message
  toast.add({
    severity: 'error',
    summary: `${action}失败`,
    detail: errorDetail,
    life: 3000
  })
}

// 提交回复
const submitReply = async () => {
  try {
    loading.value = true
    
    const response = await axios.post(`/api/forum/topics/${route.params.id}/replies`, {
      content: newReply.value,
      reply_to: replyingTo.value?.id || null
    }, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    toast.add({
      severity: 'success',
      summary: '回复成功',
      detail: '您的回复已提交',
      life: 3000
    })
    
    newReply.value = ''
    replyingTo.value = null
    await loadTopic(route.params.id)
  } catch (err) {
    console.error('提交回复失败:', err)
    
    const errorDetail = err.response?.data?.message || err.message
    
    toast.add({
      severity: 'error',
      summary: '回复失败',
      detail: `无法提交回复: ${errorDetail}`,
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

// 更新回复
const updateReply = async (replyId, content) => {
  try {
    loading.value = true
    
    await axios.put(`/api/forum/replies/${replyId}`, {
      content: content
    })
    
    toast.add({
      severity: 'success',
      summary: '更新成功',
      detail: '回复已更新',
      life: 3000
    })
    
    await loadTopic()
  } catch (err) {
    console.error('更新回复失败:', err)
    
    const errorDetail = err.response?.data?.message || err.message
    
    toast.add({
      severity: 'error',
      summary: '更新失败',
      detail: `无法更新回复: ${errorDetail}`,
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

// 删除回复
const deleteReply = async (replyId) => {
  try {
    loading.value = true
    
    await axios.delete(`/api/forum/replies/${replyId}`)
    
    toast.add({
      severity: 'success',
      summary: '删除成功',
      detail: '回复已删除',
      life: 3000
    })
    
    await loadTopic()
  } catch (err) {
    console.error('删除回复失败:', err)
    
    const errorDetail = err.response?.data?.message || err.message
    
    toast.add({
      severity: 'error',
      summary: '删除失败',
      detail: `无法删除回复: ${errorDetail}`,
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCategories()
  loadTopic()
})

// 加载分类数据
async function loadCategories() {
  try {
    const response = await axios.get('/api/forum/categories')
    // 确保API返回的数据结构包含id和name字段
    if (Array.isArray(response.data)) {
      categories.value = response.data.map(category => ({
        id: category.id,
        name: category.name
      }))
    } else {
      throw new Error('Invalid categories data format')
    }
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: '加载失败',
      detail: '无法加载分类数据: ' + err.message,
      life: 3000
    })
  }
}


</script>

<style scoped>
.reply-input {
  min-height: 10rem;
}

.reply-editor {
  min-height: 10rem;
  padding: 1rem;
  border-radius: 0.5rem;
  border: 1px solid #e5e7eb;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.reply-editor:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
  outline: none;
}

.like-button-container {
  display: flex;
  align-items: center;
  margin-right: 12px;
}

.like-button.liked,
.dislike-button.disliked {
  transform-origin: center;
}

.animate-like {
  animation: likeAnimation 0.5s ease;
}

.animate-dislike {
  animation: dislikeAnimation 0.5s ease;
}

.animate-count {
  animation: countAnimation 0.5s ease;
}

@keyframes likeAnimation {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.3);
  }
  100% {
    transform: scale(1);
  }
}

@keyframes dislikeAnimation {
  0% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.3);
  }
  100% {
    transform: scale(1);
  }
}

@keyframes countAnimation {
  0% {
    transform: scale(1);
    opacity: 0.7;
  }
  50% {
    transform: scale(1.2);
    opacity: 1;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}
</style>