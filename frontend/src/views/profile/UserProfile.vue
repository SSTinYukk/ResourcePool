<template>
  <div class="bg-white rounded-lg shadow-md overflow-hidden">
    <!-- 个人资料头部 -->
    <div class="bg-blue-600 text-white p-8 relative">
      <div class="flex flex-col md:flex-row items-center">
        <div class="relative mb-4 md:mb-0 md:mr-6">
          <Avatar 
            :image="user.avatar || undefined" 
            :label="!user.avatar ? user.username?.charAt(0).toUpperCase() : undefined"
            shape="circle" 
            size="xlarge" 
            class="w-24 h-24 border-4 border-white"
          />
          <Button 
            icon="pi pi-camera" 
            class="p-button-rounded p-button-sm absolute bottom-0 right-0 bg-white text-blue-600" 
            @click="openFileUpload"
          />
          <input type="file" ref="fileUpload" class="hidden" accept="image/*" @change="uploadAvatar" />
        </div>
        <div class="text-center md:text-left">
          <h1 class="text-2xl font-bold">{{ user.username }}</h1>
          <p class="text-blue-100">{{ user.email }}</p>
          <div class="mt-2 flex flex-wrap justify-center md:justify-start gap-2">
            <Tag value="积分: 350" severity="success" />
            <Tag value="资源: 12" severity="info" />
            <Tag value="论坛等级: 3" severity="warning" />
          </div>
        </div>
      </div>
    </div>

    <!-- 个人资料内容 -->
    <div class="p-6">
      <TabView>
        <!-- 基本信息 -->
        <TabPanel header="基本信息">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 class="text-lg font-semibold mb-4">个人信息</h3>
              <div class="p-fluid">
                <div class="field mb-4">
                  <label for="username" class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
                  <InputText id="username" v-model="user.username" class="w-full" />
                </div>
                <div class="field mb-4">
                  <label for="email" class="block text-sm font-medium text-gray-700 mb-1">电子邮箱</label>
                  <InputText id="email" v-model="user.email" class="w-full" disabled />
                </div>
                <div class="field mb-4">
                  <label for="bio" class="block text-sm font-medium text-gray-700 mb-1">个人简介</label>
                  <Textarea id="bio" v-model="user.bio" rows="4" class="w-full" />
                </div>
                <div class="field mb-4">
                  <label for="website" class="block text-sm font-medium text-gray-700 mb-1">个人网站</label>
                  <InputText id="website" v-model="user.website" class="w-full" />
                </div>
                <div class="field mb-4">
                  <label for="registerTime" class="block text-sm font-medium text-gray-700 mb-1">注册时间</label>
                  <InputText id="registerTime" v-model="user.registerTime" class="w-full" disabled />
                </div>
                <Button label="保存修改" icon="pi pi-check" class="mt-4" @click="saveProfile" />
              </div>
            </div>
            <div>
              <h3 class="text-lg font-semibold mb-4">账户安全</h3>
              <div class="p-fluid">
                <div class="field mb-4">
                  <label for="currentPassword" class="block text-sm font-medium text-gray-700 mb-1">当前密码</label>
                  <Password id="currentPassword" v-model="passwords.current" toggleMask class="w-full" />
                </div>
                <div class="field mb-4">
                  <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">新密码</label>
                  <Password id="newPassword" v-model="passwords.new" toggleMask class="w-full" />
                </div>
                <div class="field mb-4">
                  <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">确认新密码</label>
                  <Password id="confirmPassword" v-model="passwords.confirm" toggleMask class="w-full" />
                </div>
                <Button label="修改密码" icon="pi pi-lock" class="mt-4 p-button-secondary" @click="changePassword" />
              </div>
            </div>
          </div>
        </TabPanel>

        <!-- 我的资源 -->
        <TabPanel header="我的资源">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold">已上传资源</h3>
            <Button label="上传新资源" icon="pi pi-upload" @click="navigateTo('/resources/upload')" />
          </div>
          
          <DataTable :value="userResources" :paginator="true" :rows="5" responsiveLayout="scroll">
            <Column field="title" header="资源名称">
              <template #body="{data}">
                <div class="cursor-pointer text-blue-600 hover:text-blue-800" @click="navigateTo(`/resources/${data.id}`)">
                  {{ data.title }}
                </div>
              </template>
            </Column>
            <Column field="category" header="分类" />
            <Column field="downloads" header="下载次数" />
            <Column field="likes" header="点赞数" />
            <Column field="createTime" header="上传时间">
              <template #body="{data}">
                {{ formatDate(data.createTime) }}
              </template>
            </Column>
            <Column header="操作">
              <template #body="{data}">
                <div class="flex space-x-2">
                  <Button icon="pi pi-pencil" class="p-button-text p-button-sm" @click="navigateTo(`/resources/edit/${data.id}`)" />
                  <Button icon="pi pi-trash" class="p-button-text p-button-sm p-button-danger" @click="confirmDeleteResource(data.id)" />
                </div>
              </template>
            </Column>
          </DataTable>
          
          <div v-if="userResources.length === 0" class="text-center py-8">
            <i class="pi pi-file-export text-5xl text-gray-300 mb-4"></i>
            <p class="text-gray-500">你还没有上传任何资源</p>
            <Button label="上传资源" icon="pi pi-upload" class="mt-4" @click="navigateTo('/resources/upload')" />
          </div>
        </TabPanel>

        <!-- 我的积分 -->
        <TabPanel header="我的积分">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
            <div class="bg-blue-50 p-4 rounded-lg">
              <div class="text-blue-500 text-sm font-medium">当前积分</div>
              <div class="text-2xl font-bold">350</div>
            </div>
            <div class="bg-green-50 p-4 rounded-lg">
              <div class="text-green-500 text-sm font-medium">本月获得</div>
              <div class="text-2xl font-bold">+45</div>
            </div>
            <div class="bg-purple-50 p-4 rounded-lg">
              <div class="text-purple-500 text-sm font-medium">积分等级</div>
              <div class="text-2xl font-bold">3级</div>
            </div>
          </div>
          
          <h3 class="text-lg font-semibold mb-4">积分记录</h3>
          <DataTable :value="pointsHistory" :paginator="true" :rows="5" responsiveLayout="scroll">
            <Column field="description" header="描述" />
            <Column field="points" header="积分变动">
              <template #body="{data}">
                <span :class="data.points > 0 ? 'text-green-500' : 'text-red-500'">
                  {{ data.points > 0 ? '+' : '' }}{{ data.points }}
                </span>
              </template>
            </Column>
            <Column field="time" header="时间">
              <template #body="{data}">
                {{ formatDate(data.time) }}
              </template>
            </Column>
          </DataTable>
        </TabPanel>
      </TabView>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Password from 'primevue/password'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'

const router = useRouter()
const userStore = useUserStore()
const fileUpload = ref(null)

// 用户信息
const user = ref({})

// 获取用户资料
const fetchUserProfile = async () => {
  try {
    const response = await fetch('/api/user/profile', {
      headers: {
        'Authorization': `Bearer ${userStore.token}`
      }
    })
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.message || '获取用户资料失败')
    }
    const data = await response.json()
    if (!data.username || !data.email) {
      throw new Error('API返回数据格式不正确')
    }
    user.value = {
      ...data,
      bio: data.bio || '',
      website: data.website || '',
      registerTime: data.registerTime || ''
    }
  } catch (error) {
    console.error('获取用户资料错误:', error)
    // 显示错误提示给用户
    alert(error.message)
  }
}

// 初始化获取数据
onMounted(() => {
  fetchUserProfile()
  fetchUserResources()
  fetchPointsHistory()
})

// 密码修改
const passwords = ref({
  current: '',
  new: '',
  confirm: ''
})

// 用户资源数据
const userResources = ref([])

// 获取用户资源
const fetchUserResources = async () => {
  try {
    const response = await fetch('/api/user/resources', {
      headers: {
        'Authorization': `Bearer ${userStore.token}`
      }
    })
    if (!response.ok) throw new Error('获取用户资源失败')
    userResources.value = await response.json()
  } catch (error) {
    console.error(error)
  }
}

// 积分历史数据
const pointsHistory = ref([])

// 获取积分历史
const fetchPointsHistory = async () => {
  try {
    const response = await fetch('/api/user/points/history', {
      headers: {
        'Authorization': `Bearer ${userStore.token}`
      }
    })
    if (!response.ok) throw new Error('获取积分历史失败')
    pointsHistory.value = await response.json()
  } catch (error) {
    console.error(error)
  }
}

// 格式化日期
function formatDate(date) {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  }).format(date)
}

// 打开文件上传对话框
function openFileUpload() {
  fileUpload.value.click()
}

// 上传头像
function uploadAvatar(event) {
  const file = event.target.files[0]
  if (!file) return
  
  // 实际应用中应调用API上传文件
  // 这里仅做模拟
  const reader = new FileReader()
  reader.onload = (e) => {
    user.value.avatar = e.target.result
  }
  reader.readAsDataURL(file)
}

// 保存个人资料
function saveProfile() {
  // 实际应用中应调用API保存数据
  alert('个人资料已更新')
}

// 修改密码
function changePassword() {
  if (!passwords.value.current) {
    alert('请输入当前密码')
    return
  }
  
  if (passwords.value.new !== passwords.value.confirm) {
    alert('两次输入的新密码不一致')
    return
  }
  
  // 实际应用中应调用API修改密码
  alert('密码已修改')
  passwords.value = {
    current: '',
    new: '',
    confirm: ''
  }
}

// 确认删除资源
function confirmDeleteResource(resourceId) {
  // 实际应用中应使用确认对话框
  if (confirm('确定要删除这个资源吗？')) {
    // 调用API删除资源
    userResources.value = userResources.value.filter(r => r.id !== resourceId)
  }
}

// 页面导航
function navigateTo(path) {
  router.push(path)
}
</script>