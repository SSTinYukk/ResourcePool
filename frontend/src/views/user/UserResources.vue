<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">我的资源</h1>
      <Button 
        label="返回" 
        icon="pi pi-arrow-left" 
        class="p-button-secondary"
        @click="goBack"
      />
    </div>

    <div class="bg-white rounded-lg shadow-md p-6">
      <DataTable 
        :value="resources" 
        :paginator="true" 
        :rows="10"
        paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
        currentPageReportTemplate="显示 {first} 到 {last} 共 {totalRecords} 条"
        responsiveLayout="scroll"
      >
        <Column field="name" header="资源名称" :sortable="true"></Column>
        <Column field="type" header="类型" :sortable="true"></Column>
        <Column field="size" header="大小" :sortable="true"></Column>
        <Column field="uploadDate" header="上传时间" :sortable="true"></Column>
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

const router = useRouter()
const userStore = useUserStore()

const resources = ref([])
const showDeleteDialog = ref(false)
const selectedResource = ref(null)

const fetchResources = async () => {
  // TODO: Replace with actual API call
  await new Promise(resolve => setTimeout(resolve, 300))
  
  resources.value = [
    { id: 1, name: '文档1.pdf', type: 'PDF', size: '2.5MB', uploadDate: '2023-05-10' },
    { id: 2, name: '图片1.jpg', type: '图片', size: '1.2MB', uploadDate: '2023-05-15' },
    { id: 3, name: '视频1.mp4', type: '视频', size: '15.7MB', uploadDate: '2023-05-20' },
  ]
}

const downloadResource = (resource) => {
  // TODO: Implement download logic
  console.log('Downloading:', resource)
}

const confirmDelete = (resource) => {
  selectedResource.value = resource
  showDeleteDialog.value = true
}

const deleteResource = async () => {
  // TODO: Implement delete logic
  showDeleteDialog.value = false
}

const goBack = () => {
  router.push('/')
}

onMounted(() => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
  } else {
    fetchResources()
  }
})
</script>