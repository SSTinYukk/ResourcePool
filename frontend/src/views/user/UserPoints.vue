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
          <p class="text-3xl font-bold text-blue-600">{{ totalPoints }}</p>
        </div>
        <Button 
          label="积分规则" 
          icon="pi pi-info-circle" 
          class="p-button-text"
          @click="showPointsRules"
        />
      </div>
    </div>

    <div class="bg-white rounded-lg shadow-md p-6">
      <DataTable 
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
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'

const router = useRouter()
const userStore = useUserStore()

const totalPoints = ref(0)
const pointsHistory = ref([])
const showRulesDialog = ref(false)

const fetchPointsData = async () => {
  // TODO: Replace with actual API call
  await new Promise(resolve => setTimeout(resolve, 300))
  
  totalPoints.value = 1250
  pointsHistory.value = [
    { id: 1, date: '2023-05-10', description: '每日签到', points: 10 },
    { id: 2, date: '2023-05-09', description: '资源被下载', points: 5 },
    { id: 3, date: '2023-05-08', description: '上传资源', points: 50 },
    { id: 4, date: '2023-05-07', description: '资源被点赞', points: 2 },
    { id: 5, date: '2023-05-06', description: '每日签到', points: 10 },
  ]
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