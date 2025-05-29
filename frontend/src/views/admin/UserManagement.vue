<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">用户管理</h1>
    </div>

    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <!-- 搜索和筛选 -->
      <div class="flex flex-col md:flex-row gap-4 mb-4">
        <div class="flex-1">
          <span class="p-input-icon-left w-full">
            <i class="pi pi-search" />
            <InputText v-model="filters.search" placeholder="搜索用户名或邮箱" class="w-full" @keyup.enter="loadUsers()" />
          </span>
        </div>
        <div class="flex gap-2">
          <Dropdown v-model="filters.role" :options="roleOptions" optionLabel="label" optionValue="value" placeholder="角色" class="w-32" />
          <Button label="搜索" icon="pi pi-search" @click="loadUsers()" />
          <Button label="重置" icon="pi pi-refresh" outlined @click="resetFilters()" />
        </div>
      </div>

      <!-- 用户表格 -->
      <DataTable 
        :value="users" 
        :loading="loading"
        :paginator="true" 
        :rows="pagination.pageSize"
        :totalRecords="pagination.total"
        :rowsPerPageOptions="[10, 20, 50]"
        v-model:first="pagination.first"
        v-model:rows="pagination.pageSize"
        @page="onPageChange($event)"
        dataKey="id"
        stripedRows
        responsiveLayout="scroll"
        class="p-datatable-sm"
      >
        <Column field="id" header="ID" :sortable="true" style="width: 80px"></Column>
        <Column field="username" header="用户名" :sortable="true">
          <template #body="{data}">
            <div class="flex items-center">
              <Avatar :image="data.avatar || '/default-avatar.png'" shape="circle" class="mr-2" size="small" />
              <span>{{ data.username }}</span>
            </div>
          </template>
        </Column>
        <Column field="email" header="邮箱" :sortable="true"></Column>
        <Column field="role" header="角色" :sortable="true">
          <template #body="{data}">
            <Tag :severity="data.role === 'admin' ? 'danger' : 'info'" :value="data.role === 'admin' ? '管理员' : '普通用户'" />
          </template>
        </Column>
        <Column field="points" header="积分" :sortable="true">
          <template #body="{data}">
            <span>{{ data.points || 0 }}</span>
          </template>
        </Column>
        <Column field="created_at" header="注册时间" :sortable="true">
          <template #body="{data}">
            {{ formatDate(data.created_at) }}
          </template>
        </Column>
        <Column header="操作" style="width: 150px">
          <template #body="{data}">
            <div class="flex gap-1">
              <Button icon="pi pi-user-edit" text rounded severity="info" @click="openEditDialog(data)" v-tooltip.top="'编辑角色'" />
              <Button icon="pi pi-trash" text rounded severity="danger" @click="confirmDelete(data)" v-tooltip.top="'删除用户'" />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- 编辑角色对话框 -->
    <Dialog v-model:visible="editDialog.visible" :header="'编辑用户角色 - ' + (editDialog.user?.username || '')" :style="{width: '450px'}" :modal="true">
      <div class="flex flex-col gap-4 p-4">
        <div class="field">
          <label for="role" class="block text-sm font-medium text-gray-700 mb-1">角色</label>
          <Dropdown id="role" v-model="editDialog.role" :options="roleOptions" optionLabel="label" optionValue="value" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" icon="pi pi-times" outlined @click="editDialog.visible = false" />
        <Button label="保存" icon="pi pi-check" @click="updateUserRole" :loading="editDialog.loading" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { adminApi } from '@/api/admin'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { format } from 'date-fns'
import { zhCN } from 'date-fns/locale'

import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Tag from 'primevue/tag'
import Avatar from 'primevue/avatar'

const router = useRouter()
const userStore = useUserStore()
const toast = useToast()
const confirm = useConfirm()

// 用户数据
const users = ref([])
const loading = ref(false)

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  first: 0
})

// 筛选条件
const filters = reactive({
  search: '',
  role: null
})

// 角色选项
const roleOptions = [
  { label: '全部角色', value: null },
  { label: '管理员', value: 'admin' },
  { label: '普通用户', value: 'user' }
]

// 编辑对话框
const editDialog = reactive({
  visible: false,
  user: null,
  role: '',
  loading: false
})

// 格式化日期
const formatDate = (date) => {
  if (!date) return ''
  return format(new Date(date), 'yyyy-MM-dd HH:mm', { locale: zhCN })
}

// 加载用户列表
const loadUsers = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      search: filters.search || undefined,
      role: filters.role || undefined
    }
    
    const response = await adminApi.getUsers(params)
    users.value = response.data.users || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('获取用户列表失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '获取用户列表失败', life: 3000 })
  } finally {
    loading.value = false
  }
}

// 页码变化
const onPageChange = (event) => {
  pagination.page = event.page + 1
  pagination.first = event.first
  pagination.pageSize = event.rows
  loadUsers()
}

// 重置筛选条件
const resetFilters = () => {
  filters.search = ''
  filters.role = null
  loadUsers()
}

// 打开编辑对话框
const openEditDialog = (user) => {
  editDialog.user = user
  editDialog.role = user.role
  editDialog.visible = true
}

// 更新用户角色
const updateUserRole = async () => {
  if (!editDialog.user) return
  
  editDialog.loading = true
  try {
    await adminApi.updateUserRole(editDialog.user.id, editDialog.role)
    toast.add({ severity: 'success', summary: '成功', detail: '用户角色已更新', life: 3000 })
    editDialog.visible = false
    loadUsers()
  } catch (error) {
    console.error('更新用户角色失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '更新用户角色失败', life: 3000 })
  } finally {
    editDialog.loading = false
  }
}

// 确认删除用户
const confirmDelete = (user) => {
  confirm.require({
    message: `确定要删除用户 "${user.username}" 吗？此操作不可恢复。`,
    header: '删除确认',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    acceptLabel: '删除',
    rejectLabel: '取消',
    accept: () => deleteUser(user.id)
  })
}

// 删除用户
const deleteUser = async (id) => {
  try {
    await adminApi.deleteUser(id)
    toast.add({ severity: 'success', summary: '成功', detail: '用户已删除', life: 3000 })
    loadUsers()
  } catch (error) {
    console.error('删除用户失败:', error)
    toast.add({ severity: 'error', summary: '错误', detail: '删除用户失败', life: 3000 })
  }
}

onMounted(() => {
  if (!userStore.isLoggedIn || userStore.user?.role !== 'admin') {
    router.push('/')
  } else {
    loadUsers()
  }
})
</script>