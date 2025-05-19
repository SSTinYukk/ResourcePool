<template>
  <div class="container mx-auto px-4 py-8 max-w-7xl">
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-2xl font-bold text-gray-800 mb-1">我的资源</h1>
      <div class="flex space-x-2">
        <Button 
          label="上传新资源" 
          icon="pi pi-upload" 
          class="p-button-primary"
          @click="navigateToUpload"
        />
        <Button 
          label="返回" 
          icon="pi pi-arrow-left" 
          class="p-button-secondary"
          @click="goBack"
        />
      </div>
    </div>

    <div class="bg-white rounded-lg shadow-lg p-6 transition-shadow duration-300 hover:shadow-xl">
      <DataTable 
        :value="resources" 
        :paginator="true" 
        :rows="10"
        paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
        currentPageReportTemplate="显示 {first} 到 {last} 共 {totalRecords} 条"
        responsiveLayout="scroll"
      >
        <Column field="title" header="资源名称" :sortable="true">
          <template #body="slotProps">
            <div class="cursor-pointer text-blue-600 hover:text-blue-800" @click="navigateToResource(slotProps.data.id)">
              {{ slotProps.data.title }}
            </div>
          </template>
        </Column>
        <Column field="category.name" header="分类" :sortable="true"></Column>
        <Column field="created_at" header="上传时间" :sortable="true">
          <template #body="slotProps">
            {{ formatDate(slotProps.data.created_at) }}
          </template>
        </Column>
        <Column field="status" header="审核状态" :sortable="true">
          <template #body="slotProps">
            <Tag 
              :value="getStatusLabel(slotProps.data.status)" 
              :severity="getStatusSeverity(slotProps.data.status)"
              class="px-3 py-1 text-sm"
              rounded
            />
          </template>
        </Column>
        <Column header="操作">
          <template #body="slotProps">
            <Button 
              icon="pi pi-download" 
              class="p-button-rounded p-button-text"
              @click="downloadResource(slotProps.data)"
            />
            <Button 
              icon="pi pi-trash" 
              class="p-button-rounded p-button-text p-button-danger"
              @click="confirmDelete(slotProps.data)"
            />
          </template>
        </Column>
      </DataTable>
    </div>

    <Dialog 
      v-model:visible="showDeleteDialog" 
      header="确认删除" 
      :modal="true"
      class="w-full max-w-md"
    >
      <div class="space-y-4">
        <p>确定要删除资源 "{{ selectedResource?.name }}" 吗？</p>
      </div>

      <template #footer>
        <Button 
          label="取消" 
          icon="pi pi-times" 
          class="p-button-text"
          @click="showDeleteDialog = false"
        />
        <Button 
          label="删除" 
          icon="pi pi-check" 
          class="p-button-danger"
          @click="deleteResource"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Tag from 'primevue/tag'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const getStatusSeverity = (status) => {
  switch(status) {
    case 'pending': return 'info'
    case 'approved': return 'success'
    case 'rejected': return 'danger'
    default: return 'warning'
  }
}

const getStatusLabel = (status) => {
  switch(status) {
    case 'pending': return '待审核'
    case 'approved': return '已通过'
    case 'rejected': return '已拒绝'
    default: return status
  }
}

const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const resources = ref([])
const showDeleteDialog = ref(false)
const selectedResource = ref(null)

const fetchResources = async () => {
  try {
    const response = await fetch('/api/user/resources', {
      headers: {
        'Authorization': `Bearer ${userStore.token}`
      }
    })
    
    if (!response.ok) throw new Error('获取资源失败')
    
    const data = await response.json()
    console.log('Fetched resources:', data.resources)
    resources.value = data.resources
  } catch (error) {
    console.error('获取资源出错:', error)
  }
}

const downloadResource = async (resource) => {
  try {
    const response = await fetch(`/api/download/${resource.id}`, {
      headers: {
        'Authorization': `Bearer ${userStore.token}`,
        'Content-Type': 'application/json'
      }
    })
    
    if (!response.ok) throw new Error('获取下载链接失败')
    
    const data = await response.json()
    
    if (data.url) {
      // 创建隐藏的a标签触发下载
      const link = document.createElement('a')
      link.href = data.url
      link.download = resource.title || 'download'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    }
  } catch (error) {
    console.error('下载资源出错:', error)
  }
}

const confirmDelete = (resource) => {
  selectedResource.value = resource
  showDeleteDialog.value = true
}

const deleteResource = async () => {
  try {
    const response = await fetch(`/api/user/resources/${selectedResource.value.id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${userStore.token}`
      }
    })
    
    if (!response.ok) throw new Error('删除资源失败')
    
    await fetchResources()
    showDeleteDialog.value = false
  } catch (error) {
    console.error('删除资源出错:', error)
  }
}

const goBack = () => {
  router.push('/')
}

const navigateToUpload = () => {
  router.push('/resources/upload')
}

const navigateToResource = (resourceId) => {
  router.push(`/resources/${resourceId}`)
}

onMounted(() => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
  } else {
    fetchResources()
  }
})
</script>