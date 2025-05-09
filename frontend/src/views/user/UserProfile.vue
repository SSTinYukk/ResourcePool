<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">用户资料</h1>
      <Button 
        label="返回" 
        icon="pi pi-arrow-left" 
        class="p-button-secondary"
        @click="goBack"
      />
    </div>

    <div class="bg-white rounded-lg shadow-md p-6">
      <div class="flex items-center mb-6">
        <Avatar 
          :image="user.avatar" 
          size="xlarge" 
          shape="circle" 
          class="mr-4"
        />
        <div>
          <h2 class="text-xl font-semibold">{{ user.name }}</h2>
          <p class="text-gray-500">{{ user.email }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="space-y-4">
          <h3 class="text-lg font-medium">基本信息</h3>
          <div class="space-y-2">
            <div>
              <label class="block text-sm font-medium text-gray-700">用户名</label>
              <InputText v-model="user.username" class="w-full" disabled />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700">注册时间</label>
              <InputText v-model="user.createdAt" class="w-full" disabled />
            </div>
          </div>
        </div>

        <div class="space-y-4">
          <h3 class="text-lg font-medium">安全设置</h3>
          <Button 
            label="修改密码" 
            icon="pi pi-lock" 
            class="p-button-outlined"
            @click="showChangePasswordDialog = true"
          />
        </div>
      </div>
    </div>

    <Dialog 
      v-model:visible="showChangePasswordDialog" 
      header="修改密码" 
      :modal="true"
      class="w-full max-w-md"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">当前密码</label>
          <Password v-model="password.current" toggleMask class="w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">新密码</label>
          <Password v-model="password.new" toggleMask class="w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">确认新密码</label>
          <Password v-model="password.confirm" toggleMask class="w-full" />
        </div>
      </div>

      <template #footer>
        <Button 
          label="取消" 
          icon="pi pi-times" 
          class="p-button-text"
          @click="showChangePasswordDialog = false"
        />
        <Button 
          label="保存" 
          icon="pi pi-check" 
          class="p-button-primary"
          @click="changePassword"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Avatar from 'primevue/avatar'
import Dialog from 'primevue/dialog'
import Password from 'primevue/password'

const router = useRouter()
const userStore = useUserStore()

const user = ref({
  name: '',
  email: '',
  username: '',
  avatar: '',
  createdAt: ''
})

const showChangePasswordDialog = ref(false)
const password = ref({
  current: '',
  new: '',
  confirm: ''
})

const fetchUserProfile = async () => {
  // TODO: Replace with actual API call
  await new Promise(resolve => setTimeout(resolve, 300))
  
  user.value = {
    name: '张三',
    email: 'zhangsan@example.com',
    username: 'zhangsan',
    avatar: '',
    createdAt: '2023-01-15 10:30:45'
  }
}

const changePassword = async () => {
  // TODO: Implement password change logic
  showChangePasswordDialog.value = false
  password.value = { current: '', new: '', confirm: '' }
}

const goBack = () => {
  router.push('/')
}

onMounted(() => {
  if (!userStore.isLoggedIn) {
    router.push('/login')
  } else {
    fetchUserProfile()
  }
})
</script>