<template>
  <div class="app-select-wrapper">
    <label v-if="label" :for="id" class="block text-sm font-medium text-gray-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500 ml-1">*</span>
    </label>
    <div class="relative">
      <Dropdown
        :id="id"
        v-model="selectedValue"
        :options="options"
        :optionLabel="optionLabel"
        :optionValue="optionValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :class="selectClass"
        :filter="filterable"
        :showClear="clearable"
        v-bind="$attrs"
        @change="handleChange"
      />
    </div>
    <small v-if="helperText" class="block mt-1 text-sm text-gray-500">{{ helperText }}</small>
    <small v-if="errorMessage" class="block mt-1 text-sm text-red-500">{{ errorMessage }}</small>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import Dropdown from 'primevue/dropdown'

const props = defineProps({
  modelValue: {
    type: [String, Number, Object, Array],
    default: null
  },
  options: {
    type: Array,
    default: () => []
  },
  optionLabel: {
    type: String,
    default: 'label'
  },
  optionValue: {
    type: String,
    default: 'value'
  },
  label: {
    type: String,
    default: ''
  },
  id: {
    type: String,
    default: () => `select-${Math.random().toString(36).substring(2, 9)}`
  },
  placeholder: {
    type: String,
    default: '请选择'
  },
  required: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  filterable: {
    type: Boolean,
    default: false
  },
  clearable: {
    type: Boolean,
    default: false
  },
  helperText: {
    type: String,
    default: ''
  },
  errorMessage: {
    type: String,
    default: ''
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

const emit = defineEmits(['update:modelValue', 'change'])

const selectedValue = ref(props.modelValue)

// 监听modelValue变化
watch(() => props.modelValue, (newValue) => {
  selectedValue.value = newValue
})

// 监听selectedValue变化
watch(() => selectedValue.value, (newValue) => {
  emit('update:modelValue', newValue)
})

// 处理选择变化事件
const handleChange = (event) => {
  emit('change', event)
}

// 计算选择框类名
const selectClass = computed(() => {
  const classes = []
  
  // 尺寸
  switch (props.size) {
    case 'small':
      classes.push('p-inputtext-sm')
      break
    case 'large':
      classes.push('p-inputtext-lg')
      break
  }
  
  // 错误状态
  if (props.errorMessage) {
    classes.push('p-invalid')
  }
  
  // 宽度
  if (props.fullWidth) {
    classes.push('w-full')
  }
  
  return classes.join(' ')
})
</script>