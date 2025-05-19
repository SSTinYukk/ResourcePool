<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">论坛主题</h1>
      <router-link 
        to="/forum/create-topic"
        class="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition duration-300"
        v-if="userStore.isLoggedIn"
      >
        创建新主题
      </router-link>
    </div>

    <div class="bg-white rounded-lg shadow-md p-6">
      <div class="mb-4">
        <InputText 
          v-model="searchQuery" 
          placeholder="搜索主题..." 
          class="w-full"
        />
      </div>

      <DataTable 
        :value="topics" 
        :paginator="true" 
        :rows="10"
        :loading="loading"
        paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink CurrentPageReport"
        currentPageReportTemplate="显示 {first} 到 {last} 共 {totalRecords} 条"
      >
        <Column field="title" header="标题">
          <template #body="{data}">
            <router-link 
              :to="`/forum/topics/${data.id}`" 
              class="text-blue-600 hover:underline"
            >
              {{ data.title }}
            </router-link>
          </template>
        </Column>
        <Column field="author" header="作者" />
        <Column field="category" header="分类" />
        <Column field="replies" header="回复" />
        <Column field="views" header="浏览" />
        <Column field="lastReply" header="最后回复">
          <template #body="{data}">
            {{ formatDate(data.lastReply) }}
          </template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import InputText from 'primevue/inputtext'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const userStore = useUserStore()
const topics = ref([])
const loading = ref(true)
const searchQuery = ref('')

const fetchTopics = async () => {
  try {
    loading.value = true
    const response = await forumApi.getTopics();
    topics.value = response.data;
  } catch (error) {
    console.error('获取主题列表失败:', error)
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})
}

onMounted(() => {
  fetchTopics()
})
</script>