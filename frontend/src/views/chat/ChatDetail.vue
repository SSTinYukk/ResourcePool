<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">聊天详情</h1>
      <Button 
        label="返回" 
        icon="pi pi-arrow-left" 
        class="p-button-secondary"
        @click="goBack"
      />
    </div>

    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <div class="mb-4">
        <h2 class="text-xl font-semibold">{{ chatTitle }}</h2>
        <p class="text-gray-500">参与者: {{ participants }}</p>
      </div>

      <div class="border rounded-lg p-4 mb-6 h-96 overflow-y-auto">
        <div 
          v-for="(message, index) in messages" 
          :key="index" 
          class="mb-4"
          :class="{ 'text-right': message.isMe }"
        >
          <div 
            class="inline-block px-4 py-2 rounded-lg"
            :class="{ 
              'bg-blue-100 text-blue-800': message.isMe, 
              'bg-gray-100 text-gray-800': !message.isMe 
            }"
          >
            {{ message.content }}
          </div>
          <div class="text-xs text-gray-500 mt-1">
            {{ message.time }}
          </div>
        </div>
      </div>

      <div class="flex gap-2">
        <InputText 
          v-model="newMessage" 
          placeholder="输入消息..." 
          class="flex-grow"
          @keyup.enter="sendMessage"
        />
        <Button 
          label="发送" 
          icon="pi pi-send" 
          class="p-button-primary"
          @click="sendMessage"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const chatTitle = ref('')
const participants = ref('')
const messages = ref([])
const newMessage = ref('')

const fetchChatDetails = async () => {
  // TODO: Replace with actual API call
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 300))
  
  chatTitle.value = '项目讨论'
  participants.value = '张三, 李四'
  messages.value = [
    { 
      content: '大家好，我们今天讨论项目进度', 
      isMe: false, 
      time: '10:30 AM' 
    },
    { 
      content: '我已经完成了前端部分', 
      isMe: true, 
      time: '10:32 AM' 
    },
    { 
      content: '后端API还需要一些调整', 
      isMe: false, 
      time: '10:35 AM' 
    }
  ]
}

const sendMessage = () => {
  if (!newMessage.value.trim()) return
  
  messages.value.push({
    content: newMessage.value,
    isMe: true,
    time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  })
  
  newMessage.value = ''
}

const goBack = () => {
  router.push('/chat')
}

onMounted(() => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
  } else {
    fetchChatDetails()
  }
})
</script>