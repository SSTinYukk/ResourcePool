<template>
  <div :class="cardClass">
    <div v-if="$slots.header" class="card-header p-4 border-b border-gray-200">
      <slot name="header"></slot>
    </div>
    <div class="card-body p-4">
      <h3 v-if="title" class="text-xl font-medium mb-2">{{ title }}</h3>
      <p v-if="subtitle" class="text-gray-600 mb-4">{{ subtitle }}</p>
      <slot></slot>
    </div>
    <div v-if="$slots.footer" class="card-footer p-4 border-t border-gray-200">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  variant: {
    type: String,
    default: 'default',
    validator: (value) => [
      'default', 'primary', 'success', 'info', 'warning', 'danger'
    ].includes(value)
  },
  elevation: {
    type: Number,
    default: 1,
    validator: (value) => value >= 0 && value <= 5
  },
  rounded: {
    type: Boolean,
    default: true
  },
  bordered: {
    type: Boolean,
    default: false
  },
  hoverable: {
    type: Boolean,
    default: false
  },
  fullWidth: {
    type: Boolean,
    default: false
  }
})

// 计算卡片类名
const cardClass = computed(() => {
  const classes = ['bg-white']
  
  // 圆角
  if (props.rounded) {
    classes.push('rounded-lg')
  }
  
  // 边框
  if (props.bordered) {
    classes.push('border border-gray-200')
  }
  
  // 阴影
  if (props.elevation > 0) {
    classes.push(`shadow-${props.elevation === 1 ? '' : props.elevation}`)
  }
  
  // 悬停效果
  if (props.hoverable) {
    classes.push('transition-shadow duration-300 hover:shadow-lg')
  }
  
  // 宽度
  if (props.fullWidth) {
    classes.push('w-full')
  }
  
  // 变体样式
  switch (props.variant) {
    case 'primary':
      classes.push('border-l-4 border-l-blue-500')
      break
    case 'success':
      classes.push('border-l-4 border-l-green-500')
      break
    case 'info':
      classes.push('border-l-4 border-l-cyan-500')
      break
    case 'warning':
      classes.push('border-l-4 border-l-yellow-500')
      break
    case 'danger':
      classes.push('border-l-4 border-l-red-500')
      break
  }
  
  return classes.join(' ')
})
</script>