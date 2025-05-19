<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">论坛交流</h1>
      <Button label="发布新主题" icon="pi pi-plus" class="p-button-primary" @click="navigateToCreate" />
    </div>

    <div class="mb-4">
      <span class="p-input-icon-left w-full md:w-96">
        <i class="pi pi-search" />
        <InputText v-model="searchQuery" placeholder="搜索主题..." class="w-full" />
      </span>
    </div>

    <!-- 分类标签 -->
    <div class="mb-6 flex flex-wrap gap-2">
      <Button v-for="category in categories" :key="category.id"
        :label="category.name"
        :class="['p-button-rounded p-button-sm', selectedCategory === category.id ? 'p-button-primary' : 'p-button-outlined']"
        @click="toggleCategory(category.id)" />
    </div>

    <!-- 主题列表 -->
    <DataTable :value="filteredTopics" :paginator="true" :rows="10"
      :loading="loading"
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[5, 10, 20, 50]"
      responsiveLayout="scroll"
      class="p-datatable-sm">
      
      <template #empty>
        <div v-if="loading" class="text-center py-8">
          <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="4" />
          <p class="mt-4 text-gray-500">加载中...</p>
        </div>
        <div v-else-if="error" class="text-center py-8">
          <i class="pi pi-exclamation-triangle text-5xl text-red-500 mb-4"></i>
          <p class="text-red-500">{{ error }}</p>
          <Button label="重试" icon="pi pi-refresh" class="p-button-text mt-4" @click="fetchTopics" />
        </div>
      </template>
      
      <Column field="title" header="主题" style="min-width: 50%">
        <template #body="{data}">
          <div class="cursor-pointer" @click="navigateToTopic(data.id)">
            <div class="font-medium text-blue-600 hover:text-blue-800">{{ data.title }}</div>
            <div class="text-sm text-gray-500 mt-1">{{ data.brief }}</div>
          </div>
        </template>
      </Column>
      
      <Column field="author" header="作者" style="min-width: 15%">
        <template #body="{data}">
          <div class="flex items-center">
            <Avatar :image="data.authorAvatar" :label="!data.authorAvatar ? data.author.charAt(0).toUpperCase() : undefined" shape="circle" size="small" class="mr-2" />
            <span>{{ data.author }}</span>
          </div>
        </template>
      </Column>
      
      <Column field="replies" header="回复/查看" style="min-width: 15%">
        <template #body="{data}">
          <div>
            <span class="text-blue-600 font-medium">{{ data.replies }}</span> / {{ data.views }}
          </div>
        </template>
      </Column>
      
      <Column field="lastReply" header="最后回复" style="min-width: 20%">
        <template #body="{data}">
          <div>
            <div class="text-sm">{{ formatDate(data.lastReplyTime) }}</div>
            <div class="text-xs text-gray-500">by {{ data.lastReplyAuthor }}</div>
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- 无数据时显示 -->
    <div v-if="filteredTopics.length === 0" class="text-center py-8">
      <i class="pi pi-comments text-5xl text-gray-300 mb-4"></i>
      <p class="text-gray-500">暂无主题，成为第一个发布者吧！</p>
      <Button label="发布新主题" icon="pi pi-plus" class="p-button-primary mt-4" @click="navigateToCreate" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Avatar from 'primevue/avatar'
import ProgressSpinner from 'primevue/progressspinner'

const router = useRouter()
const searchQuery = ref('')
const selectedCategory = ref(null)

const loading = ref(false)
const error = ref(null)
const categories = ref([])
const topics = ref([])

// 获取分类数据
const fetchCategories = async () => {
  try {
    const response = await axios.get('/api/forum/categories')
    categories.value = response.data
  } catch (err) {
    error.value = '获取分类数据失败'
    console.error(err)
  }
}

// 获取主题列表
const fetchTopics = async () => {
  loading.value = true
  error.value = null
  try {
    console.log('正在请求主题列表API...')
    const response = await axios.get('/api/forum/topics')
    console.log('主题列表API响应:', response)
    topics.value = Array.isArray(response.data.topics) ? response.data.topics.map(topic => ({
      id: topic.id,
      title: topic.title,
      brief: topic.content.substring(0, 50),
      author: topic.user?.username || '匿名',
      authorAvatar: topic.user?.avatar,
      replies: 0,
      views: 0,
      lastReplyTime: new Date(),
      lastReplyAuthor: '',
      category: 1
    })) : []
  } catch (err) {
    error.value = `获取主题列表失败: ${err.message}`
    console.error('获取主题列表API错误:', {
      message: err.message,
      response: err.response,
      config: err.config
    })
  } finally {
    loading.value = false
  }
}

// 初始化数据
onMounted(() => {
  fetchCategories()
  fetchTopics()
})

// 根据搜索和分类过滤主题
// 添加缺失的导入
import axios from 'axios'
import { onMounted } from 'vue'

// 修改filteredTopics计算属性中的topics引用
const filteredTopics = computed(() => {
  let result = [...topics.value]  // 修改为topics.value
  
  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(topic => 
      topic.title.toLowerCase().includes(query) || 
      topic.brief.toLowerCase().includes(query)
    )
  }
  
  // 分类过滤
  if (selectedCategory.value && selectedCategory.value !== 1) {
    result = result.filter(topic => topic.category === selectedCategory.value)
  }
  
  return result
})

// 切换分类
function toggleCategory(categoryId) {
  selectedCategory.value = selectedCategory.value === categoryId ? null : categoryId
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

// 导航到主题详情
function navigateToTopic(topicId) {
  router.push(`/forum/topic/${topicId}`)
}

// 导航到创建主题
function navigateToCreate() {
  router.push('/forum/create').catch(err => {
    console.error('导航失败:', err)
    // 可以添加用户提示
  })
}
</script>