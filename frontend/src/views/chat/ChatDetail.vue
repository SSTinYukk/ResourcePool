<template>
  <div class="container mx-auto py-4 px-4 h-screen flex flex-col">
    <!-- 头部 -->
    <div class="flex justify-between items-center mb-4">
      <div class="flex items-center">
        <Button icon="pi pi-arrow-left" class="p-button-text mr-2" @click="goBack" />
        <h1 class="text-xl font-bold">{{ sessionTitle }}</h1>
      </div>
      <Button icon="pi pi-ellipsis-v" class="p-button-text p-button-rounded" @click="toggleMenu" />
      <Menu ref="menu" :model="menuItems" :popup="true" />
    </div>

    <!-- 消息列表 -->
    <div class="flex-grow overflow-y-auto mb-6 bg-gray-50 rounded-lg p-4" ref="messageContainer">
      <div v-if="loading" class="flex justify-center py-8">
        <i class="pi pi-spin pi-spinner text-4xl text-blue-500"></i>
      </div>

      <div v-else-if="error" class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4">
        <p>{{ error }}</p>
        <Button label="重试" icon="pi pi-refresh" class="mt-2" @click="fetchMessages" />
      </div>

      <div v-else-if="messages.length === 0" class="text-center py-8 text-gray-500">
        <p>开始与AI助理对话吧！</p>
      </div>

      <div v-else class="space-y-4">
        <div 
          v-for="(message, index) in messages" 
          :key="index"
          :class="{
            'flex': true,
            'justify-end': message.role === 'user',
            'justify-start': message.role === 'assistant'
          }"
        >
          <div 
            :class="{
              'max-w-3/4 rounded-xl p-4': true,
              'bg-gradient-to-r from-blue-500 to-blue-600 text-white shadow-md': message.role === 'user',
              'bg-gradient-to-r from-gray-50 to-white border border-gray-100 shadow-sm': message.role === 'assistant'
            }"
          >
            <div class="whitespace-pre-wrap" v-html="parseMarkdown(message.content)"></div>
            <div 
              :class="{
                'text-xs mt-2 opacity-80': true,
                'text-blue-100': message.role === 'user',
                'text-gray-500': message.role === 'assistant'
              }"
            >
              {{ formatTime(message.createdAt) }}
            </div>
          </div>
        </div>
      </div>

      <div v-if="sending" class="flex justify-start mt-4">
        <div class="bg-white border border-gray-100 rounded-xl p-4 shadow-sm flex items-center">
          <i class="pi pi-spin pi-spinner mr-3 text-blue-500"></i>
          <span class="text-gray-600">AI正在思考中...</span>
        </div>
      </div>
    </div>

    <!-- 输入框 -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4 mt-2 mb-14">
      <div class="flex items-center space-x-2">
        <InputText 
          v-model="newMessage" 
          placeholder="输入消息..." 
          class="flex-grow rounded-full border-gray-200" 
          @keyup.enter="sendMessage"
        />
        <Button 
          icon="pi pi-send" 
          :disabled="!newMessage.trim() || sending" 
          @click="sendMessage" 
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Menu from 'primevue/menu'
import { chatApi } from '@/api/chat'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const menu = ref(null)
const messageContainer = ref(null)

const sessionId = computed(() => route.params.id)
const sessionTitle = ref('AI对话')
const messages = ref([])
const newMessage = ref('')
const loading = ref(true)
const sending = ref(false)
const error = ref(null)

// 菜单选项
const menuItems = [
  {
    label: '清空对话',
    icon: 'pi pi-trash',
    command: () => clearMessages()
  },
  {
    label: '返回列表',
    icon: 'pi pi-list',
    command: () => router.push('/chat')
  }
]

// 获取消息列表
const fetchMessages = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response = await chatApi.getSessionMessages(sessionId.value)
    messages.value = response.data.messages
    sessionTitle.value = response.data.title || 'AI对话'
  } catch (err) {
    console.error('获取消息失败:', err)
    error.value = '获取消息失败，请稍后再试'
  } finally {
    loading.value = false
    scrollToBottom()
  }
}

// 发送消息
const sendMessage = async () => {
  if (!newMessage.value.trim() || sending.value) return
  
  const messageContent = newMessage.value
  newMessage.value = ''
  sending.value = true
  
  // 先添加用户消息到列表
  messages.value.push({
    role: 'user',
    content: messageContent,
    createdAt: new Date().toISOString()
  })
  
  scrollToBottom()
  
  try {
    const response = await chatApi.sendMessage({
        session_id: Number(sessionId.value),
        content: messageContent
    })
    console.log("DEBUG:",sessionId.value,messageContent)
    
    // 添加AI回复到列表
    messages.value.push({
      role: 'assistant',
      content: response.data.content,
      createdAt: response.data.createdAt || new Date().toISOString()
    })
    
    // 如果是新会话，更新标题
    if (!sessionTitle.value || sessionTitle.value === 'AI对话') {
      sessionTitle.value = response.data.sessionTitle || 'AI对话'
    }
  } catch (err) {
    console.error('发送消息失败:', err)
    toast.add({
      severity: 'error',
      summary: '发送失败',
      detail: '无法发送消息，请稍后再试',
      life: 3000
    })
  } finally {
    sending.value = false
    scrollToBottom()
  }
}

// 清空消息
// 解析Markdown内容
const parseMarkdown = (content) => {
  return DOMPurify.sanitize(marked.parse(content))
}

const clearMessages = async () => {
  try {
    // 这里可以添加API调用来清空服务器端的消息
    // 暂时只清空本地消息
    messages.value = []
    toast.add({
      severity: 'success',
      summary: '已清空',
      detail: '对话已清空',
      life: 3000
    })
  } catch (err) {
    console.error('清空消息失败:', err)
    toast.add({
      severity: 'error',
      summary: '操作失败',
      detail: '无法清空消息，请稍后再试',
      life: 3000
    })
  }
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }
  })
}

// 显示菜单
const toggleMenu = (event) => {
  menu.value.toggle(event)
}

// 返回上一页
const goBack = () => {
  router.push('/chat')
}

// 格式化时间
const formatTime = (dateString) => {
  if (!dateString) return ''
  
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

onMounted(() => {
  fetchMessages()
})
</script>