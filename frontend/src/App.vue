<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- 导航栏 -->
    <Navbar />
    
    <div class="flex flex-1">
      <!-- 侧边栏 -->
      <Sidebar v-model:isOpen="sidebarOpen" />
      
      <!-- 主要内容区域 -->
      <main class="flex-1 p-4 transition-all duration-300" :class="{ 'md:ml-64': sidebarOpen }">
        <div class="container mx-auto">
          <!-- 页面内容 -->
          <Toast />
          <ConfirmDialog />
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </main>
    </div>
    
    <!-- 页脚 -->
    <Footer />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from './stores/user'
import Navbar from './components/layout/Navbar.vue'
import Sidebar from './components/layout/Sidebar.vue'
import Footer from './components/layout/Footer.vue'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'

const route = useRoute()
const userStore = useUserStore()
const sidebarOpen = ref(false)
</script>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>