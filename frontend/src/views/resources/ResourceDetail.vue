<template>
  <div v-if="isLoading" class="flex justify-center items-center h-screen">
    <ProgressSpinner />
  </div>
  <div v-else-if="error" class="flex flex-col justify-center items-center h-screen">
    <i class="pi pi-exclamation-triangle text-5xl text-red-500 mb-4"></i>
    <p class="text-xl text-gray-700">{{ error }}</p>
    <Button label="重试" class="mt-4" @click="$router.go(0)" />
  </div>
  <div v-else class="bg-white rounded-lg shadow-md overflow-hidden">
    <Toast />
    <!-- 资源详情头部 -->
    <div class="bg-blue-600 text-white p-6">
      <div class="container mx-auto">
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center">
          <div>
            <h1 class="text-2xl font-bold mb-2">{{ resource.title }}</h1>
            <div class="flex flex-wrap items-center gap-2 text-sm">
              <Tag v-if="resource.category" :value="resource.category.name" severity="info" />
              <span class="flex items-center"><i class="pi pi-eye mr-1"></i> {{ resource.download_count || 0 }} 次查看</span>
              <span class="flex items-center"><i class="pi pi-download mr-1"></i> {{ resource.download_count || 0 }} 次下载</span>
              <span class="flex items-center"><i class="pi pi-thumbs-up mr-1"></i> {{ resource.likes || 0 }} 人点赞</span>
            </div>
          </div>
          <div class="mt-4 md:mt-0 flex space-x-2">
            <Button label="下载" icon="pi pi-download" class="p-button-outlined p-button-light" @click="downloadResource" />
            <Button 
              :icon="likeStatus ? 'pi pi-thumbs-up-fill' : 'pi pi-thumbs-up'" 
              class="p-button-outlined p-button-light" 
              @click="toggleLike"
              v-tooltip.top="likeStatus ? '取消点赞' : '点赞'"
            />
            <Button 
              :icon="isFavorited ? 'pi pi-star-fill' : 'pi pi-star'" 
              :class="isFavorited ? 'p-button-outlined p-button-warning' : 'p-button-outlined p-button-light'"
              @click="toggleFavorite"
              v-tooltip.top="isFavorited ? '取消收藏' : '收藏'"
            />
            <Button 
              icon="pi pi-share-alt" 
              class="p-button-outlined p-button-light" 
              v-tooltip.top="'分享'" 
              @click="shareResource" 
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 资源详情内容 -->
    <div class="container mx-auto p-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- 左侧资源信息 -->
        <div class="md:col-span-2">
          <Card>
            <template #title>
              <div class="flex justify-between items-center">
                <span>资源详情</span>
                <Tag :value="resource.fileType" severity="success" />
              </div>
            </template>
            <template #content>
              <div class="prose max-w-none" v-html="resource.description"></div>
              
              <Divider />
              
              <h3 class="text-lg font-semibold mb-3">资源预览</h3>
              <div v-if="resource.file_path">
                <div class="border border-gray-200 rounded-lg overflow-hidden">
                  <img :src="'/storage/' + resource.file_path" :alt="resource.title" class="w-full h-auto" />
                </div>
              </div>
              <div v-else class="text-center py-8 bg-gray-50 rounded-lg">
                <i class="pi pi-image text-5xl text-gray-300 mb-4"></i>
                <p class="text-gray-500">暂无预览图</p>
              </div>
              
              <Divider />
              
              <h3 class="text-lg font-semibold mb-3">文件信息</h3>
              <ul class="space-y-2">
                <li class="flex justify-between">
                  <span class="text-gray-600">文件大小：</span>
                  <span class="font-medium">{{ formatFileSize(resource.file_size) }}</span>
                </li>
                <li class="flex justify-between">
                  <span class="text-gray-600">文件格式：</span>
                  <span class="font-medium">{{ formatFileType(resource.file_type) }}</span>
                </li>
                <li class="flex justify-between">
                  <span class="text-gray-600">上传时间：</span>
                  <span class="font-medium">{{ formatDate(resource.created_at) }}</span>
                </li>
                <li class="flex justify-between">
                  <span class="text-gray-600">最后更新：</span>
                  <span class="font-medium">{{ resource.updated_at ? formatDate(resource.updated_at) : '无更新记录' }}</span>
                </li>
              </ul>
            </template>
          </Card>
          
          <!-- 评论区 -->
          <Card class="mt-6">
            <template #title>
              <div class="flex justify-between items-center">
                <span>评论 ({{ resource.comments ? resource.comments.length : 0 }})</span>
                <Button label="发表评论" icon="pi pi-comment" class="p-button-sm" @click="showCommentDialog = true" />
              </div>
            </template>
            <template #content>
              <div v-if="resource.comments && resource.comments.length > 0">
                <div v-for="(comment, index) in resource.comments" :key="index" class="border-b border-gray-200 py-4 last:border-b-0">
                  <div class="flex">
                    <Avatar :image="comment.userAvatar" :label="!comment.userAvatar ? comment.username.charAt(0).toUpperCase() : undefined" shape="circle" class="mr-3" />
                    <div class="flex-1">
                      <div class="flex justify-between">
                        <div>
                          <span class="font-medium">{{ comment.username }}</span>
                          <span class="text-gray-500 text-sm ml-2">{{ formatDate(comment.time) }}</span>
                        </div>
                        <div class="flex items-center">
                          <Rating v-model="comment.rating" :readonly="true" :cancel="false" />
                          <Button 
                            v-if="comment.userId === currentUser?.id"
                            icon="pi pi-trash" 
                            class="p-button-text p-button-danger ml-2" 
                            @click="deleteComment(comment.id, index)"
                            v-tooltip.top="'删除评论'"
                          />
                        </div>
                      </div>
                      <p class="mt-2">{{ comment.content }}</p>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="text-center py-8">
                <i class="pi pi-comments text-5xl text-gray-300 mb-4"></i>
                <p class="text-gray-500">暂无评论，成为第一个评论者吧！</p>
              </div>
            </template>
          </Card>
        </div>
        
        <!-- 右侧信息 -->
        <div>
          <!-- 上传者信息 -->
          <Card>
            <template #title>上传者信息</template>
            <template #content>
              <div class="flex items-center mb-4">
                <Avatar :image="resource.user.avatar" :label="!resource.user.avatar ? resource.user.username.charAt(0).toUpperCase() : undefined" shape="circle" class="mr-3" />
                <div>
                  <div class="font-medium">{{ resource.user.username }}</div>
                  <div class="text-sm text-gray-500">{{ resource.user.role }}</div>
                </div>
              </div>
              <p class="text-sm text-gray-600 mb-4">{{ resource.user.email }}</p>
              <Button label="查看更多资源" icon="pi pi-search" class="p-button-outlined w-full" @click="viewUploaderResources" />
            </template>
          </Card>
          
          <!-- 相关资源 -->
          <Card class="mt-6">
            <template #title>相关资源</template>
            <template #content>
              <div v-if="relatedResources.length > 0">
                <div v-for="(item, index) in relatedResources" :key="index" class="mb-4 last:mb-0">
                  <div class="cursor-pointer hover:text-blue-600" @click="navigateToResource(item.id)">
                    <div class="font-medium">{{ item.title }}</div>
                    <div class="text-sm text-gray-500 mt-1 flex justify-between">
                      <span>{{ item.category }}</span>
                      <span>{{ item.downloads }} 下载</span>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="text-center py-4">
                <p class="text-gray-500 text-sm">暂无相关资源</p>
              </div>
            </template>
          </Card>
          
          <!-- 下载须知 -->
          <Card class="mt-6">
            <template #title>下载须知</template>
            <template #content>
              <ul class="text-sm space-y-2">
                <li class="flex items-start">
                  <i class="pi pi-info-circle mt-0.5 mr-2 text-blue-500"></i>
                  <span>下载需要消耗 <span class="font-bold text-blue-600">{{ resource.points_required }}</span> 积分</span>
                </li>
                <li class="flex items-start">
                  <i class="pi pi-info-circle mt-0.5 mr-2 text-blue-500"></i>
                  <span>资源仅供学习交流使用，请勿用于商业用途</span>
                </li>
                <li class="flex items-start">
                  <i class="pi pi-info-circle mt-0.5 mr-2 text-blue-500"></i>
                  <span>如有侵权问题，请联系管理员处理</span>
                </li>
              </ul>
              <Divider />
              <Button label="立即下载" icon="pi pi-download" class="w-full" @click="downloadResource" />
            </template>
          </Card>
        </div>
      </div>
    </div>
    
    <!-- 评论对话框 -->
    <Dialog v-model:visible="showCommentDialog" header="发表评论" :style="{width: '500px'}" :modal="true">
      <div class="p-fluid">
        <div class="field mb-4">
          <label for="rating" class="block text-sm font-medium text-gray-700 mb-1">评分</label>
          <Rating v-model="newComment.rating" :cancel="false" />
        </div>
        <div class="field">
          <label for="comment" class="block text-sm font-medium text-gray-700 mb-1">评论内容</label>
          <Textarea v-model="newComment.content" rows="5" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" icon="pi pi-times" class="p-button-text" @click="showCommentDialog = false" />
        <Button label="提交" icon="pi pi-check" @click="submitComment" />
      </template>
    </Dialog>
    
    <!-- 分享对话框 -->
    <Dialog v-model:visible="showShareDialog" header="分享资源" :style="{width: '500px'}" :modal="true">
      <div class="p-fluid">
        <div class="field mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">资源链接</label>
          <div class="p-inputgroup">
            <InputText v-model="shareLink" readonly />
            <Button icon="pi pi-copy" @click="copyShareLink" />
          </div>
        </div>
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700 mb-3">分享到社交媒体</label>
          <div class="flex space-x-4">
            <Button icon="pi pi-facebook" class="p-button-rounded p-button-primary" />
            <Button icon="pi pi-twitter" class="p-button-rounded p-button-info" />
            <Button icon="pi pi-linkedin" class="p-button-rounded p-button-help" />
            <Button icon="pi pi-whatsapp" class="p-button-rounded p-button-success" />
          </div>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import Button from 'primevue/button'
import Card from 'primevue/card'
import Avatar from 'primevue/avatar'
import Tag from 'primevue/tag'
import Divider from 'primevue/divider'
import Carousel from 'primevue/carousel'
import Rating from 'primevue/rating'
import Dialog from 'primevue/dialog'
import Textarea from 'primevue/textarea'
import InputText from 'primevue/inputtext'
import ProgressSpinner from 'primevue/progressspinner'
import Toast from 'primevue/toast'
import Tooltip from 'primevue/tooltip'
import { useToast } from 'primevue/usetoast'
const toast = useToast()

const route = useRoute()
const router = useRouter()
const showCommentDialog = ref(false)
const showShareDialog = ref(false)
const shareLink = ref('')
// 新评论
const newComment = ref({
  rating: 5,
  content: ''
})

// 资源数据
const resource = ref({})
const currentUser = JSON.parse(localStorage.getItem('user'))

// 加载状态
const isLoading = ref(true)
const error = ref(null)

// 点赞状态
const likeStatus = ref(false)

// 初始化点赞状态
const initLikeStatus = async () => {
  try {
    const response = await axios.get(`/api/resources/${route.params.id}/like-status`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    likeStatus.value = response.data.isLiked
  } catch (error) {
    console.error('获取点赞状态失败:', error)
    likeStatus.value = false
    resource.value.likes = 0
  }
}

// 相关资源
const relatedResources = ref([
  {
    id: 2,
    title: 'Vue 3 + TypeScript实战指南',
    category: '前端开发',
    downloads: 325
  },
  {
    id: 3,
    title: 'Pinia状态管理详解',
    category: '前端开发',
    downloads: 189
  },
  {
    id: 4,
    title: 'Vue 3性能优化最佳实践',
    category: '前端开发',
    downloads: 276
  }
])

// 格式化日期
function formatDate(date) {
  if (!date) return '未知日期'
  
  try {
    const parsedDate = new Date(date)
    if (isNaN(parsedDate.getTime())) return '无效日期'
    
    return new Intl.DateTimeFormat('zh-CN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    }).format(parsedDate)
  } catch (e) {
    console.error('日期格式化失败:', e)
    return '日期错误'
  }
}

// 格式化文件大小
function formatFileSize(bytes) {
  if (!bytes) return '未知大小';
  
  const units = ['B', 'KB', 'MB', 'GB'];
  let size = Number(bytes);
  let unitIndex = 0;
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }
  
  return `${size.toFixed(1)} ${units[unitIndex]}`;
}

// 格式化文件类型
function formatFileType(mimeType) {
  const typeMap = {
    'application/msword': 'Word文档',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document': 'Word文档',
    'application/vnd.ms-excel': 'Excel表格',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': 'Excel表格',
    'application/vnd.ms-powerpoint': 'PPT演示',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation': 'PPT演示',
    'application/pdf': 'PDF文档',
    'text/plain': '文本文件',
    'application/zip': '压缩文件',
    'application/x-rar-compressed': '压缩文件'
  };
  
  return typeMap[mimeType] || mimeType || '未知类型';
}

// 下载资源
async function downloadResource() {
  try {
    const response = await axios.get(`/api/download/${route.params.id}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      baseURL: 'http://localhost:8080'
    })
    
    if (response.data.url) {
      // 创建隐藏的a标签触发下载
      const link = document.createElement('a')
      link.href = response.data.url
      link.download = resource.value.title || 'download'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      
      // 更新下载次数
      resource.value.download_count += 1
    }
  } catch (error) {
    console.error('下载失败:', error)
    toast.add({
      severity: 'error',
      summary: '错误',
      detail: '下载失败: ' + (error.response?.data?.message || error.message),
      life: 5000
    })
  }
}

// 点赞/取消点赞
async function toggleLike() {
  try {
    let response;
    const wasLiked = likeStatus.value;
    console.log('当前点赞状态:', wasLiked);
    
    // 立即更新UI状态
    likeStatus.value = !wasLiked;
    resource.value.likes = wasLiked ? resource.value.likes - 1 : resource.value.likes + 1;
    
    try {
      if (wasLiked) {
        response = await axios.delete(`/api/resources/${resource.value.id}/dislike`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'application/json'
          },
          baseURL: 'http://localhost:8080'
        });
      } else {
        response = await axios.post(`/api/resources/${resource.value.id}/like`, {}, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
            'Content-Type': 'application/json'
          },
          baseURL: 'http://localhost:8080'
        });
      }
      
      console.log('API点赞响应原始数据:', response.data);
      
      // 确保最终状态与API响应一致，兼容多种可能的响应格式
      if (response.data.isLiked !== undefined) {
        likeStatus.value = response.data.isLiked;
      } else if (response.data.is_liked !== undefined) {
        likeStatus.value = response.data.is_liked;
      } else if (typeof response.data === 'boolean') {
        likeStatus.value = response.data;
      } else {
        console.warn('无法从API响应中确定点赞状态，保持当前UI状态');
      }
      
      console.log('更新后的点赞状态:', likeStatus.value);
      
      toast.add({
        severity: 'success',
        summary: '成功',
        detail: likeStatus.value ? '已点赞' : '已取消点赞',
        life: 3000
      });
    } catch (error) {
      console.error('点赞操作失败:', error);
      // 出错时恢复原状态
      likeStatus.value = wasLiked;
      
      toast.add({
        severity: 'error',
        summary: '错误',
        detail: '点赞操作失败: ' + (error.response?.data?.error || error.response?.data?.message || error.message),
        life: 5000
      });
    }
  } catch (error) {
    console.error('点赞操作失败:', error);
    // 出错时恢复原状态
    likeStatus.value = !likeStatus.value;
    
    toast.add({
      severity: 'error',
      summary: '错误',
      detail: '点赞操作失败: ' + (error.response?.data?.error || error.response?.data?.message || error.message),
      life: 5000
    });
  }
}

// 收藏状态
const isFavorited = ref(false);

// 收藏/取消收藏
async function toggleFavorite() {
  try {
    let response;
    const wasFavorited = isFavorited.value;
    console.log('当前收藏状态:', wasFavorited);
    
    // 立即更新UI状态
    isFavorited.value = !wasFavorited;
    
    if (wasFavorited) {
      response = await axios.delete(`/api/resources/${resource.value.id}/favorite`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
    } else {
      response = await axios.post(`/api/resources/${resource.value.id}/favorite`, {}, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
    }
    
    console.log('API收藏响应原始数据:', response.data);
    
    // 确保最终状态与API响应一致，兼容多种可能的响应格式
    if (response.data.isFavorited !== undefined) {
      isFavorited.value = response.data.isFavorited;
    } else if (response.data.is_favorite !== undefined) {
      isFavorited.value = response.data.is_favorite;
    } else if (typeof response.data === 'boolean') {
      isFavorited.value = response.data;
    } else {
      console.warn('无法从API响应中确定收藏状态，保持当前UI状态');
    }
    
    console.log('更新后的收藏状态:', isFavorited.value);
    
    toast.add({
      severity: 'success',
      summary: '成功',
      detail: isFavorited.value ? '已添加到收藏' : '已从收藏中移除',
      life: 3000
    });
  } catch (error) {
    console.error('收藏操作失败:', error);
    // 出错时恢复原状态
    isFavorited.value = !isFavorited.value;
    toast.add({
      severity: 'error',
      summary: '错误',
      detail: '收藏操作失败: ' + (error.response?.data?.error || error.response?.data?.message || error.message),
      life: 5000
    });
  }
}

// 分享资源
function shareResource() {
  shareLink.value = `https://localhost:5173/resources/${resource.value.id}`
  showShareDialog.value = true
}

// 复制分享链接
function copyShareLink() {
  navigator.clipboard.writeText(shareLink.value)
    .then(() => {
      alert('链接已复制到剪贴板')
    })
    .catch(err => {
      console.error('复制失败:', err)
    })
}

// 查看上传者的更多资源
function viewUploaderResources() {
  // 实际应用中应导航到上传者的资源列表页
  router.push(`/resources?uploader=${resource.value.uploader.username}`)
}

// 导航到相关资源
function navigateToResource(resourceId) {
  router.push(`/resources/${resourceId}`)
}

// 提交评论
async function deleteComment(commentId, index) {
  try {
    // 确保使用正确的评论ID格式
    await axios.delete(`/api/resources/${resource.value.id}/comments/${commentId}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    resource.value.comments.splice(index, 1)
    
    toast.add({
      severity: 'success',
      summary: '成功',
      detail: '评论已删除',
      life: 3000
    })
  } catch (error) {
    console.error('删除评论失败:', error)
    toast.add({
      severity: 'error',
      summary: '错误',
      detail: '删除评论失败: ' + (error.response?.data?.message || error.message),
      life: 5000
    })
  }
}

async function submitComment() {
  if (!newComment.value.content.trim()) {
    alert('评论内容不能为空')
    return
  }
  
  try {
    const response = await axios.post(`/api/resources/${resource.value.id}/comments`, {
      rating: newComment.value.rating,
      content: newComment.value.content
    }, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    // 更新前端评论列表
    if (!resource.value.comments) {
      resource.value.comments = []
    }
    
    // 处理后端返回的评论数据
    const comment = response.data.comment || response.data;
    resource.value.comments.unshift({
      id: comment.ID || comment.id,
      content: comment.Content || comment.content,
      rating: comment.Rating || comment.rating,
      userId: comment.UserID || comment.user_id,
      userAvatar: comment.User?.avatar || comment.user?.avatar || currentUser?.avatar,
      username: comment.User?.username || comment.user?.username || currentUser?.username,
      time: comment.CreatedAt || comment.Time || comment.created_at
    })
    
    // 重置表单并关闭对话框
    newComment.value = {
      rating: 5,
      content: ''
    }
    showCommentDialog.value = false
    
    // 显示成功提示
    toast.add({
      severity: 'success',
      summary: '成功',
      detail: '评论已提交',
      life: 3000
    })
  } catch (error) {
    console.error('提交评论失败:', error)
    toast.add({
      severity: 'error',
      summary: '错误',
      detail: '提交评论失败: ' + (error.response?.data?.message || error.message),
      life: 5000
    })
  }
}

onMounted(async () => {
  const resourceId = route.params.id
  if (!resourceId) {
    error.value = '无效的资源ID'
    isLoading.value = false
    return
  }
  
  isLoading.value = true
  error.value = null
  
  try {
    // 获取资源详情、评论和收藏状态
    const [resourceResponse, commentsResponse, favoriteResponse] = await Promise.all([
      axios.get(`/api/resources/${resourceId}`),
      axios.get(`/api/resources/${resourceId}/comments`),
      axios.get(`/api/resources/${resourceId}/favorite-status`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      }).catch(() => ({ data: { isFavorited: false } }))
    ])
    
    // 初始化点赞状态
    await initLikeStatus()
    
    console.log('commentsResponse.data:', commentsResponse.data)
    console.log('点赞状态:', likeStatus.data)
    
    if (resourceResponse.data) {
      resource.value = resourceResponse.data
      if (commentsResponse.data && commentsResponse.data.comments && Array.isArray(commentsResponse.data.comments)) {
        console.log('原始评论数据:', commentsResponse.data.comments)
        resource.value.comments = commentsResponse.data.comments.map(comment => ({
          id: comment.ID,
          content: comment.Content,
          rating: comment.Rating,
          userId: comment.UserID,
          userAvatar: comment.User?.avatar,
          username: comment.User?.username || '匿名用户',
          time: comment.CreatedAt || comment.Time
        }))
        console.log('处理后评论数据:', resource.value.comments)
      }

    } else {
      error.value = '资源不存在'
      resource.value = {}
    }


  // 设置收藏状态
  console.log('收藏状态原始响应:', favoriteResponse?.data);
  if (favoriteResponse.data.isFavorited !== undefined) {
    isFavorited.value = favoriteResponse.data.isFavorited;
  } else if (favoriteResponse.data.is_favorite !== undefined) {
    isFavorited.value = favoriteResponse.data.is_favorite;
  } else if (typeof favoriteResponse.data === 'boolean') {
    isFavorited.value = favoriteResponse.data;
  } else {
    isFavorited.value = false;
    console.warn('无法从API响应中确定初始收藏状态，默认设置为未收藏');
  }
  console.log('初始化收藏状态:', isFavorited.value);

  } catch (err) {
    error.value = err.message || '加载资源详情失败'
    console.error('加载失败:', err)
    resource.value = {}
  } finally {
    isLoading.value = false
  }
})
</script>