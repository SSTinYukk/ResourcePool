<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">Admin Dashboard</h1>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="bg-white rounded-lg shadow-md p-6">
        <h2 class="text-lg font-semibold text-gray-700 mb-2">Total Users</h2>
        <p class="text-3xl font-bold text-blue-600">{{ userCount }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-md p-6">
        <h2 class="text-lg font-semibold text-gray-700 mb-2">Total Resources</h2>
        <p class="text-3xl font-bold text-green-600">{{ resourceCount }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-md p-6">
        <h2 class="text-lg font-semibold text-gray-700 mb-2">Recent Activity</h2>
        <p class="text-3xl font-bold text-purple-600">{{ activityCount }}</p>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow-md p-6">
      <h2 class="text-lg font-semibold text-gray-700 mb-4">Recent Users</h2>
      <DataTable 
        :value="recentUsers" 
        :paginator="true" 
        :rows="5"
        responsiveLayout="scroll"
      >
        <Column field="name" header="Name" :sortable="true"></Column>
        <Column field="email" header="Email" :sortable="true"></Column>
        <Column field="joinDate" header="Join Date" :sortable="true"></Column>
      </DataTable>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const router = useRouter()
const userStore = useUserStore()

const userCount = ref(0)
const resourceCount = ref(0)
const activityCount = ref(0)
const recentUsers = ref([])

const fetchAdminData = async () => {
  // TODO: Replace with actual API calls
  await new Promise(resolve => setTimeout(resolve, 300))
  
  userCount.value = 125
  resourceCount.value = 342
  activityCount.value = 28
  recentUsers.value = [
    { id: 1, name: 'User 1', email: 'user1@example.com', joinDate: '2023-05-10' },
    { id: 2, name: 'User 2', email: 'user2@example.com', joinDate: '2023-05-09' },
    { id: 3, name: 'User 3', email: 'user3@example.com', joinDate: '2023-05-08' },
    { id: 4, name: 'User 4', email: 'user4@example.com', joinDate: '2023-05-07' },
    { id: 5, name: 'User 5', email: 'user5@example.com', joinDate: '2023-05-06' },
  ]
}

onMounted(() => {
  if (!userStore.isLoggedIn || !userStore.isAdmin) {
    router.push('/')
  } else {
    fetchAdminData()
  }
})
</script>