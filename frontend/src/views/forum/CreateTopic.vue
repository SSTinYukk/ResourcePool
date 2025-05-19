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
import { useToast } from 'primevue/usetoast'
import { forumApi } from '@/api/forum'

const router = useRouter()
const userStore = useUserStore()
const toast = useToast()

const form = ref({
  title: '',
  category: null,
  content: ''
})

const categories = ref([])

const loadCategories = async () => {
  try {
    categories.value = [
      { id: 1, name: '技术讨论' },
      { id: 2, name: '学习资源' },
      { id: 3, name: '经验分享' },
      { id: 4, name: '求助问答' },
      { id: 5, name: '活动公告' }
    ]
  } catch (error) {
    console.error('获取分类失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '获取分类失败', life: 3000 })
  }
}

const submitting = ref(false)

const submitTopic = async () => {
  try {
    submitting.value = true
    
    // 表单验证
    if (!form.value.title || !form.value.category || !form.value.content) {
      throw new Error('请填写所有必填字段')
    }
    
    // 调用API创建主题
    const response = await forumApi.createTopic({
      title: form.value.title,
      category_id: form.value.category,
      content: form.value.content,
      author_id: userStore.user.id
    })
    
    // 提交成功后跳转到新创建的主题
    if (response && response.data && response.data.id) {
      router.push(`/forum/topic/${response.data.id}`)
    } else {
      throw new Error('创建主题失败，请稍后重试')
    }
  } catch (error) {
    console.error('提交主题失败:', error)
    // 显示错误提示
    if (error.response && error.response.data && error.response.data.message) {
      toast.add({ severity: 'error', summary: '错误', detail: error.response.data.message, life: 3000 })
    } else if (error.message) {
      toast.add({ severity: 'error', summary: '错误', detail: error.message, life: 3000 })
    } else {
      toast.add({ severity: 'error', summary: '错误', detail: '提交主题失败，请稍后重试', life: 3000 })
    }
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
  } else {
    loadCategories()
  }
})
</script>
