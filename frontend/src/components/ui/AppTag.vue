<template>
  <Tag
    :value="value"
    :severity="severity"
    :rounded="rounded"
    :icon="icon"
    :class="tagClass"
    v-bind="$attrs"
  />
</template>

<script setup>
import { computed } from 'vue'
import Tag from 'primevue/tag'

const props = defineProps({
  value: {
    type: String,
    default: ''
  },
  severity: {
    type: String,
    default: null,
    validator: (value) => [
      null, 'success', 'info', 'warning', 'danger'
    ].includes(value)
  },
  color: {
    type: String,
    default: ''
  },
  icon: {
    type: String,
    default: ''
  },
  rounded: {
    type: Boolean,
    default: false
  },
  size: {
    type: String,
    default: 'medium',
    validator: (value) => ['small', 'medium', 'large'].includes(value)
  },
  outlined: {
    type: Boolean,
    default: false
  },
  removable: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['remove'])

// 计算标签类名
const tagClass = computed(() => {
  const classes = []
  
  // 尺寸
  switch (props.size) {
    case 'small':
      classes.push('text-xs py-1 px-2')
      break
    case 'large':
      classes.push('text-base py-2 px-4')
      break
  }
  
  // 自定义颜色
  if (props.color && !props.severity) {
    if (props.outlined) {
      classes.push(`border border-${props.color} text-${props.color} bg-transparent`)
    } else {
      classes.push(`bg-${props.color} text-white`)
    }
  }
  
  // 轮廓样式
  if (props.outlined && !props.color) {
    classes.push('p-tag-outlined')
  }
  
  // 可移除
  if (props.removable) {
    classes.push('pr-8 relative')
  }
  
  return classes.join(' ')
})
</script>

<style scoped>
.p-tag {
  display: inline-flex;
  align-items: center;
}

/* 可移除标签的样式 */
.relative {
  position: relative;
}

.relative::after {
  content: '×';
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  cursor: pointer;
  font-size: 1.2em;
}
</style>