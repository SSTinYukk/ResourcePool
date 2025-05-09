<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">资源中心</h1>
      <router-link 
        v-if="userStore.isLoggedIn" 
        to="/resources/upload" 
        class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition-colors"
      >
        上传资源
      </router-link>
    </div>
    
    <!-- 搜索和筛选 -->
    <div class="bg-white p-4 rounded-lg shadow-sm mb-6">
      <div class="flex flex-col md:flex-row gap-4">
        <div class="flex-1">
          <div class="relative">
            <input 
              type="text" 
              v-model="searchQuery" 
              @keyup.enter="searchResources"
              placeholder="搜索资源..." 
              class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
            <div class="absolute left-3 top-2.5 text-gray-400">
              <i class="pi pi-search"></i>
            </div>
          </div>
        </div>
        
        <div class="flex gap-4">
          <MultiSelect 
            v-model="selectedCategories" 
            :options="categories" 
            optionLabel="name" 
            optionValue="id"
            placeholder="选择分类"
            display="chip"
            class="w-full md:w-80"
            @change="loadResources"
          />
          
          <Dropdown 
            v-model="sortBy" 
            :options="sortOptions" 
            optionLabel="label" 
            optionValue="value"
            placeholder="排序方式"
            class="w-full md:w-80"
            @change="loadResources"
          />
          
          <Dropdown 
            v-model="priceFilter" 
            :options="priceOptions" 
            optionLabel="label" 
            optionValue="value"
            placeholder="价格筛选"
            class="w-full md:w-80"
            @change="loadResources"
          />
        </div>
      </div>
    </div>
    
    <!-- 资源列表 -->
    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500 mb-2"></div>
      <p class="text-gray-600">正在加载资源，请稍候...</p>
    </div>
    
    <div v-else-if="resources.length === 0" class="bg-gray-50 rounded-lg p-12 text-center">
      <p class="text-gray-500 mb-4">暂无资源</p>
      <router-link 
        v-if="userStore.isLoggedIn" 
        to="/resources/upload" 
        class="text-blue-600 hover:underline"
      >
        上传第一个资源
      </router-link>
    </div>
    
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div 
        v-for="resource in resources" 
        :key="resource.id" 
        class="bg-white rounded-lg shadow-md overflow-hidden hover-card-effect"
      >
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
          
          <div class="mt-4 pt-4 border-t border-gray-100 flex justify-between items-center">
            <span class="text-sm">
              <span v-if="resource.points_required > 0" class="text-orange-600">
                {{ resource.points_required }} 积分
              </span>
              <span v-else class="text-green-600">免费</span>
            </span>
            
            <router-link 
              :to="`/resources/${resource.id}`" 
              class="text-blue-600 hover:underline text-sm"
            >
              查看详情
            </router-link>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 分页 -->
    <div v-if="resources.length > 0" class="mt-8 flex justify-center">
      <div class="flex space-x-1">
        <button 
          @click="changePage(currentPage - 1)" 
          :disabled="currentPage === 1"
          class="px-4 py-2 border rounded-md" 
          :class="currentPage === 1 ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-white hover:bg-gray-50'"
        >
          上一页
        </button>
        
        <button 
          v-for="page in totalPages" 
          :key="page" 
          @click="changePage(page)"
          class="px-4 py-2 border rounded-md" 
          :class="currentPage === page ? 'bg-blue-600 text-white' : 'bg-white hover:bg-gray-50'"
        >
          {{ page }}
        </button>
        
        <button 
          @click="changePage(currentPage + 1)" 
          :disabled="currentPage === totalPages"
          class="px-4 py-2 border rounded-md" 
          :class="currentPage === totalPages ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-white hover:bg-gray-50'"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { resourceApi } from '../../api';
import { useUserStore } from '../../stores/user';

const userStore = useUserStore();

const resources = ref([]);
const categories = ref([]);
const loading = ref(true);

const searchQuery = ref('');
const selectedCategories = ref([]);
const sortBy = ref('created_at:desc');
const priceFilter = ref('all');

const sortOptions = ref([
  { label: '最新上传', value: 'created_at:desc' },
  { label: '下载最多', value: 'download_count:desc' },
  { label: '标题 A-Z', value: 'title:asc' },
  { label: '评分最高', value: 'rating:desc' }
]);

const priceOptions = ref([
  { label: '全部', value: 'all' },
  { label: '免费', value: 'free' },
  { label: '付费', value: 'paid' }
]);

const currentPage = ref(1);
const totalPages = ref(1);
const pageSize = 12;

onMounted(async () => {
  await Promise.all([loadResources(), loadCategories()]);
});

const loadResources = async () => {
  loading.value = true;
  
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize,
      sort: sortBy.value
    };
    
    if (selectedCategories.value.length > 0) {
      params.category_ids = selectedCategories.value.join(',');
    }
    
    if (priceFilter.value === 'free') {
      params.points_required = 0;
    } else if (priceFilter.value === 'paid') {
      params.points_required = '>0';
    }
    
    if (searchQuery.value) {
      params.query = searchQuery.value;
    }
    
    console.log('请求参数:', params);
    const response = await resourceApi.getResources(params);
    
    if (!response || !response.data) {
      throw new Error(`返回数据格式错误: ${JSON.stringify(response)}`);
    }
    
    console.log('API响应:', response);
    resources.value = response.data.resources;
    totalPages.value = Math.ceil(response.data.total / pageSize);
    
    // 如果数据为空，重置为第一页
    if (resources.value.length === 0 && currentPage.value > 1) {
      currentPage.value = 1;
      await loadResources();
      return;
    }
  } catch (error) {
    console.error('获取资源列表失败:', error);
    resources.value = [];
    totalPages.value = 1;
    
    // 显示错误提示
    if (error.response) {
      console.error('API错误详情:', {
        status: error.response.status,
        data: error.response.data
      });
    }
  } finally {
    loading.value = false;
  }
};

const loadCategories = async () => {
  try {
    const response = await resourceApi.getCategories();
    categories.value = response.data;
  } catch (error) {
    console.error('获取分类失败:', error);
  }
};

const searchResources = () => {
  currentPage.value = 1; // 重置到第一页
  loadResources();
};

const changePage = (page) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  loadResources();
};

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric' });
};
</script>