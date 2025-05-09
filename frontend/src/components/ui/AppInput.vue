<template>
  <div class="app-input-wrapper">
    <label v-if="label" :for="id" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500 ml-1">*</span>
    </label>
    <div class="relative">
      <span v-if="iconLeft" class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">
        <i :class="`pi ${iconLeft}`"></i>
      </span>
      <InputText
        :id="id"
        v-model="inputValue"
        :type="type"
        :placeholder="placeholder"
        :disabled="disabled"
        :class="inputClass"
        v-bind="$attrs"
        @input="handleInput"
        @blur="handleBlur"
      />
      <span v-if="iconRight" class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-500">
        <i :class="`pi ${iconRight}`"></i>
      </span>
    </div>
    <small v-if="helperText" class="block mt-1 text-sm text-gray-500">{{ helperText }}</small>
    <small v-if="errorMessage" class="block mt-1 text-sm text-red-500">{{ errorMessage }}</small>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import InputText from 'primevue/inputtext'

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  label: {
    type: String,
    default: ''
  },
  id: {
    type: String,
    default: () => `input-${Math.random().toString(36).substring(2, 9)}`
  },
  type: {
    type: String,
    default: 'text'
  },
  placeholder: {
    type: String,
    default: ''
  },
  required: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  iconLeft: {
    type: String,
    default: ''
  },
  iconRight: {
    type: String,
    default: ''
  },
  helperText: {
    type: String,
    default: ''
  },
  errorMessage: {
    type: String,
    default: ''
  },
  variant: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'filled', 'outlined'].includes(value)
  },
  size: {
    type: String,
    default: 'medium',
    validator: (value) => ['small', 'medium', 'large'].includes(value)
  },
  fullWidth: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'input', 'blur'])

const inputValue = ref(props.modelValue)

// 监听modelValue变化
watch(() => props.modelValue, (newValue) => {
  inputValue.value = newValue
})

// 监听inputValue变化
watch(() => inputValue.value, (newValue) => {
  emit('update:modelValue', newValue)
})

// 处理输入事件
const handleInput = (event) => {
  emit('input', event)
}

// 处理失焦事件
const handleBlur = (event) => {
  emit('blur', event)
}

// 计算输入框类名
const inputClass = computed(() => {
  const classes = []
  
  // 左侧图标padding
  if (props.iconLeft) {
    classes.push('pl-10')
  }
  
  // 右侧图标padding
  if (props.iconRight) {
    classes.push('pr-10')
  }
  
  // 尺寸
  switch (props.size) {
    case 'small':
      classes.push('p-inputtext-sm')
      break
    case 'large':
      classes.push('p-inputtext-lg')
      break
  }
  
  // 变体样式
  switch (props.variant) {
    case 'filled':
      classes.push('bg-gray-100 border-gray-300')
      break
    case 'outlined':
      classes.push('border-2 border-gray-300')
      break
  }
  
  // 错误状态
  if (props.errorMessage) {
    classes.push('p-invalid border-red-500')
  }
  
  // 宽度
  if (props.fullWidth) {
    classes.push('w-full')
  }
  
  return classes.join(' ')
})
</script>