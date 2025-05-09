<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <!-- 主题标题和操作 -->
    <div class="mb-6">
      <div class="flex justify-between items-start">
        <div>
          <h1 class="text-2xl font-bold text-gray-800 mb-2">{{ topic.title }}</h1>
          <div class="flex items-center text-sm text-gray-500">
            <span>{{ formatDate(topic.createTime) }}</span>
            <span class="mx-2">•</span>
            <span>{{ topic.views }} 次查看</span>
            <span class="mx-2">•</span>
            <span>{{ topic.replies.length }} 条回复</span>
            <Tag :value="getCategoryName(topic.category)" class="ml-2" severity="info" />
          </div>
        </div>
        <div class="flex space-x-2">
          <Button icon="pi pi-bookmark" class="p-button-text p-button-rounded" v-tooltip.top="'收藏'" />
          <Button icon="pi pi-share-alt" class="p-button-text p-button-rounded" v-tooltip.top="'分享'" />
          <Button icon="pi pi-ellipsis-v" class="p-button-text p-button-rounded" @click="toggleMenu($event)" />
          <Menu ref="menu" :model="menuItems" :popup="true" />
        </div>
      </div>
    </div>

    <!-- 主题内容 -->
    <div class="border-t border-b border-gray-200 py-6">
      <div class="flex mb-4">
        <div class="mr-4">
          <Avatar :image="topic.authorAvatar" :label="!topic.authorAvatar ? topic.author.charAt(0).toUpperCase() : undefined" shape="circle" size="large" />
          <div class="text-center mt-2 text-sm font-medium">{{ topic.author }}</div>
        </div>
        <div class="flex-1">
          <div class="prose max-w-none" v-html="topic.content"></div>
        </div>
      </div>
      <div class="flex justify-between items-center mt-4">
        <div class="flex space-x-2">
          <Button icon="pi pi-thumbs-up" :label="`${topic.likes}`" class="p-button-text p-button-sm" />
          <Button icon="pi pi-thumbs-down" :label="`${topic.dislikes}`" class="p-button-text p-button-sm" />
        </div>
        <Button icon="pi pi-flag" class="p-button-text p-button-sm" label="举报" />
      </div>
    </div>

    <!-- 回复列表 -->
    <div class="mt-8">
      <h2 class="text-xl font-bold text-gray-800 mb-4">{{ topic.replies.length }} 条回复</h2>
      
      <div v-for="(reply, index) in topic.replies" :key="index" class="border-b border-gray-200 py-4 last:border-b-0">
        <div class="flex">
          <div class="mr-4">
            <Avatar :image="reply.authorAvatar" :label="!reply.authorAvatar ? reply.author.charAt(0).toUpperCase() : undefined" shape="circle" />
            <div class="text-center mt-2 text-sm font-medium">{{ reply.author }}</div>
          </div>
          <div class="flex-1">
            <div class="flex justify-between items-start mb-2">
              <div class="text-sm text-gray-500">
                <span>#{{ index + 1 }}</span>
                <span class="mx-2">•</span>
                <span>{{ formatDate(reply.createTime) }}</span>
              </div>
              <Button icon="pi pi-ellipsis-h" class="p-button-text p-button-rounded p-button-sm" />
            </div>
            <div class="prose max-w-none" v-html="reply.content"></div>
            <div class="flex justify-between items-center mt-4">
              <div class="flex space-x-2">
                <Button icon="pi pi-thumbs-up" :label="`${reply.likes}`" class="p-button-text p-button-sm" />
                <Button icon="pi pi-thumbs-down" :label="`${reply.dislikes}`" class="p-button-text p-button-sm" />
              </div>
              <Button icon="pi pi-reply" class="p-button-text p-button-sm" label="回复" @click="replyTo(reply)" />
            </div>
          </div>
        </div>
      </div>

      <!-- 无回复时显示 -->
      <div v-if="topic.replies.length === 0" class="text-center py-8">
        <i class="pi pi-comments text-5xl text-gray-300 mb-4"></i>
        <p class="text-gray-500">暂无回复，成为第一个回复者吧！</p>
      </div>
    </div>

    <!-- 回复编辑器 -->
    <div class="mt-8">
      <h3 class="text-lg font-bold text-gray-800 mb-4">发表回复</h3>
      <div v-if="replyingTo" class="bg-blue-50 p-3 rounded-md mb-4 flex justify-between items-center">
        <div>回复 <span class="font-medium">{{ replyingTo.author }}</span>: {{ replyingTo.content.substring(0, 50) }}{{ replyingTo.content.length > 50 ? '...' : '' }}</div>
        <Button icon="pi pi-times" class="p-button-text p-button-rounded p-button-sm" @click="cancelReply" />
      </div>
      <Editor v-model="newReply" editorStyle="height: 200px" />
      <div class="flex justify-end mt-4">
        <Button label="发表回复" icon="pi pi-send" @click="submitReply" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Tag from 'primevue/tag'
import Menu from 'primevue/menu'
import Editor from 'primevue/editor'

const route = useRoute()
const router = useRouter()
const menu = ref(null)
const newReply = ref('')
const replyingTo = ref(null)

// 模拟数据 - 实际应用中应从API获取
const categories = [
  { id: 1, name: '全部' },
  { id: 2, name: '技术讨论' },
  { id: 3, name: '学习资源' },
  { id: 4, name: '经验分享' },
  { id: 5, name: '求助问答' },
  { id: 6, name: '活动公告' }
]

// 模拟主题数据
const topic = ref({
  id: 1,
  title: '如何高效学习Vue 3？',
  author: '张三',
  authorAvatar: null,
  category: 3,
  createTime: new Date(2023, 8, 15, 10, 30),
  views: 356,
  likes: 42,
  dislikes: 3,
  content: `
    <p>大家好，我最近开始学习Vue 3，想请教一下有什么高效的学习方法和资源推荐？</p>
    <p>我已经有了一些JavaScript和前端开发的基础，但是Vue 3的Composition API和一些新特性对我来说还是比较陌生的。</p>
    <p>目前我找到了以下资源：</p>
    <ul>
      <li>Vue 3官方文档</li>
      <li>Vue Mastery的视频教程</li>
      <li>一些GitHub上的开源项目</li>
    </ul>
    <p>请问大家有没有其他推荐的学习资源或者学习路线图？感谢！</p>
  `,
  replies: [
    {
      author: '李四',
      authorAvatar: null,
      createTime: new Date(2023, 8, 15, 11, 15),
      content: `
        <p>我觉得学习Vue 3最好的方式是先通读官方文档，然后做一些小项目来实践。</p>
        <p>特别推荐Vue 3官方文档中的Composition API部分，写得非常详细。</p>
        <p>另外，可以看看这个YouTube频道：Vue School，他们有很多高质量的Vue 3教程。</p>
      `,
      likes: 15,
      dislikes: 0
    },
    {
      author: '王五',
      authorAvatar: null,
      createTime: new Date(2023, 8, 15, 14, 30),
      content: `
        <p>我建议你可以先了解一下Vue 3的核心概念，比如响应式系统、组件化、指令等。</p>
        <p>然后再深入学习Composition API，这是Vue 3最大的变化之一。</p>
        <p>我个人比较喜欢通过实战项目学习，可以尝试重构一个你熟悉的项目，或者跟着教程做一个新项目。</p>
        <p>推荐一个资源：<a href="https://github.com/vuejs/awesome-vue" target="_blank">awesome-vue</a>，里面有很多Vue相关的资源。</p>
      `,
      likes: 8,
      dislikes: 1
    }
  ]
})

// 菜单项
const menuItems = [
  {
    label: '编辑主题',
    icon: 'pi pi-pencil',
    command: () => router.push(`/forum/edit/${topic.value.id}`)
  },
  {
    label: '删除主题',
    icon: 'pi pi-trash',
    command: () => confirmDelete()
  }
]

// 获取分类名称
function getCategoryName(categoryId) {
  const category = categories.find(c => c.id === categoryId)
  return category ? category.name : '未分类'
}

// 格式化日期
function formatDate(date) {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

// 显示菜单
function toggleMenu(event) {
  menu.value.toggle(event)
}

// 确认删除
function confirmDelete() {
  // 实际应用中应使用确认对话框
  if (confirm('确定要删除这个主题吗？')) {
    // 调用API删除主题
    router.push('/forum')
  }
}

// 回复某人
function replyTo(reply) {
  replyingTo.value = reply
  // 滚动到编辑器
  document.querySelector('.p-editor-container').scrollIntoView({ behavior: 'smooth' })
}

// 取消回复
function cancelReply() {
  replyingTo.value = null
}

// 提交回复
function submitReply() {
  if (!newReply.value.trim()) {
    alert('回复内容不能为空')
    return
  }
  
  // 实际应用中应调用API提交回复
  topic.value.replies.push({
    author: '当前用户',
    authorAvatar: null,
    createTime: new Date(),
    content: newReply.value,
    likes: 0,
    dislikes: 0
  })
  
  // 清空编辑器
  newReply.value = ''
  replyingTo.value = null
}

onMounted(() => {
  // 实际应用中应根据路由参数获取主题详情
  const topicId = route.params.id
  console.log('加载主题ID:', topicId)
  // 模拟API调用
  // loadTopic(topicId)
})
</script>