<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">管理员仪表盘</h1>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="bg-white rounded-lg shadow-md p-6">
        <h2 class="text-lg font-semibold text-gray-700 mb-2">用户总数</h2>
        <p class="text-3xl font-bold text-blue-600">{{ stats.user_count || 0 }}</p>
        <p class="text-sm text-gray-500 mt-2">今日新增: {{ stats.today_user_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-md p-6">
        <h2 class="text-lg font-semibold text-gray-700 mb-2">资源总数</h2>
        <p class="text-3xl font-bold text-green-600">{{ stats.resource_count || 0 }}</p>
        <p class="text-sm text-gray-500 mt-2">待审核: {{ stats.pending_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-md p-6">
        <h2 class="text-lg font-semibold text-gray-700 mb-2">系统状态</h2>
        <p class="text-3xl font-bold text-purple-600">正常</p>
        <p class="text-sm text-gray-500 mt-2">最近更新: {{ formatDate(new Date()) }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
      <div class="bg-white rounded-lg shadow-md p-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-lg font-semibold text-gray-700">最近注册用户</h2>
          <Button label="查看全部" icon="pi pi-external-link" text @click="navigateTo('admin-users')" />
        </div>
        <DataTable 
          :value="recentUsers" 
          :loading="loadingUsers"
          :paginator="false"
          class="p-datatable-sm"
          responsiveLayout="scroll"
          stripedRows
        >
          <Column field="username" header="用户名" :sortable="true">
            <template #body="{data}">
              <div class="flex items-center">
                <Avatar :image="data.avatar || '/default-avatar.png'" shape="circle" class="mr-2" />
                <span>{{ data.username }}</span>
              </div>
            </template>
          </Column>
          <Column field="email" header="邮箱" :sortable="true"></Column>
          <Column field="created_at" header="注册时间" :sortable="true">
            <template #body="{data}">
              {{ formatDate(data.created_at) }}
            </template>
          </Column>
        </DataTable>
      </div>

      <div class="bg-white rounded-lg shadow-md p-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-lg font-semibold text-gray-700">待审核资源</h2>
          <Button label="查看全部" icon="pi pi-external-link" text @click="navigateTo('admin-resources')" />
        </div>
        <DataTable 
          :value="pendingResources" 
          :loading="loadingResources"
          :paginator="false"
          class="p-datatable-sm"
          responsiveLayout="scroll"
          stripedRows
        >
          <Column field="title" header="标题" :sortable="true"></Column>
          <Column field="user.username" header="上传者" :sortable="true"></Column>
          <Column field="created_at" header="上传时间" :sortable="true">
            <template #body="{data}">
              {{ formatDate(data.created_at) }}
            </template>
          </Column>
          <Column header="操作">
            <template #body="{data}">
              <Button icon="pi pi-eye" text rounded @click="reviewResource(data.id)" />
            </template>
          </Column>
        </DataTable>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { adminApi } from '@/api/admin'
import { format } from 'date-fns'
import { zhCN } from 'date-fns/locale'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Avatar from 'primevue/avatar'

const router = useRouter()
const userStore = useUserStore()

const stats = ref({
  user_count: 0,
  resource_count: 0,
  pending_count: 0,
  today_user_count: 0
})

const recentUsers = ref([])
const pendingResources = ref([])
const loadingUsers = ref(true)
const loadingResources = ref(true)

// 格式化日期
const formatDate = (date) => {
  if (!date) return ''
  return format(new Date(date), 'yyyy-MM-dd HH:mm', { locale: zhCN })
}

// 获取统计数据
const fetchStats = async () => {
  try {
    const response = await adminApi.getStats()
    stats.value = response.data
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 获取最近用户
const fetchRecentUsers = async () => {
  loadingUsers.value = true
  try {
    const response = await adminApi.getUsers({ page: 1, pageSize: 5, sort: 'created_at:desc' })
    recentUsers.value = response.data.users || []
  } catch (error) {
    console.error('获取最近用户失败:', error)
  } finally {
    loadingUsers.value = false
  }
}

// 获取待审核资源
const fetchPendingResources = async () => {
  loadingResources.value = true
  try {
    const response = await adminApi.getResources({ page: 1, pageSize: 5, status: 'pending' })
    pendingResources.value = response.data.resources || []
  } catch (error) {
    console.error('获取待审核资源失败:', error)
  } finally {
    loadingResources.value = false
  }
}

// 跳转到资源审核页面
const reviewResource = (id) => {
  router.push(`/admin/resources?id=${id}`)
}

// 导航到其他管理页面
const navigateTo = (routeName) => {
  router.push({ name: routeName })
}

// 获取所有数据
const fetchAdminData = async () => {
  await Promise.all([
    fetchStats(),
    fetchRecentUsers(),
    fetchPendingResources()
  ])
}

onMounted(() => {
  if (!userStore.isLoggedIn || userStore.user?.role !== 'admin') {
    router.push('/')
  } else {
    fetchAdminData()
  }
})
</script>