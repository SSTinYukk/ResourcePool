<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold">论坛分类</h2>
    </div>
    
    <div v-if="loading" class="text-center py-8">
      <p>加载中...</p>
    </div>
    
    <div v-else-if="categories.length === 0" class="bg-gray-50 rounded-lg p-8 text-center">
      <p class="text-gray-500">暂无分类</p>
    </div>
    
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div 
        v-for="category in categories" 
        :key="category.id" 
        class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow"
      >
        <div class="p-5">
          <h3 class="text-lg font-semibold mb-2">
            <router-link 
              :to="`/forum/categories/${category.id}`" 
              class="hover:text-blue-600"
            >
              {{ category.name }}
            </router-link>
          </h3>
          <p class="text-gray-600 text-sm mb-4 line-clamp-2">{{ category.description }}</p>
          <div class="flex justify-between items-center text-sm text-gray-500">
            <span>{{ category.topic_count }} 个话题</span>
            <span>{{ formatDate(category.updated_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { forumApi } from '../../api';

const categories = ref([]);
const loading = ref(true);

onMounted(async () => {
  try {
    const response = await forumApi.getCategories();
    categories.value = response.data;
  } catch (error) {
    console.error('获取分类失败:', error);
  } finally {
    loading.value = false;
  }
});

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric' });
};
</script>