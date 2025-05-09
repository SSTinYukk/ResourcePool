<template>
  <div class="app-table-wrapper">
    <div v-if="title || $slots.header" class="flex justify-between items-center mb-4">
      <h3 v-if="title" class="text-lg font-medium">{{ title }}</h3>
      <slot name="header"></slot>
    </div>
    
    <DataTable
      :value="data"
      :paginator="paginator"
      :rows="rows"
      :rowsPerPageOptions="rowsPerPageOptions"
      :loading="loading"
      :filters="filters"
      :globalFilterFields="globalFilterFields"
      :emptyMessage="emptyMessage"
      :class="tableClass"
      v-bind="$attrs"
      @row-click="handleRowClick"
    >
      <template v-if="$slots.empty && !data?.length" #empty>
        <slot name="empty"></slot>
      </template>
      
      <slot></slot>
      
      <template v-if="$slots.paginatorLeft" #paginatorLeft>
        <slot name="paginatorLeft"></slot>
      </template>
      
      <template v-if="$slots.paginatorRight" #paginatorRight>
        <slot name="paginatorRight"></slot>
      </template>
    </DataTable>
    
    <div v-if="$slots.footer" class="mt-4">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import DataTable from 'primevue/datatable'

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  },
  title: {
    type: String,
    default: ''
  },
  paginator: {
    type: Boolean,
    default: false
  },
  rows: {
    type: Number,
    default: 10
  },
  rowsPerPageOptions: {
    type: Array,
    default: () => [5, 10, 20, 50]
  },
  loading: {
    type: Boolean,
    default: false
  },
  filters: {
    type: Object,
    default: () => ({})
  },
  globalFilterFields: {
    type: Array,
    default: () => []
  },
  emptyMessage: {
    type: String,
    default: '暂无数据'
  },
  striped: {
    type: Boolean,
    default: true
  },
  bordered: {
    type: Boolean,
    default: false
  },
  hoverable: {
    type: Boolean,
    default: true
  },
  size: {
    type: String,
    default: 'default',
    validator: (value) => ['small', 'default', 'large'].includes(value)
  },
  responsive: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['row-click'])

// 处理行点击事件
const handleRowClick = (event) => {
  emit('row-click', event)
}

// 计算表格类名
const tableClass = computed(() => {
  const classes = []
  
  // 条纹样式
  if (props.striped) {
    classes.push('p-datatable-striped')
  }
  
  // 边框
  if (props.bordered) {
    classes.push('p-datatable-bordered')
  }
  
  // 悬停效果
  if (props.hoverable) {
    classes.push('p-datatable-hoverable-rows')
  }
  
  // 尺寸
  switch (props.size) {
    case 'small':
      classes.push('p-datatable-sm')
      break
    case 'large':
      classes.push('p-datatable-lg')
      break
  }
  
  // 响应式
  if (props.responsive) {
    classes.push('p-datatable-responsive')
  }
  
  return classes.join(' ')
})
</script>