<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">数据统计</h1>
      <div class="flex gap-2">
        <Dropdown v-model="timeRange" :options="timeRangeOptions" optionLabel="name" optionValue="value" placeholder="时间范围" @change="loadAllStats()" />
        <Button icon="pi pi-refresh" @click="loadAllStats()" v-tooltip.top="'刷新数据'" />
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center items-center py-8">
      <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="4" fill="var(--surface-ground)" animationDuration=".5s" />
      <span class="ml-3 text-lg">加载统计数据...</span>
    </div>

    <div v-else>
      <!-- 总览卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold text-gray-700">总用户数</h3>
            <i class="pi pi-users text-blue-500 text-2xl"></i>
          </div>
          <div class="text-3xl font-bold">{{ overview.totalUsers }}</div>
          <div class="text-sm text-gray-500 mt-2">
            <span class="text-green-500">+{{ overview.newUsers }}</span> 新增用户
          </div>
        </div>

        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold text-gray-700">总资源数</h3>
            <i class="pi pi-file text-orange-500 text-2xl"></i>
          </div>
          <div class="text-3xl font-bold">{{ overview.totalResources }}</div>
          <div class="text-sm text-gray-500 mt-2">
            <span class="text-green-500">+{{ overview.newResources }}</span> 新增资源
          </div>
        </div>

        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold text-gray-700">总话题数</h3>
            <i class="pi pi-comments text-purple-500 text-2xl"></i>
          </div>
          <div class="text-3xl font-bold">{{ overview.totalTopics }}</div>
          <div class="text-sm text-gray-500 mt-2">
            <span class="text-green-500">+{{ overview.newTopics }}</span> 新增话题
          </div>
        </div>

        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold text-gray-700">总下载次数</h3>
            <i class="pi pi-download text-green-500 text-2xl"></i>
          </div>
          <div class="text-3xl font-bold">{{ overview.totalDownloads }}</div>
          <div class="text-sm text-gray-500 mt-2">
            <span class="text-green-500">+{{ overview.newDownloads }}</span> 新增下载
          </div>
        </div>
      </div>

      <!-- 用户增长图表 -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-6">
        <h3 class="text-lg font-semibold text-gray-700 mb-4">用户增长趋势</h3>
        <Chart type="line" :data="userGrowthData" :options="userGrowthOptions" style="height: 300px" />
      </div>

      <!-- 资源和话题统计 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <!-- 资源分类统计 -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold text-gray-700 mb-4">资源分类统计</h3>
          <Chart type="pie" :data="resourceCategoryData" :options="pieChartOptions" style="height: 300px" />
        </div>

        <!-- 话题分类统计 -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold text-gray-700 mb-4">话题分类统计</h3>
          <Chart type="pie" :data="topicCategoryData" :options="pieChartOptions" style="height: 300px" />
        </div>
      </div>

      <!-- 活跃度统计 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <!-- 资源上传和下载趋势 -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold text-gray-700 mb-4">资源上传和下载趋势</h3>
          <Chart type="bar" :data="resourceActivityData" :options="barChartOptions" style="height: 300px" />
        </div>

        <!-- 话题和回复趋势 -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold text-gray-700 mb-4">话题和回复趋势</h3>
          <Chart type="bar" :data="forumActivityData" :options="barChartOptions" style="height: 300px" />
        </div>
      </div>

      <!-- 热门资源和话题 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- 热门资源 -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold text-gray-700 mb-4">热门资源</h3>
          <DataTable :value="popularResources" responsiveLayout="scroll" class="p-datatable-sm">
            <Column field="title" header="资源名称">
              <template #body="{data}">
                <div class="truncate max-w-xs">{{ data.title }}</div>
              </template>
            </Column>
            <Column field="downloads" header="下载次数">
              <template #body="{data}">
                <Badge :value="data.downloads" severity="success" />
              </template>
            </Column>
            <Column field="category" header="分类">
              <template #body="{data}">
                <Tag :value="data.category" />
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- 热门话题 -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold text-gray-700 mb-4">热门话题</h3>
          <DataTable :value="popularTopics" responsiveLayout="scroll" class="p-datatable-sm">
            <Column field="title" header="话题标题">
              <template #body="{data}">
                <div class="truncate max-w-xs">{{ data.title }}</div>
              </template>
            </Column>
            <Column field="views" header="浏览次数">
              <template #body="{data}">
                <Badge :value="data.views" severity="info" />
              </template>
            </Column>
            <Column field="replies" header="回复数">
              <template #body="{data}">
                <Badge :value="data.replies" severity="warning" />
              </template>
            </Column>
          </DataTable>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { adminApi } from '@/api/admin'
import { useToast } from 'primevue/usetoast'

import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import Chart from 'primevue/chart'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Badge from 'primevue/badge'
import ProgressSpinner from 'primevue/progressspinner'

const router = useRouter()
const userStore = useUserStore()
const toast = useToast()

// 加载状态
const loading = ref(false)

// 时间范围
const timeRange = ref('last30days')
const timeRangeOptions = [
  { name: '最近7天', value: 'last7days' },
  { name: '最近30天', value: 'last30days' },
  { name: '最近90天', value: 'last90days' },
  { name: '今年', value: 'thisyear' },
  { name: '全部', value: 'all' }
]

// 总览数据
const overview = reactive({
  totalUsers: 0,
  newUsers: 0,
  totalResources: 0,
  newResources: 0,
  totalTopics: 0,
  newTopics: 0,
  totalDownloads: 0,
  newDownloads: 0
})

// 热门资源和话题
const popularResources = ref([])
const popularTopics = ref([])

// 用户增长数据
const userGrowthData = ref({
  labels: [],
  datasets: [
    {
      label: '新增用户',
      data: [],
      fill: false,
      borderColor: '#4CAF50',
      tension: 0.4
    },
    {
      label: '累计用户',
      data: [],
      fill: false,
      borderColor: '#2196F3',
      tension: 0.4
    }
  ]
})

// 资源分类数据
const resourceCategoryData = ref({
  labels: [],
  datasets: [
    {
      data: [],
      backgroundColor: [
        '#FF6384',
        '#36A2EB',
        '#FFCE56',
        '#4BC0C0',
        '#9966FF',
        '#FF9F40',
        '#C9CBCF'
      ]
    }
  ]
})

// 话题分类数据
const topicCategoryData = ref({
  labels: [],
  datasets: [
    {
      data: [],
      backgroundColor: [
        '#FF6384',
        '#36A2EB',
        '#FFCE56',
        '#4BC0C0',
        '#9966FF',
        '#FF9F40',
        '#C9CBCF'
      ]
    }
  ]
})

// 资源活动数据
const resourceActivityData = ref({
  labels: [],
  datasets: [
    {
      label: '资源上传',
      backgroundColor: '#42A5F5',
      data: []
    },
    {
      label: '资源下载',
      backgroundColor: '#66BB6A',
      data: []
    }
  ]
})

// 论坛活动数据
const forumActivityData = ref({
  labels: [],
  datasets: [
    {
      label: '新增话题',
      backgroundColor: '#7E57C2',
      data: []
    },
    {
      label: '新增回复',
      backgroundColor: '#FFB74D',
      data: []
    }
  ]
})

// 图表选项
const userGrowthOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top'
    }
  },
  scales: {
    y: {
      beginAtZero: true
    }
  }
}

const pieChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'right'
    }
  }
}

const barChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top'
    }
  },
  scales: {
    y: {
      beginAtZero: true
    }
  }
}

// 加载所有统计数据
const loadAllStats = async () => {
  loading.value = true
  try {
    // 加载总览数据
    await loadOverviewStats()
    
    // 加载用户增长数据
    await loadUserGrowthStats()
    
    // 加载资源分类数据
    await loadResourceCategoryStats()
    
    // 加载话题分类数据
    await loadTopicCategoryStats()
    
    // 加载资源活动数据
    await loadResourceActivityStats()
    
    // 加载论坛活动数据
    await loadForumActivityStats()
    
    // 加载热门资源和话题
    await loadPopularItems()
  } catch (error) {
    console.error('加载统计数据失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '加载统计数据失败', life: 3000 })
  } finally {
    loading.value = false
  }
}

// 加载总览统计
const loadOverviewStats = async () => {
  try {
    const response = await adminApi.getOverviewStats({ timeRange: timeRange.value })
    const data = response.data
    
    // 更新总览数据
    overview.totalUsers = data.totalUsers || 0
    overview.newUsers = data.newUsers || 0
    overview.totalResources = data.totalResources || 0
    overview.newResources = data.newResources || 0
    overview.totalTopics = data.totalTopics || 0
    overview.newTopics = data.newTopics || 0
    overview.totalDownloads = data.totalDownloads || 0
    overview.newDownloads = data.newDownloads || 0
  } catch (error) {
    console.error('加载总览统计失败:', error)
    throw error
  }
}

// 加载用户增长统计
const loadUserGrowthStats = async () => {
  try {
    const response = await adminApi.getUserGrowthStats({ timeRange: timeRange.value })
    const data = response.data
    
    // 更新用户增长数据
    userGrowthData.value.labels = data.labels || []
    userGrowthData.value.datasets[0].data = data.newUsers || []
    userGrowthData.value.datasets[1].data = data.totalUsers || []
  } catch (error) {
    console.error('加载用户增长统计失败:', error)
    throw error
  }
}

// 加载资源分类统计
const loadResourceCategoryStats = async () => {
  try {
    const response = await adminApi.getResourceCategoryStats({ timeRange: timeRange.value })
    const data = response.data
    
    // 更新资源分类数据
    resourceCategoryData.value.labels = data.categories || []
    resourceCategoryData.value.datasets[0].data = data.counts || []
  } catch (error) {
    console.error('加载资源分类统计失败:', error)
    throw error
  }
}

// 加载话题分类统计
const loadTopicCategoryStats = async () => {
  try {
    const response = await adminApi.getTopicCategoryStats({ timeRange: timeRange.value })
    const data = response.data
    
    // 更新话题分类数据
    topicCategoryData.value.labels = data.categories || []
    topicCategoryData.value.datasets[0].data = data.counts || []
  } catch (error) {
    console.error('加载话题分类统计失败:', error)
    throw error
  }
}

// 加载资源活动统计
const loadResourceActivityStats = async () => {
  try {
    const response = await adminApi.getResourceActivityStats({ timeRange: timeRange.value })
    const data = response.data
    
    // 更新资源活动数据
    resourceActivityData.value.labels = data.labels || []
    resourceActivityData.value.datasets[0].data = data.uploads || []
    resourceActivityData.value.datasets[1].data = data.downloads || []
  } catch (error) {
    console.error('加载资源活动统计失败:', error)
    throw error
  }
}

// 加载论坛活动统计
const loadForumActivityStats = async () => {
  try {
    const response = await adminApi.getForumActivityStats({ timeRange: timeRange.value })
    const data = response.data
    
    // 更新论坛活动数据
    forumActivityData.value.labels = data.labels || []
    forumActivityData.value.datasets[0].data = data.topics || []
    forumActivityData.value.datasets[1].data = data.replies || []
  } catch (error) {
    console.error('加载论坛活动统计失败:', error)
    throw error
  }
}

// 加载热门资源和话题
const loadPopularItems = async () => {
  try {
    // 加载热门资源
    const resourceResponse = await adminApi.getPopularResources({ timeRange: timeRange.value, limit: 5 })
    popularResources.value = resourceResponse.data.resources || []
    
    // 加载热门话题
    const topicResponse = await adminApi.getPopularTopics({ timeRange: timeRange.value, limit: 5 })
    popularTopics.value = topicResponse.data.topics || []
  } catch (error) {
    console.error('加载热门项目失败:', error)
    throw error
  }
}

onMounted(async () => {
  if (!userStore.isLoggedIn || userStore.user?.role !== 'admin') {
    router.push('/')
    return
  }
  
  // 加载所有统计数据
  loadAllStats()
})
</script>