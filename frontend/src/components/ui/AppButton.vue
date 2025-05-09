<template>
  <Button
    :label="label"
    :icon="icon"
    :iconPos="iconPos"
    :class="buttonClass"
    :disabled="disabled"
    :loading="loading"
    v-bind="$attrs"
    @click="$emit('click', $event)"
  />
</template>

<script setup>
import { computed } from 'vue'
import Button from 'primevue/button'

const props = defineProps({
  label: {
    type: String,
    default: ''
  },
  icon: {
    type: String,
    default: ''
  },
  iconPos: {
    type: String,
    default: 'left'
  },
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => [
      'primary', 'secondary', 'success', 'info', 'warning', 'danger', 'light', 'dark', 'outline', 'text'
    ].includes(value)
  },
  size: {
    type: String,
    default: 'normal',
    validator: (value) => ['small', 'normal', 'large'].includes(value)
  },
  rounded: {
    type: Boolean,
    default: false
  },
  block: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])

// 计算按钮类名
const buttonClass = computed(() => {
  const classes = []
  
  // 变体样式
  switch (props.variant) {
    case 'primary':
      classes.push('p-button-primary')
      break
    case 'secondary':
      classes.push('p-button-secondary')
      break
    case 'success':
      classes.push('p-button-success')
      break
    case 'info':
      classes.push('p-button-info')
      break
    case 'warning':
      classes.push('p-button-warning')
      break
    case 'danger':
      classes.push('p-button-danger')
      break
    case 'light':
      classes.push('p-button-light')
      break
    case 'dark':
      classes.push('p-button-dark')
      break
    case 'outline':
      classes.push('p-button-outlined')
      break
    case 'text':
      classes.push('p-button-text')
      break
  }
  
  // 尺寸
  if (props.size === 'small') {
    classes.push('p-button-sm')
  } else if (props.size === 'large') {
    classes.push('p-button-lg')
  }
  
  // 圆角
  if (props.rounded) {
    classes.push('p-button-rounded')
  }
  
  // 块级
  if (props.block) {
    classes.push('w-full')
  }
  
  return classes.join(' ')
})
</script>