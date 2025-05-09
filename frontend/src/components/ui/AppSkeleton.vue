<template>
  <div :class="['app-skeleton', { 'animate-pulse': animation === 'pulse' }]">
    <div v-if="type === 'text'" 
      :class="['bg-gray-200 rounded', sizeClass, widthClass]"
      :style="customStyle">
    </div>
    
    <div v-else-if="type === 'circle'" 
      :class="['bg-gray-200 rounded-full', sizeClass]"
      :style="customStyle">
    </div>
    
    <div v-else-if="type === 'rectangle'" 
      :class="['bg-gray-200 rounded', sizeClass, widthClass]"
      :style="customStyle">
    </div>
    
    <div v-else-if="type === 'card'" class="bg-gray-200 rounded p-4">
      <div class="flex items-center mb-4">
        <div class="w-12 h-12 bg-gray-300 rounded-full mr-3"></div>
        <div class="flex-1">
          <div class="h-4 bg-gray-300 rounded w-3/4 mb-2"></div>
          <div class="h-3 bg-gray-300 rounded w-1/2"></div>
        </div>
      </div>
      <div class="space-y-2">
        <div class="h-3 bg-gray-300 rounded"></div>
        <div class="h-3 bg-gray-300 rounded"></div>
        <div class="h-3 bg-gray-300 rounded w-4/5"></div>
      </div>
    </div>
    
    <div v-else-if="type === 'list'" class="space-y-3">
      <div v-for="i in count" :key="i" class="flex items-center">
        <div v-if="showAvatar" class="w-10 h-10 bg-gray-300 rounded-full mr-3"></div>
        <div class="flex-1">
          <div class="h-4 bg-gray-300 rounded w-3/4 mb-2"></div>
          <div class="h-3 bg-gray-300 rounded w-1/2"></div>
        </div>
      </div>
    </div>
    
    <div v-else-if="type === 'table'" class="w-full">
      <div class="flex mb-3">
        <div v-for="i in columns" :key="i" class="flex-1 h-8 bg-gray-300 rounded mr-2"></div>
      </div>
      <div v-for="i in count" :key="i" class="flex mb-2">
        <div v-for="j in columns" :key="j" class="flex-1 h-6 bg-gray-200 rounded mr-2"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'text',
    validator: (value) => ['text', 'circle', 'rectangle', 'card', 'list', 'table'].includes(value)
  },
  size: {
    type: String,
    default: 'medium',
    validator: (value) => ['small', 'medium', 'large', 'custom'].includes(value)
  },
  width: {
    type: String,
    default: 'full',
    validator: (value) => ['full', 'half', 'quarter', 'auto', 'custom'].includes(value)
  },
  height: {
    type: String,
    default: ''
  },
  customWidth: {
    type: String,
    default: ''
  },
  animation: {
    type: String,
    default: 'pulse',
    validator: (value) => ['pulse', 'none'].includes(value)
  },
  count: {
    type: Number,
    default: 3
  },
  columns: {
    type: Number,
    default: 4
  },
  showAvatar: {
    type: Boolean,
    default: true
  }
})

// 计算尺寸类名
const sizeClass = computed(() => {
  switch (props.type) {
    case 'text':
      switch (props.size) {
        case 'small': return 'h-3'
        case 'medium': return 'h-4'
        case 'large': return 'h-6'
        case 'custom': return ''
        default: return 'h-4'
      }
    case 'circle':
      switch (props.size) {
        case 'small': return 'w-8 h-8'
        case 'medium': return 'w-12 h-12'
        case 'large': return 'w-16 h-16'
        case 'custom': return ''
        default: return 'w-12 h-12'
      }
    case 'rectangle':
      switch (props.size) {
        case 'small': return 'h-16'
        case 'medium': return 'h-24'
        case 'large': return 'h-32'
        case 'custom': return ''
        default: return 'h-24'
      }
    default:
      return ''
  }
})

// 计算宽度类名
const widthClass = computed(() => {
  if (props.width === 'custom') return ''
  
  switch (props.width) {
    case 'full': return 'w-full'
    case 'half': return 'w-1/2'
    case 'quarter': return 'w-1/4'
    case 'auto': return 'w-auto'
    default: return 'w-full'
  }
})

// 自定义样式
const customStyle = computed(() => {
  const style = {}
  
  if (props.size === 'custom' && props.height) {
    style.height = props.height
  }
  
  if (props.width === 'custom' && props.customWidth) {
    style.width = props.customWidth
  }
  
  return style
})
</script>

<style scoped>
.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>