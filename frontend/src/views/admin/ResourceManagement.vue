<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">资源管理</h1>
    </div>

    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <!-- 搜索和筛选 -->
      <div class="flex flex-col md:flex-row gap-4 mb-4">
        <div class="flex-1">
          <span class="p-input-icon-left w-full">
            <i class="pi pi-search" />
            <InputText v-model="filters.search" placeholder="搜索资源标题或描述" class="w-full" @keyup.enter="loadResources()" />
          </span>
        </div>
        <div class="flex gap-2">
          <Dropdown v-model="filters.status" :options="statusOptions" optionLabel="label" optionValue="value" placeholder="状态" class="w-32" />
          <Dropdown v-model="filters.category" :options="categoryOptions" optionLabel="name" optionValue="id" placeholder="分类" class="w-32" />
          <Button label="搜索" icon="pi pi-search" @click="loadResources()" />
          <Button label="重置" icon="pi pi-refresh" outlined @click="resetFilters()" />
        </div>
      </div>

      <!-- 资源表格 -->
      <DataTable 
        :value="resources" 
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
              <img :src="data.cover || '/default-resource.png'" class="w-10 h-10 object-cover rounded mr-2" />
              <span>{{ data.title }}</span>
            </div>
          </template>
        </Column>
        <Column field="user.username" header="上传者" :sortable="true">
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
        <Column field="status" header="状态" :sortable="true">
          <template #body="{data}">
            <Tag :severity="getStatusSeverity(data.status)" :value="getStatusLabel(data.status)" />
          </template>
        </Column>
        <Column field="created_at" header="上传时间" :sortable="true">
          <template #body="{data}">
            {{ formatDate(data.created_at) }}
          </template>
        </Column>
        <Column header="操作" style="width: 180px">
          <template #body="{data}">
            <div class="flex gap-1">
              <Button icon="pi pi-eye" text rounded @click="viewResource(data)" v-tooltip.top="'查看详情'" />
              <Button v-if="data.status === 'pending'" icon="pi pi-check" text rounded severity="success" @click="openReviewDialog(data, 'approved')" v-tooltip.top="'通过审核'" />
              <Button v-if="data.status === 'pending'" icon="pi pi-times" text rounded severity="danger" @click="openReviewDialog(data, 'rejected')" v-tooltip.top="'拒绝审核'" />
              <Button icon="pi pi-trash" text rounded severity="danger" @click="confirmDelete(data)" v-tooltip.top="'删除资源'" />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- 资源详情对话框 -->
    <Dialog v-model:visible="detailDialog.visible" :header="detailDialog.resource?.title || '资源详情'" :style="{width: '700px'}" :modal="true">
      <div v-if="detailDialog.resource" class="p-4">
        <div class="flex mb-4">
          <img :src="detailDialog.resource.cover || '/default-resource.png'" class="w-32 h-32 object-cover rounded mr-4" />
          <div class="flex-1">
            <h3 class="text-xl font-bold mb-2">{{ detailDialog.resource.title }}</h3>
            <div class="flex items-center text-sm text-gray-600 mb-1">
              <i class="pi pi-user mr-1"></i>
              <span>{{ detailDialog.resource.user?.username }}</span>
            </div>
            <div class="flex items-center text-sm text-gray-600 mb-1">
              <i class="pi pi-tag mr-1"></i>
              <span>{{ detailDialog.resource.category?.name || '未分类' }}</span>
            </div>
            <div class="flex items-center text-sm text-gray-600">
              <i class="pi pi-calendar mr-1"></i>
              <span>{{ formatDate(detailDialog.resource.created_at) }}</span>
            </div>
          </div>
        </div>
        
        <div class="mb-4">
          <h4 class="font-semibold mb-2">描述</h4>
          <p class="text-gray-700 whitespace-pre-line">{{ detailDialog.resource.description || '无描述' }}</p>
        </div>
        
        <div class="mb-4">
          <h4 class="font-semibold mb-2">下载信息</h4>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <span class="text-gray-600">下载次数:</span>
              <span class="ml-2">{{ detailDialog.resource.downloads || 0 }}</span>
            </div>
            <div>
              <span class="text-gray-600">文件大小:</span>
              <span class="ml-2">{{ formatFileSize(detailDialog.resource.file_size) }}</span>
            </div>
            <div>
              <span class="text-gray-600">文件类型:</span>
              <span class="ml-2">{{ detailDialog.resource.file_type || '未知' }}</span>
            </div>
            <div>
              <span class="text-gray-600">状态:</span>
              <Tag class="ml-2" :severity="getStatusSeverity(detailDialog.resource.status)" :value="getStatusLabel(detailDialog.resource.status)" />
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <Button label="关闭" icon="pi pi-times" @click="detailDialog.visible = false" />
      </template>
    </Dialog>

    <!-- 审核对话框 -->
    <Dialog v-model:visible="reviewDialog.visible" :header="reviewDialog.action === 'approved' ? '通过审核' : '拒绝审核'" :style="{width: '500px'}" :modal="true">
      <div class="p-4">
        <div class="mb-4">
          <h3 class="font-semibold mb-2">资源信息</h3>
          <p><span class="text-gray-600">标题:</span> {{ reviewDialog.resource?.title }}</p>
          <p><span class="text-gray-600">上传者:</span> {{ reviewDialog.resource?.user?.username }}</p>
        </div>
        
        <div class="field">
          <label for="message" class="block text-sm font-medium text-gray-700 mb-1">审核留言 (可选)</label>
          <Textarea id="message" v-model="reviewDialog.message" rows="3" class="w-full" placeholder="输入审核留言，将通知给资源上传者" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" icon="pi pi-times" outlined @click="reviewDialog.visible = false" />
        <Button 
          :label="reviewDialog.action === 'approved' ? '通过' : '拒绝'" 
          :icon="reviewDialog.action === 'approved' ? 'pi pi-check' : 'pi pi-times'" 
          :severity="reviewDialog.action === 'approved' ? 'success' : 'danger'" 
          @click="reviewResource" 
          :loading="reviewDialog.loading" 
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
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
import Textarea from 'primevue/textarea'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const toast = useToast()
const confirm = useConfirm()

// 资源数据
const resources = ref([])
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
  status: 'pending',  // 默认显示待审核
  category: null
})

// 状态选项
const statusOptions = [
  { label: '全部状态', value: null },
  { label: '待审核', value: 'pending' },
  { label: '已通过', value: 'approved' },
  { label: '已拒绝', value: 'rejected' }
]

// 分类选项
const categoryOptions = ref([])

// 详情对话框
const detailDialog = reactive({
  visible: false,
  resource: null
})

// 审核对话框
const reviewDialog = reactive({
  visible: false,
  resource: null,
  action: 'approved',  // 'approved' 或 'rejected'
  message: '',
  loading: false
})

// 格式化日期
const formatDate = (date) => {
  if (!date) return ''
  return format(new Date(date), 'yyyy-MM-dd HH:mm', { locale: zhCN })
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 获取状态标签
const getStatusLabel = (status) => {
  switch (status) {
    case 'pending': return '待审核'
    case 'approved': return '已通过'
    case 'rejected': return '已拒绝'
    default: return status
  }
}

// 获取状态样式
const getStatusSeverity = (status) => {
  switch (status) {
    case 'pending': return 'warning'
    case 'approved': return 'success'
    case 'rejected': return 'danger'
    default: return 'info'
  }
}

// 加载资源分类
const loadCategories = async () => {
  try {
    // 这里使用资源控制器的API获取分类列表
    const response = await fetch('/api/resources/categories')
    const data = await response.json()
    categoryOptions.value = [{ id: null, name: '全部分类' }, ...data]
  } catch (error) {
    console.error('获取资源分类失败:', error)
  }
}

// 加载资源列表
const loadResources = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      search: filters.search || undefined,
      status: filters.status || undefined,
      category_id: filters.category || undefined
    }
    
    const response = await adminApi.getResources(params)
    resources.value = response.data.resources || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('获取资源列表失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '获取资源列表失败', life: 3000 })
  } finally {
    loading.value = false
  }
}

// 页码变化
const onPageChange = (event) => {
  pagination.page = event.page + 1
  pagination.first = event.first
  pagination.pageSize = event.rows
  loadResources()
}

// 重置筛选条件
const resetFilters = () => {
  filters.search = ''
  filters.status = 'pending'
  filters.category = null
  loadResources()
}

// 查看资源详情
const viewResource = (resource) => {
  detailDialog.resource = resource
  detailDialog.visible = true
}

// 打开审核对话框
const openReviewDialog = (resource, action) => {
  reviewDialog.resource = resource
  reviewDialog.action = action
  reviewDialog.message = ''
  reviewDialog.visible = true
}

// 审核资源
const reviewResource = async () => {
  if (!reviewDialog.resource) return
  
  reviewDialog.loading = true
  try {
    await adminApi.reviewResource(reviewDialog.resource.id, {
      status: reviewDialog.action,
      message: reviewDialog.message
    })
    
    toast.add({ 
      severity: 'success', 
      summary: '成功', 
      detail: reviewDialog.action === 'approved' ? '资源已通过审核' : '资源已被拒绝', 
      life: 3000 
    })
    
    reviewDialog.visible = false
    loadResources()
  } catch (error) {
    console.error('审核资源失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '审核资源失败', life: 3000 })
  } finally {
    reviewDialog.loading = false
  }
}

// 确认删除资源
const confirmDelete = (resource) => {
  confirm.require({
    message: `确定要删除资源 "${resource.title}" 吗？此操作不可恢复。`,
    header: '删除确认',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    acceptLabel: '删除',
    rejectLabel: '取消',
    accept: () => deleteResource(resource.id)
  })
}

// 删除资源
const deleteResource = async (id) => {
  try {
    await adminApi.deleteResource(id)
    toast.add({ severity: 'success', summary: '成功', detail: '资源已删除', life: 3000 })
    loadResources()
  } catch (error) {
    console.error('删除资源失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '删除资源失败', life: 3000 })
  }
}

onMounted(async () => {
  if (!userStore.isLoggedIn || userStore.user?.role !== 'admin') {
    router.push('/')
    return
  }
  
  // 加载分类
  await loadCategories()
  
  // 如果URL中有资源ID参数，直接打开该资源的详情
  const resourceId = route.query.id
  if (resourceId) {
    try {
      const response = await adminApi.getResources({ id: resourceId })
      if (response.data.resources && response.data.resources.length > 0) {
        viewResource(response.data.resources[0])
      }
    } catch (error) {
      console.error('获取资源详情失败:', error)
    }
  }
  
  // 加载资源列表
  loadResources()
})
</script>