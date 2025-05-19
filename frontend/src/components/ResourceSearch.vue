<template>
  <div class="resource-search">
    <form @submit.prevent="handleSearch">
      <div class="search-form">
        <input 
          v-model="searchQuery.q" 
          type="text" 
          placeholder="输入关键词搜索"
          class="search-input"
        />
        
        <select v-model="searchQuery.category_id" class="category-select">
          <option value="">所有分类</option>
          <option v-for="category in categories" :key="category.id" :value="category.id">
            {{ category.name }}
          </option>
        </select>
        
        <select v-model="searchQuery.sort" class="sort-select">
          <option value="created_at:desc">最新</option>
          <option value="download_count:desc">最热</option>
          <option value="title:asc">名称(A-Z)</option>
          <option value="rating:desc">评分最高</option>
        </select>
        
        <button type="submit" class="search-button" :disabled="loading">
          {{ loading ? '搜索中...' : '搜索' }}
        </button>
      </div>
    </form>
    
    <div v-if="loading" class="loading-indicator">加载中...</div>
    
    <div v-if="error" class="error-message">{{ error }}</div>
    
    <div v-if="resources.length > 0" class="search-results">
      <div class="resource-list">
        <div v-for="resource in resources" :key="resource.id" class="resource-item">
          <h3>{{ resource.title }}</h3>
          <p>{{ resource.description }}</p>
          <div class="resource-meta">
            <span>分类: {{ resource.category.name }}</span>
            <span>上传者: {{ resource.user.username }}</span>
          </div>
        </div>
      </div>
      
      <div class="pagination">
        <button 
          @click="changePage(currentPage - 1)" 
          :disabled="currentPage === 1"
        >
          上一页
        </button>
        
        <span>第 {{ currentPage }} 页 / 共 {{ totalPages }} 页</span>
        
        <button 
          @click="changePage(currentPage + 1)" 
          :disabled="currentPage === totalPages"
        >
          下一页
        </button>
      </div>
    </div>
    
    <div v-else-if="!loading && !error" class="no-results">
      没有找到匹配的资源
    </div>
  </div>
</template>

<script>
export default {
  name: 'ResourceSearch',
  data() {
    return {
      searchQuery: {
        q: '',
        category_id: '',
        tags: '',
        sort: 'created_at:desc',
        price_range: 'all',
        page: 1,
        pageSize: 12
      },
      resources: [],
      categories: [],
      total: 0,
      loading: false,
      error: ''
    }
  },
  computed: {
    totalPages() {
      return Math.ceil(this.total / this.searchQuery.pageSize)
    },
    currentPage() {
      return this.searchQuery.page
    }
  },
  created() {
    this.fetchCategories()
    this.searchResources()
  },
  methods: {
    async fetchCategories() {
      try {
        const response = await this.$axios.get('/api/resources/categories')
        this.categories = response.data
      } catch (error) {
        console.error('获取分类失败:', error)
      }
    },
    async searchResources() {
      this.loading = true
      this.error = ''
      
      try {
        const params = {
          q: this.searchQuery.q,
          category_id: this.searchQuery.category_id,
          tags: this.searchQuery.tags,
          sort: this.searchQuery.sort,
          price_range: this.searchQuery.price_range,
          page: this.searchQuery.page,
          pageSize: this.searchQuery.pageSize
        }
        
        const response = await this.$axios.get('/api/resources/search', { params })
        this.resources = response.data.resources
        this.total = response.data.total
      } catch (error) {
        console.error('搜索失败:', error)
        this.error = '搜索失败，请稍后重试'
      } finally {
        this.loading = false
      }
    },
    handleSearch() {
      this.searchQuery.page = 1
      this.searchResources()
    },
    changePage(page) {
      if (page > 0 && page <= this.totalPages) {
        this.searchQuery.page = page
        this.searchResources()
      }
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

.search-form {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.search-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.category-select, .sort-select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.search-button {
  padding: 8px 16px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.search-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.loading-indicator, .error-message, .no-results {
  text-align: center;
  padding: 20px;
}

.error-message {
  color: #f44336;
}

.resource-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.resource-item {
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 15px;
}

.resource-meta {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
  font-size: 0.8em;
  color: #666;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  margin-top: 20px;
}

.pagination button {
  padding: 5px 10px;
  border: 1px solid #ddd;
  background-color: #fff;
  cursor: pointer;
}

.pagination button:disabled {
  color: #999;
  cursor: not-allowed;
}
</style>