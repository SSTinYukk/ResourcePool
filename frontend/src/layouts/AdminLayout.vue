<template>
  <div class="min-h-screen bg-gray-100 flex flex-col">
    <!-- 管理员导航栏 -->
    <div class="bg-white shadow-sm z-10">
      <div class="container mx-auto px-4 py-2 flex justify-between items-center">
        <div class="flex items-center">
          <router-link to="/" class="flex items-center">
            <span class="text-xl font-bold text-primary">资源池</span>
          </router-link>
          <span class="ml-4 text-gray-500">管理后台</span>
        </div>
        
        <div class="flex items-center gap-4">
          <span class="text-sm text-gray-600">欢迎，{{ userStore.user?.username || '管理员' }}</span>
          <Button icon="pi pi-sign-out" text severity="secondary" aria-label="退出" @click="logout" v-tooltip.bottom="'退出登录'" />
        </div>
      </div>
    </div>
    
    <div class="flex flex-1 overflow-hidden">
      <!-- 侧边导航 -->
      <div class="w-64 bg-white shadow-sm flex-shrink-0 h-[calc(100vh-4rem)] overflow-y-auto">
        <div class="p-4">
          <ul class="space-y-1">
            <li>
              <router-link to="/admin" class="flex items-center p-2 rounded-lg hover:bg-gray-100" :class="{ 'bg-primary/10 text-primary': isActive('/admin') }">
                <i class="pi pi-home mr-2"></i>
                <span>仪表盘</span>
              </router-link>
            </li>
            <li>
              <router-link to="/admin/users" class="flex items-center p-2 rounded-lg hover:bg-gray-100" :class="{ 'bg-primary/10 text-primary': isActive('/admin/users') }">
                <i class="pi pi-users mr-2"></i>
                <span>用户管理</span>
              </router-link>
            </li>
            <li>
              <router-link to="/admin/resources" class="flex items-center p-2 rounded-lg hover:bg-gray-100" :class="{ 'bg-primary/10 text-primary': isActive('/admin/resources') }">
                <i class="pi pi-database mr-2"></i>
                <span>资源管理</span>
              </router-link>
            </li>
            <li>
              <router-link to="/admin/forum" class="flex items-center p-2 rounded-lg hover:bg-gray-100" :class="{ 'bg-primary/10 text-primary': isActive('/admin/forum') }">
                <i class="pi pi-comments mr-2"></i>
                <span>论坛管理</span>
              </router-link>
            </li>
            <li>
              <router-link to="/admin/points" class="flex items-center p-2 rounded-lg hover:bg-gray-100" :class="{ 'bg-primary/10 text-primary': isActive('/admin/points') }">
                <i class="pi pi-wallet mr-2"></i>
                <span>积分管理</span>
              </router-link>
            </li>
            <li>
              <router-link to="/admin/stats" class="flex items-center p-2 rounded-lg hover:bg-gray-100" :class="{ 'bg-primary/10 text-primary': isActive('/admin/stats') }">
                <i class="pi pi-chart-bar mr-2"></i>
                <span>数据统计</span>
              </router-link>
            </li>
          </ul>
        </div>
      </div>
      
      <!-- 主要内容区域 -->
      <div class="flex-1 p-6 overflow-y-auto">
        <Toast />
        <ConfirmDialog />
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const confirm = useConfirm()
const toast = useToast()

// 检查当前路由是否激活
const isActive = (path) => {
  return route.path === path || route.path.startsWith(`${path}/`)
}

// 退出登录
const logout = () => {
  confirm.require({
    message: '确定要退出登录吗？',
    header: '退出确认',
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: '确定',
    rejectLabel: '取消',
    accept: () => {
      userStore.logout()
      toast.add({ severity: 'success', summary: '成功', detail: '已退出登录', life: 3000 })
      router.push('/')
    }
  })
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>