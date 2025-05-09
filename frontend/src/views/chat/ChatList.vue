<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">聊天列表</h1>
      <Button 
        label="新建聊天" 
        icon="pi pi-plus" 
        class="p-button-primary"
        @click="createNewChat"
      />
    </div>

    <div class="bg-white rounded-lg shadow-md p-6">
      <div class="mb-4">
        <InputText 
          v-model="searchQuery" 
          placeholder="搜索聊天..." 
          class="w-full"
        />
      </div>

      <DataTable 
        :value="filteredChats" 
        :paginator="true" 
        :rows="10"
        :loading="loading"
        paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport"
        currentPageReportTemplate="显示 {first} 到 {last} 共 {totalRecords} 条"
      >
        <Column field="title" header="聊天名称">
          <template #body="{data}">
            <router-link 
              :to="`/chat/${data.id}`" 
              class="text-blue-600 hover:underline"
            >
              {{ data.title }}
            </router-link>
          </template>
        </Column>
        <Column field="participants" header="参与者" />
        <Column field="lastMessage" header="最后消息" />
        <Column field="time" header="时间" />
      </DataTable>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const searchQuery = ref('')
const chats = ref([])

const filteredChats = computed(() => {
  if (!searchQuery.value) return chats.value
  return chats.value.filter(chat => 
    chat.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    chat.participants.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const fetchChats = async () => {
  try {
    loading.value = true
    // TODO: Replace with actual API call
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 500))
    chats.value = [
      { id: 1, title: '项目讨论', participants: '张三, 李四', lastMessage: '我们明天开会', time: '10分钟前' },
      { id: 2, title: '技术支持', participants: '王五', lastMessage: '问题已解决', time: '1小时前' },
      { id: 3, title: '团队会议', participants: '全体成员', lastMessage: '下周计划', time: '昨天' },
    ]
  } catch (error) {
    console.error('获取聊天列表失败:', error)
  } finally {
    loading.value = false
  }
}

const createNewChat = () => {
  router.push('/chat/new')
}

onMounted(() => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
  } else {
    fetchChats()
  }
})
</script>