<template>
  <div v-if="visible" :class="messageClass" class="app-message p-4 mb-4 flex items-start">
    <i v-if="showIcon" :class="iconClass" class="mr-3 text-lg"></i>
    <div class="flex-grow">
      <div v-if="title" class="font-medium mb-1">{{ title }}</div>
      <div class="message-content">
        <slot>{{ content }}</slot>
      </div>
    </div>
    <button v-if="closable" @click="handleClose" class="ml-3 text-gray-500 hover:text-gray-700">
      <i class="pi pi-times"></i>
    </button>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'info',
    validator: (value) => ['success', 'info', 'warning', 'error'].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  content: {
    type: String,
    default: ''
  },
  showIcon: {
    type: Boolean,
    default: true
  },
  closable: {
    type: Boolean,
    default: false
  },
  autoClose: {
    type: Boolean,
    default: false
  },
  duration: {
    type: Number,
    default: 3000
  },
  bordered: {
    type: Boolean,
    default: true
  },
  filled: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const visible = ref(true)
let timer = null

// 监听autoClose和duration变化
watch(
  () => [props.autoClose, props.duration],
  ([autoClose, duration]) => {
    clearTimeout(timer)
    if (autoClose && visible.value) {
      timer = setTimeout(() => {
        visible.value = false
        emit('close')
      }, duration)
    }
  },
  { immediate: true }
)

// 处理关闭事件
const handleClose = () => {
  visible.value = false
  emit('close')
}

// 计算消息类名
const messageClass = computed(() => {
  const classes = ['rounded-md']
  
  if (props.bordered) {
    classes.push('border')
  }
  
  // 根据类型设置样式
  switch (props.type) {
    case 'success':
      classes.push(props.filled ? 'bg-green-500 text-white' : 'bg-green-50 text-green-800 border-green-200')
      break
    case 'info':
      classes.push(props.filled ? 'bg-blue-500 text-white' : 'bg-blue-50 text-blue-800 border-blue-200')
      break
    case 'warning':
      classes.push(props.filled ? 'bg-yellow-500 text-white' : 'bg-yellow-50 text-yellow-800 border-yellow-200')
      break
    case 'error':
      classes.push(props.filled ? 'bg-red-500 text-white' : 'bg-red-50 text-red-800 border-red-200')
      break
  }
  
  return classes.join(' ')
})

// 计算图标类名
const iconClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'pi pi-check-circle'
    case 'info':
      return 'pi pi-info-circle'
    case 'warning':
      return 'pi pi-exclamation-triangle'
    case 'error':
      return 'pi pi-times-circle'
    default:
      return 'pi pi-info-circle'
  }
})
</script>