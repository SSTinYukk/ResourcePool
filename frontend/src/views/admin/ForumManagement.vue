<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">论坛管理</h1>
    </div>

    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <!-- 搜索和筛选 -->
      <div class="flex flex-col md:flex-row gap-4 mb-4">
        <div class="flex-1">
          <span class="p-input-icon-left w-full">
            <i class="pi pi-search" />
            <InputText v-model="filters.search" placeholder="搜索话题标题或内容" class="w-full" @keyup.enter="loadTopics()" />
          </span>
        </div>
        <div class="flex gap-2">
          <Dropdown v-model="filters.category" :options="categoryOptions" optionLabel="name" optionValue="id" placeholder="分类" class="w-32" />
          <Button label="搜索" icon="pi pi-search" @click="loadTopics()" />
          <Button label="重置" icon="pi pi-refresh" outlined @click="resetFilters()" />
        </div>
      </div>

      <!-- 话题表格 -->
      <DataTable 
        :value="topics" 
        :loading="loading"
        :paginator="true" 
        :rows="pagination.pageSize"
        :totalRecords="pagination.total"
        :rowsPerPageOptions="[10, 20, 50]"
        v-model:first="pagination.first"
        v-model:rows="pagination.pageSize"
        @page="onPageChange($event)"
        dataKey="id"
        stripedRows
        responsiveLayout="scroll"
        class="p-datatable-sm"
      >
        <Column field="id" header="ID" :sortable="true" style="width: 80px"></Column>
        <Column field="title" header="标题" :sortable="true">
          <template #body="{data}">
            <div class="flex items-center">
              <span>{{ data.title }}</span>
            </div>
          </template>
        </Column>
        <Column field="user.username" header="发布者" :sortable="true">
          <template #body="{data}">
            <div class="flex items-center">
              <Avatar :image="data.user?.avatar || '/default-avatar.png'" shape="circle" class="mr-2" size="small" />
              <span>{{ data.user?.username }}</span>
            </div>
          </template>
        </Column>
        <Column field="category.name" header="分类" :sortable="true">
          <template #body="{data}">
            <Tag :value="data.category?.name || '未分类'" />
          </template>
        </Column>
        <Column field="views" header="浏览量" :sortable="true">
          <template #body="{data}">
            <span>{{ data.views || 0 }}</span>
          </template>
        </Column>
        <Column field="replies_count" header="回复数" :sortable="true">
          <template #body="{data}">
            <span>{{ data.replies_count || 0 }}</span>
          </template>
        </Column>
        <Column field="created_at" header="发布时间" :sortable="true">
          <template #body="{data}">
            {{ formatDate(data.created_at) }}
          </template>
        </Column>
        <Column header="操作" style="width: 150px">
          <template #body="{data}">
            <div class="flex gap-1">
              <Button icon="pi pi-eye" text rounded @click="viewTopic(data)" v-tooltip.top="'查看详情'" />
              <Button icon="pi pi-trash" text rounded severity="danger" @click="confirmDelete(data)" v-tooltip.top="'删除话题'" />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- 话题详情对话框 -->
    <Dialog v-model:visible="detailDialog.visible" :header="detailDialog.topic?.title || '话题详情'" :style="{width: '700px'}" :modal="true">
      <div v-if="detailDialog.topic" class="p-4">
        <div class="mb-4">
          <div class="flex items-center mb-2">
            <Avatar :image="detailDialog.topic.user?.avatar || '/default-avatar.png'" shape="circle" class="mr-2" />
            <span class="font-medium">{{ detailDialog.topic.user?.username }}</span>
            <span class="text-sm text-gray-500 ml-2">{{ formatDate(detailDialog.topic.created_at) }}</span>
          </div>
          <Tag :value="detailDialog.topic.category?.name || '未分类'" class="mb-2" />
          <div class="mt-2 p-3 bg-gray-50 rounded">
            <p class="whitespace-pre-line">{{ detailDialog.topic.content }}</p>
          </div>
        </div>
        
        <div class="mb-4">
          <h4 class="font-semibold mb-2">统计信息</h4>
          <div class="grid grid-cols-3 gap-4">
            <div>
              <span class="text-gray-600">浏览量:</span>
              <span class="ml-2">{{ detailDialog.topic.views || 0 }}</span>
            </div>
            <div>
              <span class="text-gray-600">回复数:</span>
              <span class="ml-2">{{ detailDialog.topic.replies_count || 0 }}</span>
            </div>
            <div>
              <span class="text-gray-600">点赞数:</span>
              <span class="ml-2">{{ detailDialog.topic.likes_count || 0 }}</span>
            </div>
          </div>
        </div>

        <div v-if="detailDialog.replies && detailDialog.replies.length > 0" class="mb-4">
          <h4 class="font-semibold mb-2">最新回复</h4>
          <div v-for="reply in detailDialog.replies" :key="reply.id" class="border-b border-gray-200 py-3 last:border-0">
            <div class="flex items-center mb-1">
              <Avatar :image="reply.user?.avatar || '/default-avatar.png'" shape="circle" size="small" class="mr-2" />
              <span class="font-medium">{{ reply.user?.username }}</span>
              <span class="text-sm text-gray-500 ml-2">{{ formatDate(reply.created_at) }}</span>
            </div>
            <p class="text-gray-700 ml-8">{{ reply.content }}</p>
          </div>
        </div>
      </div>
      <template #footer>
        <Button label="关闭" icon="pi pi-times" @click="detailDialog.visible = false" />
        <Button label="删除" icon="pi pi-trash" severity="danger" @click="confirmDeleteFromDialog()" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { adminApi } from '@/api/admin'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { format } from 'date-fns'
import { zhCN } from 'date-fns/locale'

import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Tag from 'primevue/tag'
import Avatar from 'primevue/avatar'

const router = useRouter()
const userStore = useUserStore()
const toast = useToast()
const confirm = useConfirm()

// 话题数据
const topics = ref([])
const loading = ref(false)

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  first: 0
})

// 筛选条件
const filters = reactive({
  search: '',
  category: null
})

// 分类选项
const categoryOptions = ref([])

// 详情对话框
const detailDialog = reactive({
  visible: false,
  topic: null,
  replies: []
})

// 格式化日期
const formatDate = (date) => {
  if (!date) return ''
  return format(new Date(date), 'yyyy-MM-dd HH:mm', { locale: zhCN })
}

// 加载论坛分类
const loadCategories = async () => {
  try {
    // 这里使用论坛控制器的API获取分类列表
    const response = await fetch('/api/forum/categories')
    const data = await response.json()
    categoryOptions.value = [{ id: null, name: '全部分类' }, ...data]
  } catch (error) {
    console.error('获取论坛分类失败:', error)
  }
}

// 加载话题列表
const loadTopics = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      search: filters.search || undefined,
      category_id: filters.category || undefined
    }
    
    const response = await adminApi.getTopics(params)
    topics.value = response.data.topics || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('获取话题列表失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '获取话题列表失败', life: 3000 })
  } finally {
    loading.value = false
  }
}

// 页码变化
const onPageChange = (event) => {
  pagination.page = event.page + 1
  pagination.first = event.first
  pagination.pageSize = event.rows
  loadTopics()
}

// 重置筛选条件
const resetFilters = () => {
  filters.search = ''
  filters.category = null
  loadTopics()
}

// 查看话题详情
const viewTopic = async (topic) => {
  detailDialog.topic = topic
  detailDialog.replies = []
  detailDialog.visible = true
  
  // 获取话题回复
  try {
    const response = await fetch(`/api/forum/topics/${topic.id}`)
    const data = await response.json()
    if (data.replies) {
      detailDialog.replies = data.replies.slice(0, 5) // 只显示前5条回复
    }
  } catch (error) {
    console.error('获取话题回复失败:', error)
  }
}

// 确认删除话题
const confirmDelete = (topic) => {
  confirm.require({
    message: `确定要删除话题 "${topic.title}" 吗？此操作将同时删除所有相关回复，且不可恢复。`,
    header: '删除确认',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    acceptLabel: '删除',
    rejectLabel: '取消',
    accept: () => deleteTopic(topic.id)
  })
}

// 从详情对话框确认删除
const confirmDeleteFromDialog = () => {
  if (!detailDialog.topic) return
  
  confirm.require({
    message: `确定要删除话题 "${detailDialog.topic.title}" 吗？此操作将同时删除所有相关回复，且不可恢复。`,
    header: '删除确认',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    acceptLabel: '删除',
    rejectLabel: '取消',
    accept: () => {
      deleteTopic(detailDialog.topic.id)
      detailDialog.visible = false
    }
  })
}

// 删除话题
const deleteTopic = async (id) => {
  try {
    await adminApi.deleteTopic(id)
    toast.add({ severity: 'success', summary: '成功', detail: '话题已删除', life: 3000 })
    loadTopics()
  } catch (error) {
    console.error('删除话题失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '删除话题失败', life: 3000 })
  }
}

onMounted(async () => {
  if (!userStore.isLoggedIn || userStore.user?.role !== 'admin') {
    router.push('/')
    return
  }
  
  // 加载分类
  await loadCategories()
  
  // 加载话题列表
  loadTopics()
})
</script>