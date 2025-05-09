<template>
  <div class="bg-gray-50 border-r border-gray-200 h-screen w-64 fixed left-0 top-0 z-10 transition-all duration-300 transform" :class="{ '-translate-x-full': !isOpen, 'translate-x-0': isOpen }">
    <div class="p-4 border-b border-gray-200">
      <div class="flex items-center justify-between">
        <router-link to="/" class="flex items-center">
          <span class="text-xl font-bold text-blue-600">V1-BK</span>
        </router-link>
        <Button icon="pi pi-times" class="p-button-text p-button-rounded" @click="toggleSidebar" />
      </div>
    </div>
    
    <div class="p-4">
      <div class="mb-6">
        <div class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2">主菜单</div>
        <ul class="space-y-2">
          <li>
            <router-link to="/" class="flex items-center px-2 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-md">
              <i class="pi pi-home mr-3"></i>
              <span>首页</span>
            </router-link>
          </li>
          <li>
            <router-link to="/resources" class="flex items-center px-2 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-md">
              <i class="pi pi-file mr-3"></i>
              <span>资源中心</span>
            </router-link>
          </li>
          <li>
            <router-link to="/forum" class="flex items-center px-2 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-md">
              <i class="pi pi-comments mr-3"></i>
              <span>论坛</span>
            </router-link>
          </li>
          <li>
            <router-link to="/chat" class="flex items-center px-2 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-md">
              <i class="pi pi-comment mr-3"></i>
              <span>聊天</span>
            </router-link>
          </li>
        </ul>
      </div>
      
      <div class="mb-6" v-if="userStore.isLoggedIn">
        <div class="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-2">个人中心</div>
        <ul class="space-y-2">
          <li>
            <router-link to="/profile" class="flex items-center px-2 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-md">
              <i class="pi pi-user mr-3"></i>
              <span>个人资料</span>
            </router-link>
          </li>
          <li>
            <router-link to="/my-resources" class="flex items-center px-2 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-md">
              <i class="pi pi-folder mr-3"></i>
              <span>我的资源</span>
            </router-link>
          </li>
          <li>
            <router-link to="/my-points" class="flex items-center px-2 py-2 text-gray-700 hover:bg-blue-50 hover:text-blue-600 rounded-md">
              <i class="pi pi-star mr-3"></i>
              <span>我的积分</span>
            </router-link>
          </li>
        </ul>
      </div>
      
      <div v-if="userStore.isLoggedIn">
        <Button label="退出登录" icon="pi pi-sign-out" class="p-button-text p-button-danger w-full" @click="logout" />
      </div>
      <div v-else class="space-y-2">
        <router-link to="/login" class="block text-center py-2 text-gray-700 hover:text-blue-600 transition duration-300">
          登录
        </router-link>
        <router-link to="/register" class="block bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition duration-300 text-center">
          注册
        </router-link>
      </div>
    </div>
  </div>
  
  <!-- 遮罩层 -->
  <div 
    v-if="isOpen" 
    class="fixed inset-0 bg-black bg-opacity-50 z-0 md:hidden"
    @click="toggleSidebar"
  ></div>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Button from 'primevue/button'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:isOpen'])
const router = useRouter()
const userStore = useUserStore()

function toggleSidebar() {
  emit('update:isOpen', !props.isOpen)
}

async function logout() {
  await userStore.logout()
  router.push('/login')
  toggleSidebar()
}
</script>