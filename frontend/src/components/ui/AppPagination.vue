<template>
  <div class="app-pagination flex items-center justify-center" :class="{ 'justify-between': showInfo }">
    <div v-if="showInfo" class="pagination-info text-sm text-gray-500">
      显示 {{ firstItem }}-{{ lastItem }} 条，共 {{ totalRecords }} 条
    </div>
    <Paginator
      :rows="rows"
      :totalRecords="totalRecords"
      :first="first"
      :pageLinkSize="pageLinkSize"
      :template="template"
      :rowsPerPageOptions="rowsPerPageOptions"
      :alwaysShow="alwaysShow"
      :class="paginatorClass"
      @page="onPageChange"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import Paginator from 'primevue/paginator'

const props = defineProps({
  currentPage: {
    type: Number,
    default: 1
  },
  rows: {
    type: Number,
    default: 10
  },
  totalRecords: {
    type: Number,
    default: 0
  },
  pageLinkSize: {
    type: Number,
    default: 5
  },
  rowsPerPageOptions: {
    type: Array,
    default: () => [10, 20, 30, 50]
  },
  template: {
    type: String,
    default: 'FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown'
  },
  alwaysShow: {
    type: Boolean,
    default: true
  },
  showInfo: {
    type: Boolean,
    default: true
  },
  size: {
    type: String,
    default: 'default',
    validator: (value) => ['small', 'default', 'large'].includes(value)
  }
})

const emit = defineEmits(['update:currentPage', 'update:rows', 'page'])

// 计算当前页的第一条记录
const firstItem = computed(() => {
  if (props.totalRecords === 0) return 0
  return (props.currentPage - 1) * props.rows + 1
})

// 计算当前页的最后一条记录
const lastItem = computed(() => {
  const last = props.currentPage * props.rows
  return last > props.totalRecords ? props.totalRecords : last
})

// 计算PrimeVue Paginator的first属性（基于0的索引）
const first = computed(() => {
  return (props.currentPage - 1) * props.rows
})

// 处理页面变化事件
const onPageChange = (event) => {
  const newPage = Math.floor(event.first / event.rows) + 1
  emit('update:currentPage', newPage)
  emit('update:rows', event.rows)
  emit('page', {
    page: newPage,
    first: event.first,
    rows: event.rows,
    pageCount: Math.ceil(props.totalRecords / event.rows)
  })
}

// 计算分页器类名
const paginatorClass = computed(() => {
  const classes = []
  
  // 尺寸
  switch (props.size) {
    case 'small':
      classes.push('p-paginator-sm')
      break
    case 'large':
      classes.push('p-paginator-lg')
      break
  }
  
  return classes.join(' ')
})
</script>