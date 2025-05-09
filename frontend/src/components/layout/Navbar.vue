<template>
  <nav class="bg-white shadow-md">
    <div class="container mx-auto px-4">
      <div class="flex justify-between items-center py-3">
        <!-- Logo和网站名称 -->
        <div class="flex items-center space-x-2">
          <router-link to="/" class="flex items-center">
            <span class="text-xl font-bold text-blue-600">V1-BK</span>
            <span class="ml-2 text-gray-600 text-sm">知识共享平台</span>
          </router-link>
        </div>

        <!-- 导航链接 -->
        <div class="hidden md:flex space-x-6">
          <router-link to="/" class="text-gray-700 hover:text-blue-600 transition duration-300">
            首页
          </router-link>
          <router-link to="/resources" class="text-gray-700 hover:text-blue-600 transition duration-300">
            资源中心
          </router-link>
          <router-link to="/forum" class="text-gray-700 hover:text-blue-600 transition duration-300">
            论坛
          </router-link>
          <router-link to="/chat" class="text-gray-700 hover:text-blue-600 transition duration-300">
            聊天
          </router-link>
        </div>

        <!-- 用户菜单 -->
        <div class="flex items-center space-x-4">
          <template v-if="userStore.isLoggedIn">
            <Button icon="pi pi-bell" class="p-button-text p-button-rounded" badge="3" badgeClass="p-badge-danger" />
            <Menu ref="menu" :model="userMenuItems" :popup="true" />
            <Avatar 
              :image="userStore.user?.avatar || undefined" 
              :label="!userStore.user?.avatar ? userStore.user?.username?.charAt(0).toUpperCase() : undefined"
              shape="circle" 
              class="cursor-pointer" 
              @click="toggleUserMenu"
            />
          </template>
          <template v-else>
            <router-link to="/login" class="text-gray-700 hover:text-blue-600 transition duration-300">
              登录
            </router-link>
            <router-link to="/register" class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition duration-300">
              注册
            </router-link>
          </template>
        </div>

        <!-- 移动端菜单按钮 -->
        <div class="md:hidden">
          <Button icon="pi pi-bars" class="p-button-text p-button-rounded" @click="toggleMobileMenu" />
        </div>
      </div>
    </div>

    <!-- 移动端菜单 -->
    <Sidebar v-model:visible="mobileMenuVisible" position="right" class="p-sidebar-md">
      <div class="flex flex-col space-y-4 p-4">
        <router-link to="/" class="text-gray-700 hover:text-blue-600 transition duration-300">
          首页
        </router-link>
        <router-link to="/resources" class="text-gray-700 hover:text-blue-600 transition duration-300">
          资源中心
        </router-link>
        <router-link to="/forum" class="text-gray-700 hover:text-blue-600 transition duration-300">
          论坛
        </router-link>
        <router-link to="/resources/upload" class="text-gray-700 hover:text-blue-600 transition duration-300" v-if="userStore.isLoggedIn">
          上传资源
        </router-link>
        <router-link to="/chat" class="text-gray-700 hover:text-blue-600 transition duration-300">
          聊天
        </router-link>
        <Divider />
        <template v-if="userStore.isLoggedIn">
          <router-link to="/profile" class="text-gray-700 hover:text-blue-600 transition duration-300">
            个人资料
          </router-link>
          <Button label="退出登录" class="p-button-text p-button-danger" @click="logout" />
        </template>
        <template v-else>
          <router-link to="/login" class="text-gray-700 hover:text-blue-600 transition duration-300">
            登录
          </router-link>
          <router-link to="/register" class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition duration-300 text-center">
            注册
          </router-link>
        </template>
      </div>
    </Sidebar>
  </nav>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import Sidebar from 'primevue/sidebar'
import Divider from 'primevue/divider'

const router = useRouter()
const userStore = useUserStore()
const menu = ref(null)
const mobileMenuVisible = ref(false)

const userMenuItems = [
  {
    label: '个人资料',
    icon: 'pi pi-user',
    command: () => router.push('/profile')
  },
  {
    label: '我的资源',
    icon: 'pi pi-file',
    command: () => router.push('/my-resources')
  },
  {
    label: '我的积分',
    icon: 'pi pi-star',
    command: () => router.push('/my-points')
  },
  {
    separator: true
  },
  {
    label: '退出登录',
    icon: 'pi pi-sign-out',
    command: logout
  }
]

function toggleUserMenu(event) {
  menu.value.toggle(event)
}

function toggleMobileMenu() {
  mobileMenuVisible.value = !mobileMenuVisible.value
}

async function logout() {
  await userStore.logout()
  router.push('/login')
}
</script>