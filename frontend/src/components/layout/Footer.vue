<template>
  <footer class="bg-gray-100 border-t border-gray-200 mt-auto">
    <div class="container mx-auto px-4 py-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <div>
          <h3 class="text-lg font-semibold text-gray-800 mb-3">关于我们</h3>
          <p class="text-gray-600 text-sm">微机原理资源池致力于为用户提供高质量的学习资源和交流环境，促进知识的分享与创新。</p>
        </div>
        
        <div>
          <h3 class="text-lg font-semibold text-gray-800 mb-3">快速链接</h3>
          <ul class="space-y-2">
            <li>
              <router-link to="/" class="text-gray-600 hover:text-blue-600 text-sm">首页</router-link>
            </li>
            <li>
              <router-link to="/resources" class="text-gray-600 hover:text-blue-600 text-sm">资源中心</router-link>
            </li>
            <li>
              <router-link to="/forum" class="text-gray-600 hover:text-blue-600 text-sm">论坛</router-link>
            </li>
            <li>
              <router-link to="/chat" class="text-gray-600 hover:text-blue-600 text-sm">聊天</router-link>
            </li>
          </ul>
        </div>
        
        <div>
          <h3 class="text-lg font-semibold text-gray-800 mb-3">联系我们</h3>
          <ul class="space-y-2">
            <li class="flex items-center text-gray-600 text-sm">
              <i class="pi pi-envelope mr-2"></i>
              <span>contact@v1bk.com</span>
            </li>
            <li class="flex items-center text-gray-600 text-sm">
              <i class="pi pi-phone mr-2"></i>
              <span>+86 123 4567 8901</span>
            </li>
            <li class="flex items-center text-gray-600 text-sm">
              <i class="pi pi-map-marker mr-2"></i>
              <span>北京市海淀区</span>
            </li>
          </ul>
        </div>
        
        <div>
          <h3 class="text-lg font-semibold text-gray-800 mb-3">关注我们</h3>
          <div class="flex space-x-4">
            <a href="#" class="text-gray-600 hover:text-blue-600">
              <i class="pi pi-facebook text-xl"></i>
            </a>
            <a href="#" class="text-gray-600 hover:text-blue-600">
              <i class="pi pi-twitter text-xl"></i>
            </a>
            <a href="#" class="text-gray-600 hover:text-blue-600">
              <i class="pi pi-instagram text-xl"></i>
            </a>
            <a href="#" class="text-gray-600 hover:text-blue-600">
              <i class="pi pi-github text-xl"></i>
            </a>
          </div>
        </div>
      </div>
      
      <div class="border-t border-gray-200 mt-6 pt-6 text-center">
        <p class="text-gray-600 text-sm">&copy; {{ new Date().getFullYear() }} 微机原理资源池. 保留所有权利.</p>
      </div>
    </div>
  </footer>
</template>

<script setup>
// 无需额外逻辑
</script>