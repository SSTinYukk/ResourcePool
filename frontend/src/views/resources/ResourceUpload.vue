<template>
  <div class="container mx-auto py-8 px-4">
    <Card>
      <template #title>
        <div class="flex items-center">
          <i class="pi pi-upload mr-2 text-blue-500"></i>
          <h1 class="text-xl font-bold">上传资源</h1>
        </div>
      </template>
      <template #content>
        <form @submit.prevent="submitForm" class="space-y-6">
          <!-- 基本信息 -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="col-span-2">
              <AppInput
                v-model="form.title"
                label="资源标题"
                placeholder="请输入资源标题"
                required
                :error-message="errors.title"
                @blur="validateField('title')"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                资源分类 <span class="text-red-500">*</span>
              </label>
              <Dropdown
                v-model="form.category_id"
                :options="categories"
                optionLabel="name"
                optionValue="value"
                placeholder="选择资源分类"
                class="w-full"
                :class="{ 'p-invalid': errors.category_id_id }"
                @change="validateField('category_id')"
              />
              <small v-if="errors.category_id" class="p-error block mt-1">{{ errors.category_id }}</small>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                资源类型 <span class="text-red-500">*</span>
              </label>
              <Dropdown
                v-model="form.fileType"
                :options="fileTypes"
                optionLabel="name"
                optionValue="value"
                placeholder="选择资源类型"
                class="w-full"
                :class="{ 'p-invalid': errors.fileType }"
                @change="validateField('fileType')"
              />
              <small v-if="errors.fileType" class="p-error block mt-1">{{ errors.fileType }}</small>
            </div>
          </div>
          
          <!-- 资源描述 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              资源描述 <span class="text-red-500">*</span>
            </label>
            <Editor
              v-model="form.description"
              editorStyle="height: 250px"
              :class="{ 'p-invalid': errors.description }"
              @blur="validateField('description')"
            />
            <small v-if="errors.description" class="p-error block mt-1">{{ errors.description }}</small>
          </div>
          
          <!-- 文件上传 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              上传文件 <span class="text-red-500">*</span>
            </label>
            <FileUpload
              ref="fileUpload"
              mode="basic"
              :multiple="false"
              accept=".pdf,.doc,.docx,.ppt,.pptx,.zip,.rar,.jpg,.png,.mp4"
              :maxFileSize="50000000"
              chooseLabel="选择文件"
              class="w-full"
              :class="{ 'p-invalid': errors.file }"
              @select="onFileSelect"
              @clear="form.file = null"
            />
            <small class="text-gray-500 block mt-1">支持的文件格式：PDF, Word, PPT, ZIP, RAR, 图片, 视频等，最大50MB</small>
            <small v-if="errors.file" class="p-error block mt-1">{{ errors.file }}</small>
          </div>
          
          <!-- 预览图片上传 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              预览图片
            </label>
            <FileUpload
              ref="previewUpload"
              mode="advanced"
              :multiple="true"
              accept="image/*"
              :maxFileSize="5000000"
              chooseLabel="选择图片"
              uploadLabel="上传"
              cancelLabel="取消"
              :auto="true"
              :customUpload="true"
              @uploader="onPreviewUpload"
              :class="{ 'p-invalid': errors.previewImages }"
            />
            <small class="text-gray-500 block mt-1">上传资源的预览图片，最多5张，每张不超过5MB</small>
            <small v-if="errors.previewImages" class="p-error block mt-1">{{ errors.previewImages }}</small>
          </div>
          
          <!-- 标签 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              资源标签
            </label>
            <Chips v-model="form.tags" placeholder="输入标签后按回车" :max="5" />
            <small class="text-gray-500 block mt-1">最多添加5个标签，每个标签按回车确认</small>
          </div>
          
          <!-- 提交按钮 -->
          <div class="flex justify-end space-x-3">
            <Button type="button" label="取消" class="p-button-outlined" @click="$router.back()" />
            <Button type="submit" label="上传资源" icon="pi pi-upload" :loading="loading" />
          </div>
        </form>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import axios from 'axios'
import { useUserStore } from '@/stores/user'
import Card from 'primevue/card'
import Dropdown from 'primevue/dropdown'
import Editor from 'primevue/editor'
import FileUpload from 'primevue/fileupload'
import Chips from 'primevue/chips'
import Button from 'primevue/button'
import AppInput from '@/components/ui/AppInput.vue'

const router = useRouter()
const toast = useToast()
const userStore = useUserStore()
const loading = ref(false)
const fileUpload = ref(null)
const previewUpload = ref(null)

// 表单数据
const form = reactive({
  title: '',
  category: '',
  fileType: '',
  description: '',
  file: null,
  previewImages: [],
  tags: [],
  category_id: ''
})

// 错误信息
const errors = reactive({
  title: '',
  category: '',
  fileType: '',
  description: '',
  file: '',
  previewImages: ''
})

// 分类选项
const categories = [
  { name: '教程资料', value: 1 },
  { name: '软件工具', value: 2 },
  { name: '学习资源', value: 3 },
  { name: '模板素材', value: 4 },
  { name: '其他资源', value: 5 }
]

// 文件类型选项
const fileTypes = [
  { name: 'PDF文档', value: 'pdf' },
  { name: 'Word文档', value: 'word' },
  { name: 'PPT演示', value: 'ppt' },
  { name: '压缩包', value: 'archive' },
  { name: '图片', value: 'image' },
  { name: '视频', value: 'video' },
  { name: '其他', value: 'other' }
]

// 文件选择处理
const onFileSelect = (event) => {
  if (event.files && event.files.length > 0) {
    form.file = event.files[0]
    validateField('file')
  }
}

// 预览图片上传处理
const onPreviewUpload = (event) => {
  // 这里应该调用实际的上传API，这里仅做模拟
  const file = event.files[0]
  
  // 模拟上传成功后添加到预览图片数组
  if (form.previewImages.length < 5) {
    // 实际项目中，这里应该是上传后返回的图片URL
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => {
      form.previewImages.push(reader.result)
    }
    
    // 显示上传成功消息
    toast.add({
      severity: 'success',
      summary: '上传成功',
      detail: `图片 ${file.name} 上传成功`,
      life: 3000
    })
  } else {
    toast.add({
      severity: 'error',
      summary: '上传失败',
      detail: '最多只能上传5张预览图片',
      life: 3000
    })
  }
  
  // 清除上传组件的文件
  previewUpload.value.clear()
}

// 字段验证
const validateField = (field) => {
  switch (field) {
    case 'title':
      errors.title = !form.title ? '请输入资源标题' : ''
      break
    case 'category_id':
      errors.category_id_id = !form.category_id ? '请选择资源分类' : ''
      break
    case 'fileType':
      errors.fileType = !form.fileType ? '请选择资源类型' : ''
      break
    case 'description':
      errors.description = !form.description ? '请输入资源描述' : ''
      break
    case 'file':
      errors.file = !form.file ? '请上传资源文件' : ''
      break
  }
  return !errors[field]
}

// 验证所有字段
const validateForm = () => {
  validateField('title')
  validateField('category_id')
  validateField('fileType')
  validateField('description')
  validateField('file')
  
  // 检查是否有错误
  return !Object.values(errors).some(error => error)
}

// 提交表单
const submitForm = async () => {
  if (!validateForm()) {
    toast.add({
      severity: 'error',
      summary: '表单验证失败',
      detail: '请检查表单填写是否正确',
      life: 3000
    })
    return
  }
  
  loading.value = true
  
  try {
    const formData = new FormData()
    formData.append('title', form.title)
    formData.append('description', form.description) 
    formData.append('category_id', form.category_id)
    formData.append('file', form.file)
    
    try {
      const response = await axios.post('/api/resources/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          'Authorization': `Bearer ${userStore.token}`
        },
        onUploadProgress: (progressEvent) => {
          const percentCompleted = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total
          )
          toast.add({
            severity: 'info',
            summary: '上传中',
            detail: `上传进度: ${percentCompleted}%`,
            life: 1000
          })
        }
      })
      
      if (response.status === 200) {
        toast.add({
          severity: 'success',
          summary: '上传成功',
          detail: '资源已成功上传',
          life: 3000
        })
        
        // 上传成功后跳转到资源列表页
        router.push('/resources')
      } else {
        throw new Error(response.data?.message || '资源上传失败')
      }
    } catch (error) {
      console.error('上传资源失败:', error)
      toast.add({
        severity: 'error',
        summary: '上传失败',
        detail: error.message || '资源上传失败，请稍后重试',
        life: 3000
      })
    } finally {
      loading.value = false
    }
  } catch (error) {
    console.error('上传资源失败:', error)
    toast.add({
      severity: 'error',
      summary: '上传失败',
      detail: error.message || '资源上传失败，请稍后重试',
      life: 3000
    })
    loading.value = false
  }
}
</script>