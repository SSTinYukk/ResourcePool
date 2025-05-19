<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">我的积分</h1>
      <Button 
        label="返回" 
        icon="pi pi-arrow-left" 
        class="p-button-secondary"
        @click="goBack"
      />
    </div>

    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold text-gray-700">当前积分</h2>
          <p v-if="isLoading" class="text-3xl font-bold text-gray-400">加载中...</p>
          <p v-else-if="error" class="text-3xl font-bold text-red-500">加载失败</p>
          <p v-else class="text-3xl font-bold text-blue-600">{{ totalPoints }}</p>
        </div>
        <Button 
          label="积分规则" 
          icon="pi pi-info-circle" 
          class="p-button-text"
          @click="showPointsRules"
        />
      </div>
      <p v-if="error" class="mt-2 text-sm text-red-500">{{ error }}</p>
    </div>

    <div class="bg-white rounded-lg shadow-md p-6">
      <div v-if="isLoading" class="flex justify-center py-8">
        <i class="pi pi-spinner pi-spin text-2xl text-blue-500"></i>
      </div>
      <div v-else-if="error" class="flex justify-center py-8 text-red-500">
        {{ error }}
      </div>
      <DataTable 
        v-else
        :value="pointsHistory" 
        :paginator="true" 
        :rows="10"
        paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport RowsPerPageDropdown"
        currentPageReportTemplate="显示 {first} 到 {last} 共 {totalRecords} 条"
        responsiveLayout="scroll"
      >
        <Column field="date" header="日期" :sortable="true"></Column>
        <Column field="description" header="说明" :sortable="true"></Column>
        <Column field="points" header="积分" :sortable="true">
          <template #body="slotProps">
            <span :class="slotProps.data.points > 0 ? 'text-green-600' : 'text-red-600'">
              {{ slotProps.data.points > 0 ? '+' : '' }}{{ slotProps.data.points }}
            </span>
          </template>
        </Column>
      </DataTable>
    </div>

    <Dialog 
      v-model:visible="showRulesDialog" 
      header="积分规则" 
      :modal="true"
      class="w-full max-w-md"
    >
      <div class="space-y-4">
        <p class="font-medium">积分获取方式:</p>
        <ul class="list-disc pl-5 space-y-2">
          <li>每日签到: +10积分</li>
          <li>上传资源: +50积分</li>
          <li>资源被下载: +5积分/次</li>
          <li>资源被点赞: +2积分/次</li>
        </ul>
      </div>

      <template #footer>
        <Button 
          label="关闭" 
          icon="pi pi-times" 
          class="p-button-text"
          @click="showRulesDialog = false"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import axios from 'axios'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'

const router = useRouter()
const userStore = useUserStore()

const totalPoints = ref(0)
const pointsHistory = ref([])
const showRulesDialog = ref(false)

const isLoading = ref(false)
const error = ref(null)

const fetchPointsData = async () => {
  isLoading.value = true
  error.value = null
  
  try {
    const { data } = await axios.get('/api/user/points')
    totalPoints.value = data.totalPoints
    pointsHistory.value = data.history
  } catch (err) {
    console.error('获取积分数据错误:', err)
    error.value = err.response?.data?.message || '获取积分数据失败'
  } finally {
    isLoading.value = false
  }
}

const showPointsRules = () => {
  showRulesDialog.value = true
}

const goBack = () => {
  router.push('/')
}

onMounted(() => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
  } else {
    fetchPointsData()
  }
})
</script>