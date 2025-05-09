<template>
  <div class="resource-search">
    <div class="search-container">
      <input 
        v-model="searchQuery" 
        type="text" 
        placeholder="搜索资源..."
        @keyup.enter="searchResources"
      />
      <button @click="searchResources">搜索</button>
    </div>
    
    <div class="filter-container bg-white p-4 rounded-lg shadow-sm mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">分类</label>
          <Dropdown
            v-model="selectedCategory"
            :options="categories"
            optionLabel="name"
            optionValue="id"
            placeholder="所有分类"
            class="w-full"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">标签</label>
          <Chips 
            v-model="selectedTags" 
            placeholder="输入标签" 
            class="w-full"
            @add="searchResources"
            @remove="searchResources"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">价格范围</label>
          <Dropdown
            v-model="priceRange"
            :options="[
              { label: '全部', value: 'all' },
              { label: '免费', value: 'free' },
              { label: '付费', value: 'paid' }
            ]"
            optionLabel="label"
            optionValue="value"
            class="w-full"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">排序</label>
          <Dropdown
            v-model="sortBy"
            :options="[
              { label: '最新', value: 'newest' },
              { label: '最热', value: 'popular' },
              { label: '下载量', value: 'downloads' },
              { label: '评分最高', value: 'rating' }
            ]"
            optionLabel="label"
            optionValue="value"
            class="w-full"
          />
        </div>
      </div>
    </div>
    
    <div class="resource-list">
      <div 
        v-for="resource in resources" 
        :key="resource.id" 
        class="resource-item"
        @click="$router.push({ name: 'resourceDetail', params: { id: resource.id } })"
      >
        <h3>{{ resource.title }}</h3>
        <p>{{ resource.description }}</p>
        <div class="meta">
          <span>{{ resource.user.username }}</span>
          <span>{{ resource.createdAt }}</span>
          <span>{{ resource.downloadCount }} 下载</span>
        </div>
      </div>
    </div>
    
    <div class="pagination" v-if="total > pageSize">
      <button 
        v-for="page in totalPages" 
        :key="page"
        :class="{ active: currentPage === page }"
        @click="goToPage(page)"
      >
        {{ page }}
      </button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../../api/axios'

export default {
  setup() {
    const router = useRouter()
    const searchQuery = ref('')
    const selectedCategory = ref(null)
    const selectedTags = ref([])
    const sortBy = ref('newest')
    const priceRange = ref('all')
    const resources = ref([])
    const categories = ref([])
    const currentPage = ref(1)
    const pageSize = ref(10)
    const total = ref(0)
    
    const fetchCategories = async () => {
      try {
        const response = await axios.get('/resources/categories')
        categories.value = response.data
      } catch (error) {
        console.error('获取分类失败:', error)
      }
    }
    
    const searchResources = async () => {
      try {
        const params = {
          q: searchQuery.value.trim(),
          category_id: selectedCategory.value,
          tags: selectedTags.value.join(','),
          sort: sortBy.value,
          price_range: priceRange.value,
          page: currentPage.value,
          pageSize: pageSize.value
        }
        
        const response = await axios.get('/resources/search', { params })
        resources.value = response.data.resources
        total.value = response.data.total
      } catch (error) {
        console.error('搜索资源失败:', error)
      }
    }
    
    const goToPage = (page) => {
      currentPage.value = page
      searchResources()
    }
    
    onMounted(() => {
      fetchCategories()
      searchResources()
    })
    
    return {
      searchQuery,
      selectedCategory,
      sortBy,
      resources,
      categories,
      currentPage,
      pageSize,
      total,
      searchResources,
      goToPage,
      totalPages: computed(() => Math.ceil(total.value / pageSize.value))
    }
  }
}
</script>

<style scoped>
.resource-search {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.search-container {
  display: flex;
  margin-bottom: 20px;
}

.search-container input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px 0 0 4px;
}

.search-container button {
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 0 4px 4px 0;
  cursor: pointer;
}

.filter-container {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.filter-container select {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.sort-options {
  display: flex;
  gap: 15px;
}

.resource-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.resource-item {
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.resource-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #4CAF50;
}

.meta {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
  font-size: 12px;
  color: #666;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
  gap: 5px;
}

.pagination button {
  padding: 5px 10px;
  border: 1px solid #ddd;
  background: white;
  cursor: pointer;
}

.pagination button.active {
  background: #4CAF50;
  color: white;
  border-color: #4CAF50;
}
</style>