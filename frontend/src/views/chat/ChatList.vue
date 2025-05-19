<template>
  <div class="container mx-auto py-8 px-4">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">AI助理</h1>
      <Button label="新建会话" icon="pi pi-plus" @click="createNewSession" />
    </div>

    <div v-if="loading" class="flex justify-center py-8">
      <i class="pi pi-spin pi-spinner text-4xl text-blue-500"></i>
    </div>

    <div v-else-if="error" class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4">
      <p>{{ error }}</p>
      <Button label="重试" icon="pi pi-refresh" class="mt-2" @click="fetchSessions" />
    </div>

    <div v-else-if="sessions.length === 0" class="text-center py-8">
      <div class="text-gray-500 mb-4">暂无聊天会话</div>
      <Button label="开始新对话" icon="pi pi-comments" @click="createNewSession" />
    </div>

    <div v-else class="grid grid-cols-1 gap-4">
      <div 
        v-for="session in sessions" 
        :key="session.id"
        class="bg-white rounded-lg shadow-md p-4 hover:shadow-lg transition-shadow cursor-pointer flex justify-between items-center"
        @click="goToChat(session.id)"
      >
        <div>
          <div class="font-medium">{{ session.title || '新对话' }}</div>
          <div class="text-sm text-gray-500">{{ formatDate(session.createdAt) }}</div>
        </div>
        <div class="flex items-center">
          <Button 
            icon="pi pi-trash" 
            class="p-button-text p-button-rounded p-button-danger" 
            @click.stop="confirmDeleteSession(session.id)" 
          />
          <i class="pi pi-angle-right ml-2 text-gray-400"></i>
        </div>
      </div>
    </div>

    <ConfirmDialog>
      <template #message>
        <div class="p-4">
          <p>确定要删除这个会话吗？此操作不可恢复。</p>
        </div>
      </template>
    </ConfirmDialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import ConfirmDialog from 'primevue/confirmdialog'
import { chatApi } from '@/api/chat'

const router = useRouter()
const confirm = useConfirm()
const toast = useToast()

const sessions = ref([])
const loading = ref(true)
const error = ref(null)

// 获取会话列表
const fetchSessions = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response = await chatApi.getSessions()
    sessions.value = response.data.sessions || response.data
    console.log('获取到的会话列表:', sessions.value)
    
  } catch (err) {
    console.error('获取会话列表失败:', err)
    error.value = '获取会话列表失败，请稍后再试'
  } finally {
    loading.value = false
  }
}

// 创建新会话
const createNewSession = async () => {
  try {
    const response = await chatApi.createSession()
    router.push(`/chat/${response.data.id}`)
  } catch (err) {
    console.error('创建会话失败:', err)
    toast.add({
      severity: 'error',
      summary: '创建失败',
      detail: '无法创建新会话，请稍后再试',
      life: 3000
    })
  }
}

// 确认删除会话
const confirmDeleteSession = (sessionId) => {
  confirm.require({
    message: '确定要删除这个会话吗？',
    header: '删除确认',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: () => deleteSession(sessionId),
    reject: () => {}
  })
}

// 删除会话
const deleteSession = async (sessionId) => {
  try {
    await chatApi.deleteSession(sessionId)
    sessions.value = sessions.value.filter(session => session.id !== sessionId)
    toast.add({
      severity: 'success',
      summary: '删除成功',
      detail: '会话已成功删除',
      life: 3000
    })
  } catch (err) {
    console.error('删除会话失败:', err)
    toast.add({
      severity: 'error',
      summary: '删除失败',
      detail: '无法删除会话，请稍后再试',
      life: 3000
    })
  }
}

// 跳转到聊天页面
const goToChat = (sessionId) => {
  router.push(`/chat/${sessionId}`)
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '未知时间'
  try {
    const date = new Date(dateString)
    return new Intl.DateTimeFormat('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date)
  } catch (e) {
    console.error('日期格式化错误:', e)
    return '无效时间'
  }
}

onMounted(() => {
  fetchSessions()
})
</script>