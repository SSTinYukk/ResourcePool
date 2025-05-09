<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex flex-col items-center justify-center min-h-[60vh]">
      <h1 class="text-6xl font-bold text-gray-800 mb-4">404</h1>
      <p class="text-xl text-gray-600 mb-8">页面未找到</p>
      <Button 
        label="返回首页" 
        icon="pi pi-home" 
        class="p-button-primary"
        @click="goHome"
      />
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import Button from 'primevue/button'

const router = useRouter()

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
/* 保持与项目一致的样式 */
</style>