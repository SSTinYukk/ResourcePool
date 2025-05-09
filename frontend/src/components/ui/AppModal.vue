<template>
  <Dialog
    v-model:visible="isVisible"
    :modal="modal"
    :header="title"
    :style="{ width: width }"
    :breakpoints="{ '960px': '75vw', '640px': '90vw' }"
    :draggable="draggable"
    :closable="closable"
    :closeOnEscape="closeOnEscape"
    :showHeader="showHeader"
    :baseZIndex="baseZIndex"
    :class="modalClass"
    @hide="handleHide"
  >
    <template #header v-if="$slots.header && showHeader">
      <slot name="header"></slot>
    </template>
    
    <div class="modal-content">
      <slot></slot>
    </div>
    
    <template #footer v-if="$slots.footer || (showFooter && (confirmLabel || cancelLabel))">
      <slot name="footer">
        <div class="flex justify-end gap-2">
          <AppButton
            v-if="cancelLabel"
            :label="cancelLabel"
            variant="outline"
            @click="handleCancel"
          />
          <AppButton
            v-if="confirmLabel"
            :label="confirmLabel"
            :loading="loading"
            @click="handleConfirm"
          />
        </div>
      </slot>
    </template>
  </Dialog>
</template>

<script setup>
import { computed, watch, ref } from 'vue'
import Dialog from 'primevue/dialog'
import AppButton from './AppButton.vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  width: {
    type: String,
    default: '500px'
  },
  modal: {
    type: Boolean,
    default: true
  },
  draggable: {
    type: Boolean,
    default: false
  },
  closable: {
    type: Boolean,
    default: true
  },
  closeOnEscape: {
    type: Boolean,
    default: true
  },
  showHeader: {
    type: Boolean,
    default: true
  },
  showFooter: {
    type: Boolean,
    default: true
  },
  confirmLabel: {
    type: String,
    default: '确认'
  },
  cancelLabel: {
    type: String,
    default: '取消'
  },
  loading: {
    type: Boolean,
    default: false
  },
  baseZIndex: {
    type: Number,
    default: 1000
  },
  position: {
    type: String,
    default: 'center',
    validator: (value) => ['center', 'top', 'bottom', 'left', 'right', 'topleft', 'topright', 'bottomleft', 'bottomright'].includes(value)
  }
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel', 'hide'])

// 内部可见状态
const isVisible = ref(props.modelValue)

// 监听modelValue变化
watch(() => props.modelValue, (newValue) => {
  isVisible.value = newValue
})

// 监听isVisible变化
watch(() => isVisible.value, (newValue) => {
  emit('update:modelValue', newValue)
})

// 处理确认事件
const handleConfirm = () => {
  emit('confirm')
}

// 处理取消事件
const handleCancel = () => {
  isVisible.value = false
  emit('cancel')
}

// 处理隐藏事件
const handleHide = () => {
  emit('hide')
}

// 计算模态框类名
const modalClass = computed(() => {
  const classes = []
  
  // 位置
  if (props.position !== 'center') {
    classes.push(`p-dialog-${props.position}`)
  }
  
  return classes.join(' ')
})
</script>

<style scoped>
.modal-content {
  padding: 0 1rem;
  max-height: 70vh;
  overflow-y: auto;
}

/* 自定义位置样式 */
:deep(.p-dialog-top) {
  margin-top: 2rem;
}

:deep(.p-dialog-bottom) {
  margin-bottom: 2rem;
}

:deep(.p-dialog-left) {
  margin-left: 2rem;
}

:deep(.p-dialog-right) {
  margin-right: 2rem;
}

:deep(.p-dialog-topleft) {
  margin-top: 2rem;
  margin-left: 2rem;
}

:deep(.p-dialog-topright) {
  margin-top: 2rem;
  margin-right: 2rem;
}

:deep(.p-dialog-bottomleft) {
  margin-bottom: 2rem;
  margin-left: 2rem;
}

:deep(.p-dialog-bottomright) {
  margin-bottom: 2rem;
  margin-right: 2rem;
}
</style>