<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold mb-6">我的收藏</h1>
    
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <!-- 侧边栏导航 -->
      <div class="bg-white rounded-lg shadow p-4">
        <h2 class="text-lg font-semibold mb-4">收藏分类</h2>
        <ul class="space-y-2">
          <li>
            <a 
              href="#" 
              @click.prevent="activeTab = 'resources'" 
              class="block px-3 py-2 rounded-md transition-colors" 
              :class="activeTab === 'resources' ? 'bg-blue-50 text-blue-600' : 'text-gray-700 hover:bg-gray-50'"
            >
              <i class="pi pi-file mr-2"></i>收藏的资源
            </a>
          </li>
          <li>
            <a 
              href="#" 
              @click.prevent="activeTab = 'posts'" 
              class="block px-3 py-2 rounded-md transition-colors" 
              :class="activeTab === 'posts' ? 'bg-blue-50 text-blue-600' : 'text-gray-700 hover:bg-gray-50'"
            >
              <i class="pi pi-comments mr-2"></i>收藏的帖子
            </a>
          </li>
        </ul>
      </div>
      
      <!-- 内容区域 -->
      <div class="md:col-span-3">
        <div class="bg-white rounded-lg shadow p-6">
          <!-- 资源收藏标签页 -->
          <div v-if="activeTab === 'resources'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-xl font-semibold">收藏的资源</h2>
              <div class="flex items-center">
                <span class="text-sm text-gray-500 mr-2">排序方式:</span>
                <Dropdown v-model="resourceSortOption" :options="sortOptions" optionLabel="label" class="w-32" />
              </div>
            </div>
            
            <div v-if="loading" class="flex justify-center py-8">
              <ProgressSpinner style="width:50px;height:50px" strokeWidth="4" />
            </div>
            
            <div v-else-if="!favoriteResources || favoriteResources.length === 0" class="text-center py-8 text-gray-500">
              <i class="pi pi-inbox text-4xl mb-4"></i>
              <p>您还没有收藏任何资源</p>
              <router-link to="/resources" class="text-blue-600 hover:underline mt-2 inline-block">
                去浏览资源
              </router-link>
            </div>
            
            <div v-else class="space-y-4">
              <div v-for="resource in favoriteResources" :key="resource.id" class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
                <div class="flex justify-between">
                  <div>
                    <router-link :to="`/resources/${resource.id}`" class="text-lg font-medium text-blue-600 hover:underline">
                      {{ resource.title }}
                    </router-link>
                    <p class="text-sm text-gray-600 mt-1" v-html="resource.description"></p>
                    <div class="flex items-center mt-2 text-sm text-gray-500">
                      <span class="flex items-center mr-4">
                        <i class="pi pi-user mr-1"></i>
                        {{ resource.user.username }}
                      </span>
                      <span class="flex items-center mr-4">
                        <i class="pi pi-calendar mr-1"></i>
                        {{ formatDate(resource.created_at) }}
                      </span>
                    </div>
                  </div>
                  <div>
                    <Button icon="pi pi-trash" class="p-button-rounded p-button-text p-button-danger" @click="confirmRemoveFavorite(resource.id, 'resource')" />
                  </div>
                </div>
              </div>
              
              <Paginator 
                v-if="totalResourcePages > 1" 
                :rows="10" 
                :totalRecords="totalResources" 
                v-model:first="resourceFirst" 
                @page="onResourcePageChange($event)"
              />
            </div>
          </div>
          
          <!-- 帖子收藏标签页 -->
          <div v-if="activeTab === 'posts'">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-xl font-semibold">收藏的帖子</h2>
              <div class="flex items-center">
                <span class="text-sm text-gray-500 mr-2">排序方式:</span>
                <Dropdown v-model="postSortOption" :options="sortOptions" optionLabel="label" class="w-32" />
              </div>
            </div>
            
            <div v-if="loading" class="flex justify-center py-8">
              <ProgressSpinner style="width:50px;height:50px" strokeWidth="4" />
            </div>
            
            <div v-else-if="favoritePosts.length === 0" class="text-center py-8 text-gray-500">
              <i class="pi pi-inbox text-4xl mb-4"></i>
              <p>您还没有收藏任何帖子</p>
              <router-link to="/forum" class="text-blue-600 hover:underline mt-2 inline-block">
                去浏览论坛
              </router-link>
            </div>
            
            <div v-else class="space-y-4">
              <div v-for="post in favoritePosts" :key="post.id" class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
                <div class="flex justify-between">
                  <div>
                    <router-link :to="`/forum/topic/${post.id}`" class="text-lg font-medium text-blue-600 hover:underline">
                      {{ post.title }}
                    </router-link>
                    <p class="text-sm text-gray-600 mt-1" v-html="post.content"></p>
                    <div class="flex items-center mt-2 text-sm text-gray-500">
                      <span class="flex items-center mr-4">
                        <i class="pi pi-user mr-1"></i>
                        {{ post.user.username }}
                      </span>
                      <span class="flex items-center mr-4">
                        <i class="pi pi-calendar mr-1"></i>
                        {{ formatDate(post.created_at) }}
                      </span>
                      <span class="flex items-center mr-4">
                        <i class="pi pi-eye mr-1"></i>
                        {{ post.view_count }} 次查看
                      </span>
                      <span class="flex items-center">
                        <i class="pi pi-comments mr-1"></i>
                        {{ post.reply_count }} 条评论
                      </span>
                    </div>
                  </div>
                  <div>
                    <Button icon="pi pi-trash" class="p-button-rounded p-button-text p-button-danger" @click="confirmRemoveFavorite(post.id, 'post')" />
                  </div>
                </div>
              </div>
              
              <Paginator 
                v-if="totalPostPages > 1" 
                :rows="10" 
                :totalRecords="totalPosts" 
                v-model:first="postFirst" 
                @page="onPostPageChange($event)"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 确认对话框 -->
    <Dialog v-model:visible="confirmDialog" header="确认移除" :style="{width: '350px'}" :modal="true">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span>确定要移除这个收藏吗？</span>
      </div>
      <template #footer>
        <Button label="取消" icon="pi pi-times" class="p-button-text" @click="confirmDialog = false" />
        <Button label="确定" icon="pi pi-check" class="p-button-danger" @click="removeFavorite" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import Paginator from 'primevue/paginator'
import ProgressSpinner from 'primevue/progressspinner'
import Dialog from 'primevue/dialog'
import axios from 'axios'

const userStore = useUserStore()
const toast = useToast()

// 标签页状态
const activeTab = ref('resources')

// 加载状态
const loading = ref(false)

// 排序选项
const sortOptions = [
  { label: '最新收藏', value: 'newest' },
  { label: '最早收藏', value: 'oldest' },
  { label: '名称排序', value: 'name' }
]

// 资源收藏相关状态
const favoriteResources = ref([])
const resourceSortOption = ref(sortOptions[0])
const totalResources = ref(0)
const totalResourcePages = ref(0)
const resourceFirst = ref(0)

// 确保数组不为undefined
watch(favoriteResources, (newVal) => {
  if (!newVal) favoriteResources.value = []
}, { immediate: true })

// 计算当前资源页码
const currentResourcePage = computed(() => {
  return resourceFirst.value / 10 + 1
})

// 帖子收藏相关状态
const favoritePosts = ref([])
const postSortOption = ref(sortOptions[0])
const totalPosts = ref(0)
const totalPostPages = ref(0)
const postFirst = ref(0)

// 确保数组不为undefined
watch(favoritePosts, (newVal) => {
  if (!newVal) favoritePosts.value = []
}, { immediate: true })

// 确认对话框状态
const confirmDialog = ref(false)
const itemToRemove = ref(null)
const itemType = ref('')

// 格式化日期
function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

// 加载收藏的资源
async function loadFavoriteResources() {
  loading.value = true
  try {
    const response = await axios.get('/api/user/favorites/resources', {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      },
      params: {
        page: currentResourcePage.value,
        sort: resourceSortOption.value.value
      }
    })
    
    favoriteResources.value = response.data.data
    console.log('获取到的资源:', response.data.data)
    console.log("赋值后的资源:", favoriteResources.value)
    totalResources.value = response.data.data.total
    totalResourcePages.value = Math.ceil(response.data.data.total / 10)
  } catch (error) {
    console.error('Failed to load favorite resources:', error)
    toast.add({ severity: 'error', summary: '加载失败', detail: '无法加载收藏的资源', life: 3000 })
  } finally {
    loading.value = false
  }
}

// 加载收藏的帖子
async function loadFavoritePosts() {
  loading.value = true
  try {
    const response = await axios.get('/api/user/favorites/posts', {
      headers: {
        'Authorization': `Bearer ${userStore.token}`
      },
      params: {
        page: postFirst.value / 10 + 1,
      }
    })
    
    favoritePosts.value = response.data.data
    console.log('赋值后的帖子:', favoritePosts.value)
    totalPosts.value = response.data.data.total
    totalPostPages.value = Math.ceil(response.data.data.total / 10)
  } catch (error) {
    console.error('Failed to load favorite posts:', error)
    toast.add({ severity: 'error', summary: '加载失败', detail: '无法加载收藏的帖子', life: 3000 })
  } finally {
    loading.value = false
  }
}

// 资源分页变化
function onResourcePageChange(event) {
  resourceFirst.value = event.first
  loadFavoriteResources()
}

// 帖子分页变化
function onPostPageChange(event) {
  postFirst.value = event.first
  loadFavoritePosts()
}

// 确认移除收藏
function confirmRemoveFavorite(id, type) {
  itemToRemove.value = id
  itemType.value = type
  confirmDialog.value = true
}

// 移除收藏
async function removeFavorite() {
  try {
    const endpoint = itemType.value === 'resource' 
      ? `/api/user/favorites/resources/${itemToRemove.value}` 
      : `/api/user/favorites/posts/${itemToRemove.value}`
    
    await axios.delete(endpoint, {
      headers: {
        'Authorization': `Bearer ${userStore.token}`
      }
    })
    
    toast.add({ severity: 'success', summary: '移除成功', detail: '已从收藏中移除', life: 3000 })
    
    // 刷新列表
    if (itemType.value === 'resource') {
      loadFavoriteResources()
    } else {
      loadFavoritePosts()
    }
  } catch (error) {
    console.error('Failed to remove favorite:', error)
    toast.add({ severity: 'error', summary: '操作失败', detail: '无法移除收藏', life: 3000 })
  } finally {
    confirmDialog.value = false
  }
}

// 监听排序选项变化
watch(resourceSortOption, () => {
  resourceFirst.value = 0
  loadFavoriteResources()
})

watch(postSortOption, () => {
  postFirst.value = 0
  loadFavoritePosts()
})

// 监听标签页变化
watch(activeTab, (newTab) => {
  if (newTab === 'resources') {
    loadFavoriteResources()
  } else if (newTab === 'posts') {
    loadFavoritePosts()
  }
})

onMounted(() => {
  // 初始加载资源收藏
  loadFavoriteResources()
})
</script>