<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">创建新主题</h1>
      <router-link 
        to="/forum"
        class="text-gray-600 hover:text-gray-800 transition duration-300"
      >
        返回论坛
      </router-link>
    </div>

    <div class="bg-white rounded-lg shadow-md p-6">
      <form @submit.prevent="submitTopic">
        <div class="mb-6">
          <label for="title" class="block text-sm font-medium text-gray-700 mb-1">标题</label>
          <InputText 
            id="title"
            v-model="form.title" 
            placeholder="输入主题标题" 
            class="w-full"
            required
          />
        </div>

        <div class="mb-6">
          <label for="category" class="block text-sm font-medium text-gray-700 mb-1">分类</label>
          <Dropdown 
            id="category"
            v-model="form.category" 
            :options="categories" 
            optionLabel="name"
            optionValue="id"
            placeholder="选择分类"
            class="w-full"
            required
          />
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-1">内容</label>
          <Editor 
            v-model="form.content" 
            editorStyle="height: 320px"
          />
        </div>

        <div class="flex justify-end">
          <Button 
            type="submit" 
            label="发布主题" 
            class="p-button-primary"
            :loading="submitting"
          />
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Editor from 'primevue/editor'
import Button from 'primevue/button'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  title: '',
  category: null,
  content: ''
})

const categories = ref([
  { id: 1, name: '前端开发' },
  { id: 2, name: '后端开发' },
  { id: 3, name: '移动开发' },
  { id: 4, name: '数据库' },
  { id: 5, name: '运维部署' },
])

const submitting = ref(false)

const submitTopic = async () => {
  try {
    submitting.value = true
    // TODO: 替换为实际API调用
    console.log('提交主题:', form.value)
    
    // 模拟API调用延迟
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 提交成功后跳转到主题详情页
    router.push('/forum/topics/123')
  } catch (error) {
    console.error('提交主题失败:', error)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
  }
})
</script>