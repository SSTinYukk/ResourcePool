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
          <Button 
            label="更换头像" 
            icon="pi pi-image" 
            class="p-button-text p-button-sm mt-2"
            @click="showAvatarDialog = true"
          />
        </div>
      </div>
      
      <Dialog v-model:visible="showAvatarDialog" header="上传头像" :modal="true" class="w-full max-w-2xl">
        <div class="flex flex-col gap-4">
          <FileUpload 
            mode="basic" 
            name="avatar" 
            :auto="false"
            :customUpload="true"
            @select="onFileSelect"
            accept="image/*"
            :maxFileSize="1000000"
            chooseLabel="选择图片"
          />
          
          <Dialog v-model:visible="showCropDialog" header="裁剪头像" :modal="true" class="w-full max-w-2xl">
            <div class="flex flex-col items-center gap-4">
              <Image :src="previewImage" alt="预览" class="max-h-96" v-if="previewImage" />
              <div class="flex gap-2">
                <Button label="取消" class="p-button-text" @click="showCropDialog = false" />
                <Button label="确认裁剪" class="p-button-primary" @click="handleCropConfirm" />
              </div>
            </div>
          </Dialog>
        </div>
      </Dialog>

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
import axios from 'axios'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Avatar from 'primevue/avatar'
import Dialog from 'primevue/dialog'
import Password from 'primevue/password'
import FileUpload from 'primevue/fileupload'
import Image from 'primevue/image'
import { useToast } from 'primevue/usetoast'

const router = useRouter()
const toast = useToast()
const userStore = useUserStore()
const showAvatarDialog = ref(false)
const showCropDialog = ref(false)
const previewImage = ref('')
const selectedFile = ref(null)
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

const onFileSelect = (event) => {
  const file = event.files[0]
  if (!file) return
  
  // 验证文件类型和大小
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    toast.add({ 
      severity: 'error', 
      summary: '错误', 
      detail: '仅支持JPEG、PNG或GIF格式的图片', 
      life: 3000 
    })
    return
  }
  
  if (file.size > 5 * 1024 * 1024) { // 5MB限制
    toast.add({ 
      severity: 'error', 
      summary: '错误', 
      detail: '头像文件大小不能超过5MB', 
      life: 3000 
    })
    return
  }
  
  selectedFile.value = file
  previewImage.value = URL.createObjectURL(file)
  showCropDialog.value = true
}

const handleCropConfirm = async () => {
  if (!selectedFile.value) return
  
  const loading = ref(true)
  try {
    // 验证文件是否存在
    if (!selectedFile.value) {
      throw new Error('未选择有效的图片文件')
    }
    // 创建canvas进行图片裁剪
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    const img = document.createElement('img')
    
    // 等待图片加载
    await new Promise((resolve) => {
      img.onload = resolve
      img.src = previewImage.value
    })
    
    // 设置裁剪区域(这里使用中心正方形裁剪作为示例)
    const size = Math.min(img.width, img.height)
    const x = (img.width - size) / 2
    const y = (img.height - size) / 2
    
    // 设置canvas尺寸并绘制裁剪后的图片
    canvas.width = 200 // 输出宽度
    canvas.height = 200 // 输出高度
    ctx.drawImage(img, x, y, size, size, 0, 0, canvas.width, canvas.height)
    
    // 将canvas转为Blob对象
    const blob = await new Promise((resolve) => {
      canvas.toBlob(resolve, 'image/jpeg', 0.9)
    })
    
    // 创建FormData并添加裁剪后的图片
    const formData = new FormData()
    formData.append('avatar', blob, 'avatar.jpg')
    
    const token = typeof localStorage !== 'undefined' ? localStorage.getItem('token') : null
    if (!token) {
      throw new Error('未找到用户token，请重新登录')
    }
    
    const response = await axios.post('/api/user/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (response.data && response.data.avatar) {
      // 更新用户头像URL并刷新显示
      user.value.avatar = response.data.avatar
      toast.add({ 
        severity: 'success', 
        summary: '成功', 
        detail: '头像上传成功', 
        life: 3000 
      })
      
      // 刷新用户数据
      await fetchUserProfile()
    }
  } catch (error) {
    console.error('头像上传失败:', error)
    toast.add({ 
      severity: 'error', 
      summary: '错误', 
      detail: '头像上传失败: ' + (error.response?.data?.error || error.message),
      life: 3000 
    })
    // 确保在错误时也重置状态
    showCropDialog.value = false
    previewImage.value = ''
    selectedFile.value = null
  } finally {
    loading.value = false
    showAvatarDialog.value = false
    showCropDialog.value = false
    previewImage.value = ''
    selectedFile.value = null
  }
}

const onAvatarUploadSuccess = (event) => {
  if (event.xhr.response) {
    const response = JSON.parse(event.xhr.response)
    if (response && response.avatar_url) {
      user.value.avatar = response.avatar_url
      toast.add({ severity: 'success', summary: '成功', detail: '头像上传成功', life: 3000 })
    }
  }
}

const onAvatarUploadError = (event) => {
  let errorMessage = '头像上传失败'
  try {
    const response = JSON.parse(event.xhr.response)
    errorMessage += ': ' + (response.message || event.xhr.statusText)
  } catch {
    errorMessage += ': ' + event.xhr.statusText
  }
  toast.add({
    severity: 'error',
    summary: '错误',
    detail: errorMessage,
    life: 3000
  })
}

const fetchUserProfile = async () => {
  try {
    const token = typeof localStorage !== 'undefined' ? localStorage.getItem('token') : null
    if (!token) {
      router.push('/login')
      return
    }
    const response = await axios.get('/api/user/profile', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    if (response.data) {
      user.value = {
        name: response.data.username,
        email: response.data.email,
        username: response.data.username,
        avatar: response.data.avatar || '',
        createdAt: new Date(response.data.created_at).toLocaleString()
      }
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
    toast.add({
      severity: 'error',
      summary: '错误',
      detail: '获取用户信息失败: ' + (error.response?.data?.message || error.message),
      life: 3000
    })
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