<template>
  <div class="bg-white rounded-lg shadow-md overflow-hidden">
    <!-- 聊天界面头部 -->
    <div class="bg-blue-600 text-white p-4">
      <div class="flex justify-between items-center">
        <h1 class="text-xl font-bold">AI助手</h1>
        <div class="flex space-x-2">
          <Button icon="pi pi-cog" class="p-button-rounded p-button-text p-button-sm" v-tooltip.left="'设置'" />
          <Button icon="pi pi-refresh" class="p-button-rounded p-button-text p-button-sm" v-tooltip.left="'新对话'" @click="startNewChat" />
        </div>
      </div>
    </div>

    <!-- 聊天内容区域 -->
    <div class="flex h-[calc(100vh-12rem)]">
      <!-- 历史记录侧边栏 -->
      <div class="w-64 bg-gray-50 border-r border-gray-200 p-4 hidden md:block">
        <div class="mb-4">
          <Button label="新建对话" icon="pi pi-plus" class="p-button-outlined w-full" @click="startNewChat" />
        </div>
        <div class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2">历史对话</div>
        <div class="space-y-2 max-h-[calc(100vh-16rem)] overflow-y-auto">
          <div 
            v-for="(chat, index) in chatHistory" 
            :key="index"
            class="p-2 rounded-md cursor-pointer hover:bg-gray-200 transition-colors"
            :class="{ 'bg-blue-100': currentChatIndex === index }"
            @click="selectChat(index)"
          >
            <div class="text-sm font-medium truncate">{{ chat.title }}</div>
            <div class="text-xs text-gray-500">{{ formatDate(chat.date) }}</div>
          </div>
        </div>
      </div>

      <!-- 聊天主区域 -->
      <div class="flex-1 flex flex-col">
        <!-- 消息列表 -->
        <div class="flex-1 p-4 overflow-y-auto" ref="messageContainer">
          <div v-if="currentChat.messages.length === 0" class="h-full flex flex-col items-center justify-center text-center p-4">
            <i class="pi pi-comments text-5xl text-gray-300 mb-4"></i>
            <h3 class="text-xl font-bold text-gray-700 mb-2">欢迎使用AI助手</h3>
            <p class="text-gray-500 mb-6 max-w-md">我可以帮助你解答问题、提供建议或者进行日常对话。请在下方输入你的问题。</p>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-2 w-full max-w-lg">
              <Button 
                v-for="(suggestion, index) in suggestions" 
                :key="index"
                :label="suggestion"
                class="p-button-outlined p-button-sm text-left"
                @click="sendMessage(suggestion)"
              />
            </div>
          </div>
          <div v-else>
            <div 
              v-for="(message, index) in currentChat.messages" 
              :key="index"
              class="mb-4 flex"
              :class="{ 'justify-end': message.sender === 'user' }"
            >
              <div 
                class="max-w-[80%] rounded-lg p-3"
                :class="message.sender === 'user' ? 'bg-blue-500 text-white rounded-tr-none' : 'bg-gray-100 text-gray-800 rounded-tl-none'"
              >
                <div class="prose prose-sm" v-html="formatMessage(message.content)"></div>
                <div class="text-xs mt-1 opacity-70 text-right">
                  {{ formatTime(message.timestamp) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="border-t border-gray-200 p-4">
          <div class="flex items-end">
            <div class="flex-1 relative">
              <Textarea 
                v-model="newMessage" 
                placeholder="输入消息..." 
                class="w-full p-3 pr-10 rounded-lg border border-gray-300 focus:border-blue-500 focus:ring focus:ring-blue-200 focus:ring-opacity-50"
                :autoResize="true"
                rows="2"
                @keydown.enter.prevent="sendMessage()"
              />
              <div class="absolute right-2 bottom-2 flex space-x-1">
                <Button icon="pi pi-paperclip" class="p-button-text p-button-rounded p-button-sm" v-tooltip.top="'上传文件'" />
                <Button icon="pi pi-image" class="p-button-text p-button-rounded p-button-sm" v-tooltip.top="'上传图片'" />
              </div>
            </div>
            <Button 
              icon="pi pi-send" 
              class="p-button-rounded ml-2" 
              :disabled="!newMessage.trim()"
              @click="sendMessage()"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick, onMounted } from 'vue'
import Button from 'primevue/button'
import Textarea from 'primevue/textarea'
const messageContainer = ref(null)
const newMessage = ref('')
const currentChatIndex = ref(0)

// 示例建议问题
const suggestions = [
  '你能介绍一下自己吗？',
  '如何学习前端开发？',
  '推荐一些学习资源',
  '如何提高编程效率？'
]

// 模拟聊天历史
const chatHistory = ref([
  {
    id: 1,
    title: '关于前端开发的对话',
    date: new Date(2023, 8, 15),
    messages: [
      {
        sender: 'user',
        content: '你好，我想学习前端开发，有什么建议吗？',
        timestamp: new Date(2023, 8, 15, 10, 30)
      },
      {
        sender: 'assistant',
        content: `学习前端开发的建议：
        <ol>
          <li>掌握HTML、CSS和JavaScript基础</li>
          <li>学习一个现代前端框架，如Vue.js或React</li>
          <li>了解响应式设计和移动优先的概念</li>
          <li>学习版本控制工具如Git</li>
          <li>实践项目，建立个人作品集</li>
        </ol>
        你有特定的学习目标吗？`,
        timestamp: new Date(2023, 8, 15, 10, 31)
      }
    ]
  },
  {
    id: 2,
    title: '数据库设计问题',
    date: new Date(2023, 8, 16),
    messages: [
      {
        sender: 'user',
        content: '在设计用户表时，应该考虑哪些字段？',
        timestamp: new Date(2023, 8, 16, 14, 20)
      },
      {
        sender: 'assistant',
        content: `用户表常见字段包括：
        <ul>
          <li>id: 唯一标识符</li>
          <li>username: 用户名</li>
          <li>email: 电子邮件</li>
          <li>password_hash: 密码哈希</li>
          <li>created_at: 创建时间</li>
          <li>updated_at: 更新时间</li>
          <li>status: 账户状态</li>
          <li>avatar: 头像URL</li>
        </ul>
        根据应用需求，可能还需要添加其他字段。`,
        timestamp: new Date(2023, 8, 16, 14, 22)
      }
    ]
  },
  {
    id: 3,
    title: '新对话',
    date: new Date(),
    messages: []
  }
])

// 当前聊天
const currentChat = computed(() => chatHistory.value[currentChatIndex.value])

// 格式化日期
function formatDate(date) {
  return new Intl.DateTimeFormat('zh-CN', {
    month: 'short',
    day: 'numeric'
  }).format(date)
}

// 格式化时间
function formatTime(timestamp) {
  return new Intl.DateTimeFormat('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  }).format(timestamp)
}

// 格式化消息内容（支持简单的Markdown和HTML）
function formatMessage(content) {
  return content
}

// 发送消息
function sendMessage(content = null) {
  const messageText = content || newMessage.value.trim()
  if (!messageText) return
  
  // 添加用户消息
  currentChat.value.messages.push({
    sender: 'user',
    content: messageText,
    timestamp: new Date()
  })
  
  // 清空输入框
  newMessage.value = ''
  
  // 滚动到底部
  scrollToBottom()
  
  // 模拟AI回复
  setTimeout(() => {
    let response
    
    if (messageText.includes('前端') || messageText.includes('学习')) {
      response = `前端开发是一个不断发展的领域，建议：
      <ul>
        <li>持续学习新技术和框架</li>
        <li>参与开源项目积累经验</li>
        <li>关注行业动态和最佳实践</li>
        <li>构建个人项目展示能力</li>
      </ul>
      你对哪个方向更感兴趣？`
    } else if (messageText.includes('Vue') || messageText.includes('框架')) {
      response = `Vue.js是一个流行的前端框架，特点包括：
      <ul>
        <li>易学易用，渐进式架构</li>
        <li>响应式数据绑定</li>
        <li>组件化开发模式</li>
        <li>强大的生态系统</li>
      </ul>
      推荐先从官方文档开始学习。`
    } else {
      response = '我理解你的问题。你能提供更多细节，这样我可以给你更具体的帮助吗？'
    }
    
    // 添加AI回复
    currentChat.value.messages.push({
      sender: 'assistant',
      content: response,
      timestamp: new Date()
    })
    
    // 更新聊天标题（如果是新对话）
    if (currentChat.value.title === '新对话' && currentChat.value.messages.length >= 2) {
      currentChat.value.title = messageText.substring(0, 20) + (messageText.length > 20 ? '...' : '')
    }
    
    // 滚动到底部
    scrollToBottom()
  }, 1000)
}

// 滚动到底部
function scrollToBottom() {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }
  })
}

// 选择聊天
function selectChat(index) {
  currentChatIndex.value = index
  scrollToBottom()
}

// 开始新对话
function startNewChat() {
  chatHistory.value.unshift({
    id: Date.now(),
    title: '新对话',
    date: new Date(),
    messages: []
  })
  currentChatIndex.value = 0
}

onMounted(() => {
  scrollToBottom()
})
</script>