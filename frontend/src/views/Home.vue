<template>
  <div>
    <!-- 欢迎横幅 -->
    <div class="bg-gradient-to-r from-blue-600 to-indigo-700 text-white py-12 px-4 rounded-lg shadow-lg mb-8">
      <div class="max-w-4xl mx-auto">
        <h1 class="text-3xl md:text-4xl font-bold mb-4">欢迎来到 V1-BK 知识共享平台</h1>
        <p class="text-xl mb-6">分享知识，共同成长。探索高质量资源，参与社区讨论，提升专业技能。</p>
        <div class="flex flex-wrap gap-4">
          <router-link to="/resources" class="bg-white text-blue-600 hover:bg-blue-50 px-6 py-2 rounded-md font-medium transition-colors">
            浏览资源
          </router-link>
          <router-link to="/forum" class="bg-transparent border border-white text-white hover:bg-white hover:text-blue-600 px-6 py-2 rounded-md font-medium transition-colors">
            参与讨论
          </router-link>
        </div>
      </div>
    </div>
    
    <!-- 平台特色 -->
    <div class="mb-12">
      <h2 class="text-2xl font-bold mb-6 text-center">平台特色</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="bg-white p-6 rounded-lg shadow-md">
          <div class="text-blue-600 text-4xl mb-4">
            <i class="pi pi-file"></i>
          </div>
          <h3 class="text-xl font-semibold mb-2">优质资源</h3>
          <p class="text-gray-600">提供丰富的学习资源，包括文档、教程、代码示例等，助力您的学习和工作。</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow-md">
          <div class="text-blue-600 text-4xl mb-4">
            <i class="pi pi-comments"></i>
          </div>
          <h3 class="text-xl font-semibold mb-2">社区交流</h3>
          <p class="text-gray-600">活跃的论坛社区，与志同道合的伙伴交流经验，解决问题，共同进步。</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow-md">
          <div class="text-blue-600 text-4xl mb-4">
            <i class="pi pi-comment"></i>
          </div>
          <h3 class="text-xl font-semibold mb-2">AI助手</h3>
          <p class="text-gray-600">智能AI助手随时为您解答问题，提供专业建议，让学习更高效。</p>
        </div>
      </div>
    </div>
    
    <!-- 最新资源 -->
    <div class="mb-12">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold">最新资源</h2>
        <router-link to="/resources" class="text-blue-600 hover:underline">查看全部</router-link>
      </div>
      
      <div v-if="loading" class="text-center py-8">
        <p>加载中...</p>
      </div>
      
      <div v-else-if="latestResources.length === 0" class="bg-gray-50 rounded-lg p-8 text-center">
        <p class="text-gray-500">暂无资源</p>
      </div>
      
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="resource in latestResources" :key="resource.id" class="bg-white rounded-lg shadow-md overflow-hidden">
          <div class="p-5">
            <div class="flex justify-between items-start mb-2">
              <h3 class="text-lg font-semibold truncate">
                <router-link :to="`/resources/${resource.id}`" class="hover:text-blue-600">
                  {{ resource.title }}
                </router-link>
              </h3>
              <span class="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">
                {{ resource.category.name }}
              </span>
            </div>
            <p class="text-gray-600 text-sm mb-4 line-clamp-2">{{ resource.description }}</p>
            <div class="flex justify-between items-center text-sm text-gray-500">
              <div class="flex items-center">
                <span>{{ resource.user.username }}</span>
              </div>
              <div class="flex items-center space-x-4">
                <span class="flex items-center">
                  <i class="pi pi-download mr-1"></i> {{ resource.download_count }}
                </span>
                <span>{{ formatDate(resource.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 热门话题 -->
    <div class="mb-12">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold">热门话题</h2>
        <router-link to="/forum" class="text-blue-600 hover:underline">查看全部</router-link>
      </div>
      
      <div v-if="loadingTopics" class="text-center py-8">
        <p>加载中...</p>
      </div>
      
      <div v-else-if="hotTopics.length === 0" class="bg-gray-50 rounded-lg p-8 text-center">
        <p class="text-gray-500">暂无话题</p>
      </div>
      
      <div v-else class="bg-white rounded-lg shadow-md overflow-hidden">
        <div class="divide-y">
          <div v-for="topic in hotTopics" :key="topic.id" class="p-4 hover:bg-gray-50">
            <div class="flex justify-between">
              <h3 class="font-medium">
                <router-link :to="`/forum/topics/${topic.id}`" class="hover:text-blue-600">
                  {{ topic.title }}
                </router-link>
              </h3>
              <span class="text-sm text-gray-500">{{ formatDate(topic.created_at) }}</span>
            </div>
            <div class="flex items-center mt-2 text-sm text-gray-500">
              <span class="mr-4">{{ topic.user.username }}</span>
              <span class="flex items-center mr-4">
                <i class="pi pi-eye mr-1"></i> {{ topic.view_count }}
              </span>
              <span class="flex items-center">
                <i class="pi pi-comments mr-1"></i> {{ topic.reply_count }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { resourceApi, forumApi } from '../api';

const latestResources = ref([]);
const hotTopics = ref([]);
const loading = ref(true);
const loadingTopics = ref(true);

onMounted(async () => {
  try {
    // 获取最新资源
    const resourceResponse = await resourceApi.getResources({ limit: 6, sort: 'created_at:desc' });
    latestResources.value = resourceResponse.data;
  } catch (error) {
    console.error('获取最新资源失败:', error);
  } finally {
    loading.value = false;
  }
  
  try {
    // 获取热门话题
    const topicResponse = await forumApi.getTopics({ limit: 5, sort: 'view_count:desc' });
    hotTopics.value = topicResponse.data;
  } catch (error) {
    console.error('获取热门话题失败:', error);
  } finally {
    loadingTopics.value = false;
  }
});

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric' });
};
</script>