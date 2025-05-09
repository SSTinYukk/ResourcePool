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
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[5, 10, 20, 50]"
      responsiveLayout="scroll"
      class="p-datatable-sm">
      
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

const router = useRouter()
const searchQuery = ref('')
const selectedCategory = ref(null)

// 模拟数据 - 实际应用中应从API获取
const categories = [
  { id: 1, name: '全部' },
  { id: 2, name: '技术讨论' },
  { id: 3, name: '学习资源' },
  { id: 4, name: '经验分享' },
  { id: 5, name: '求助问答' },
  { id: 6, name: '活动公告' }
]

const topics = [
  {
    id: 1,
    title: '如何高效学习Vue 3？',
    brief: '分享一些学习Vue 3的经验和资源推荐',
    author: '张三',
    authorAvatar: null,
    category: 3,
    replies: 24,
    views: 356,
    lastReplyTime: new Date(2023, 8, 15, 14, 30),
    lastReplyAuthor: '李四'
  },
  {
    id: 2,
    title: 'Tailwind CSS使用技巧分享',
    brief: '整理了一些Tailwind CSS的常用技巧和最佳实践',
    author: '王五',
    authorAvatar: null,
    category: 4,
    replies: 18,
    views: 245,
    lastReplyTime: new Date(2023, 8, 16, 9, 15),
    lastReplyAuthor: '赵六'
  },
  {
    id: 3,
    title: '求助：Go语言并发编程问题',
    brief: '在使用goroutine时遇到了一些问题，求大神指点',
    author: '李四',
    authorAvatar: null,
    category: 5,
    replies: 12,
    views: 189,
    lastReplyTime: new Date(2023, 8, 16, 16, 45),
    lastReplyAuthor: '张三'
  },
  {
    id: 4,
    title: '线上技术分享会：微服务架构实践',
    brief: '下周三晚8点，欢迎参加我们的线上技术分享会',
    author: '管理员',
    authorAvatar: null,
    category: 6,
    replies: 5,
    views: 120,
    lastReplyTime: new Date(2023, 8, 14, 11, 20),
    lastReplyAuthor: '王五'
  },
  {
    id: 5,
    title: 'Laravel与Vue.js前后端分离实践',
    brief: '分享一个完整的前后端分离项目的开发经验',
    author: '赵六',
    authorAvatar: null,
    category: 4,
    replies: 30,
    views: 412,
    lastReplyTime: new Date(2023, 8, 17, 10, 5),
    lastReplyAuthor: '张三'
  }
]

// 根据搜索和分类过滤主题
const filteredTopics = computed(() => {
  let result = [...topics]
  
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
  router.push('/forum/create')
}
</script>