<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">积分管理</h1>
      <Button label="调整积分" icon="pi pi-plus" @click="openAdjustDialog()" />
    </div>

    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <!-- 搜索和筛选 -->
      <div class="flex flex-col md:flex-row gap-4 mb-4">
        <div class="flex-1">
          <span class="p-input-icon-left w-full">
            <i class="pi pi-search" />
            <InputText v-model="filters.search" placeholder="搜索用户名" class="w-full" @keyup.enter="loadPointsRecords()" />
          </span>
        </div>
        <div class="flex gap-2">
          <Dropdown v-model="filters.type" :options="typeOptions" optionLabel="name" optionValue="value" placeholder="类型" class="w-32" />
          <Button label="搜索" icon="pi pi-search" @click="loadPointsRecords()" />
          <Button label="重置" icon="pi pi-refresh" outlined @click="resetFilters()" />
        </div>
      </div>

      <!-- 积分记录表格 -->
      <DataTable 
        :value="pointsRecords" 
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
        <Column field="user.username" header="用户" :sortable="true">
          <template #body="{data}">
            <div class="flex items-center">
              <Avatar :image="data.user?.avatar || '/default-avatar.png'" shape="circle" class="mr-2" size="small" />
              <span>{{ data.user?.username }}</span>
            </div>
          </template>
        </Column>
        <Column field="type" header="类型" :sortable="true">
          <template #body="{data}">
            <Tag :severity="getTypeSeverity(data.type)" :value="getTypeLabel(data.type)" />
          </template>
        </Column>
        <Column field="points" header="积分变动" :sortable="true">
          <template #body="{data}">
            <span :class="{'text-green-600': data.points > 0, 'text-red-600': data.points < 0}">
              {{ data.points > 0 ? '+' : '' }}{{ data.points }}
            </span>
          </template>
        </Column>
        <Column field="description" header="描述" :sortable="false"></Column>
        <Column field="created_at" header="时间" :sortable="true">
          <template #body="{data}">
            {{ formatDate(data.created_at) }}
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- 积分调整对话框 -->
    <Dialog v-model:visible="adjustDialog.visible" header="调整用户积分" :style="{width: '500px'}" :modal="true">
      <div class="p-fluid">
        <div class="field mb-4">
          <label for="username" class="block mb-2">用户名</label>
          <AutoComplete 
            id="username" 
            v-model="adjustDialog.username" 
            :suggestions="userSuggestions" 
            @complete="searchUsers($event)" 
            field="username" 
            placeholder="输入用户名搜索"
            class="w-full"
          />
        </div>
        <div class="field mb-4">
          <label for="points" class="block mb-2">积分变动</label>
          <InputNumber 
            id="points" 
            v-model="adjustDialog.points" 
            placeholder="输入积分变动值（正数增加，负数减少）" 
            class="w-full"
          />
        </div>
        <div class="field mb-4">
          <label for="description" class="block mb-2">描述</label>
          <Textarea 
            id="description" 
            v-model="adjustDialog.description" 
            rows="3" 
            placeholder="输入积分调整原因" 
            class="w-full"
          />
        </div>
      </div>
      <template #footer>
        <Button label="取消" icon="pi pi-times" outlined @click="adjustDialog.visible = false" />
        <Button label="确认" icon="pi pi-check" @click="submitPointsAdjustment()" :loading="adjustDialog.loading" />
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
import { format } from 'date-fns'
import { zhCN } from 'date-fns/locale'

import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Tag from 'primevue/tag'
import Avatar from 'primevue/avatar'
import AutoComplete from 'primevue/autocomplete'

const router = useRouter()
const userStore = useUserStore()
const toast = useToast()

// 积分记录数据
const pointsRecords = ref([])
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
  type: null
})

// 类型选项
const typeOptions = ref([
  { name: '全部类型', value: null },
  { name: '资源上传', value: 'resource_upload' },
  { name: '资源下载', value: 'resource_download' },
  { name: '话题发布', value: 'topic_create' },
  { name: '回复发布', value: 'reply_create' },
  { name: '管理员调整', value: 'admin_adjustment' },
  { name: '系统奖励', value: 'system_reward' }
])

// 积分调整对话框
const adjustDialog = reactive({
  visible: false,
  username: '',
  points: null,
  description: '',
  loading: false
})

// 用户搜索建议
const userSuggestions = ref([])

// 格式化日期
const formatDate = (date) => {
  if (!date) return ''
  return format(new Date(date), 'yyyy-MM-dd HH:mm', { locale: zhCN })
}

// 获取类型标签
const getTypeLabel = (type) => {
  const option = typeOptions.value.find(opt => opt.value === type)
  return option ? option.name : type
}

// 获取类型样式
const getTypeSeverity = (type) => {
  switch (type) {
    case 'resource_upload':
      return 'success'
    case 'resource_download':
      return 'info'
    case 'topic_create':
      return 'warning'
    case 'reply_create':
      return 'warning'
    case 'admin_adjustment':
      return 'danger'
    case 'system_reward':
      return 'success'
    default:
      return null
  }
}

// 加载积分记录
const loadPointsRecords = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      search: filters.search || undefined,
      type: filters.type || undefined
    }
    
    const response = await adminApi.getPointsRecords(params)
    pointsRecords.value = response.data.records || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('获取积分记录失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '获取积分记录失败', life: 3000 })
  } finally {
    loading.value = false
  }
}

// 页码变化
const onPageChange = (event) => {
  pagination.page = event.page + 1
  pagination.first = event.first
  pagination.pageSize = event.rows
  loadPointsRecords()
}

// 重置筛选条件
const resetFilters = () => {
  filters.search = ''
  filters.type = null
  loadPointsRecords()
}

// 打开积分调整对话框
const openAdjustDialog = () => {
  adjustDialog.username = ''
  adjustDialog.points = null
  adjustDialog.description = ''
  adjustDialog.visible = true
}

// 搜索用户
const searchUsers = async (event) => {
  try {
    const response = await adminApi.getUsers({
      search: event.query,
      pageSize: 5
    })
    userSuggestions.value = response.data.users || []
  } catch (error) {
    console.error('搜索用户失败:', error)
  }
}

// 提交积分调整
const submitPointsAdjustment = async () => {
  // 表单验证
  if (!adjustDialog.username) {
    toast.add({ severity: 'error', summary: '错误', detail: '请输入用户名', life: 3000 })
    return
  }
  
  if (adjustDialog.points === null || adjustDialog.points === 0) {
    toast.add({ severity: 'error', summary: '错误', detail: '请输入有效的积分变动值', life: 3000 })
    return
  }
  
  if (!adjustDialog.description) {
    toast.add({ severity: 'error', summary: '错误', detail: '请输入积分调整原因', life: 3000 })
    return
  }
  
  adjustDialog.loading = true
  try {
    // 如果username是对象（从AutoComplete选择），则获取用户ID
    const username = typeof adjustDialog.username === 'object' 
      ? adjustDialog.username.username 
      : adjustDialog.username
    
    await adminApi.adjustPoints({
      username,
      points: adjustDialog.points,
      description: adjustDialog.description
    })
    
    toast.add({ severity: 'success', summary: '成功', detail: '积分调整成功', life: 3000 })
    adjustDialog.visible = false
    loadPointsRecords()
  } catch (error) {
    console.error('积分调整失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '积分调整失败: ' + (error.message || '未知错误'), life: 3000 })
  } finally {
    adjustDialog.loading = false
  }
}

onMounted(async () => {
  if (!userStore.isLoggedIn || userStore.user?.role !== 'admin') {
    router.push('/')
    return
  }
  
  // 加载积分记录
  loadPointsRecords()
})
</script>