<template>
  <form @submit.prevent="handleSubmit" :class="formClass">
    <div v-if="title" class="mb-4">
      <h3 class="text-xl font-medium">{{ title }}</h3>
      <p v-if="subtitle" class="text-gray-600 mt-1">{{ subtitle }}</p>
    </div>
    
    <slot></slot>
    
    <div v-if="$slots.actions || submitLabel || resetLabel" class="form-actions mt-6 flex items-center" :class="{ 'justify-end': actionsAlign === 'right', 'justify-between': actionsAlign === 'between' }">
      <slot name="actions">
        <div class="flex gap-3">
          <AppButton 
            v-if="resetLabel" 
            :label="resetLabel" 
            variant="outline" 
            type="button"
            @click="handleReset"
          />
          <AppButton 
            v-if="submitLabel" 
            :label="submitLabel" 
            :loading="loading"
            :disabled="disabled || loading"
            type="submit"
          />
        </div>
      </slot>
    </div>
  </form>
</template>

<script setup>
import { computed } from 'vue'
import AppButton from './AppButton.vue'

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  submitLabel: {
    type: String,
    default: ''
  },
  resetLabel: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  layout: {
    type: String,
    default: 'vertical',
    validator: (value) => ['vertical', 'horizontal', 'inline'].includes(value)
  },
  actionsAlign: {
    type: String,
    default: 'right',
    validator: (value) => ['left', 'right', 'between', 'center'].includes(value)
  },
  gap: {
    type: String,
    default: 'normal',
    validator: (value) => ['small', 'normal', 'large'].includes(value)
  }
})

const emit = defineEmits(['submit', 'reset'])

// 处理表单提交
const handleSubmit = (event) => {
  emit('submit', event)
}

// 处理表单重置
const handleReset = () => {
  emit('reset')
}

// 计算表单类名
const formClass = computed(() => {
  const classes = []
  
  // 布局
  switch (props.layout) {
    case 'vertical':
      classes.push('flex flex-col')
      break
    case 'horizontal':
      classes.push('grid grid-cols-12 gap-4')
      break
    case 'inline':
      classes.push('flex flex-row items-center')
      break
  }
  
  // 间距
  switch (props.gap) {
    case 'small':
      classes.push('space-y-2')
      break
    case 'normal':
      classes.push('space-y-4')
      break
    case 'large':
      classes.push('space-y-6')
      break
  }
  
  return classes.join(' ')
})
</script>