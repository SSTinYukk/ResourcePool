<template>
  <div class="max-w-md mx-auto bg-white rounded-lg shadow-md overflow-hidden mt-10">
    <div class="py-4 px-6 bg-blue-600 text-white text-center">
      <h2 class="text-2xl font-bold">用户注册</h2>
    </div>
    
    <form @submit.prevent="submitForm" class="py-6 px-8">
      <div class="mb-4" v-if="errorMessage">
        <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {{ errorMessage }}
        </div>
      </div>
      
      <div class="mb-4">
        <label for="username" class="block text-gray-700 font-medium mb-2">用户名</label>
        <input 
          type="text" 
          id="username" 
          v-model="formData.username" 
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          :class="{ 'border-red-500': v$.username.$error }"
        >
        <div v-if="v$.username.$error" class="text-red-500 text-sm mt-1">
          {{ v$.username.$errors[0].$message }}
        </div>
      </div>
      
      <div class="mb-4">
        <label for="email" class="block text-gray-700 font-medium mb-2">电子邮箱</label>
        <input 
          type="email" 
          id="email" 
          v-model="formData.email" 
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          :class="{ 'border-red-500': v$.email.$error }"
        >
        <div v-if="v$.email.$error" class="text-red-500 text-sm mt-1">
          {{ v$.email.$errors[0].$message }}
        </div>
      </div>
      
      <div class="mb-4">
        <label for="password" class="block text-gray-700 font-medium mb-2">密码</label>
        <input 
          type="password" 
          id="password" 
          v-model="formData.password" 
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          :class="{ 'border-red-500': v$.password.$error }"
        >
        <div v-if="v$.password.$error" class="text-red-500 text-sm mt-1">
          {{ v$.password.$errors[0].$message }}
        </div>
      </div>
      
      <div class="mb-6">
        <label for="confirmPassword" class="block text-gray-700 font-medium mb-2">确认密码</label>
        <input 
          type="password" 
          id="confirmPassword" 
          v-model="formData.confirmPassword" 
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          :class="{ 'border-red-500': v$.confirmPassword.$error }"
        >
        <div v-if="v$.confirmPassword.$error" class="text-red-500 text-sm mt-1">
          {{ v$.confirmPassword.$errors[0].$message }}
        </div>
      </div>
      
      <div>
        <button 
          type="submit" 
          class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
          :disabled="isSubmitting"
        >
          {{ isSubmitting ? '注册中...' : '注册' }}
        </button>
      </div>
      
      <div class="mt-4 text-center">
        <p class="text-gray-600">
          已有账号? 
          <router-link to="/login" class="text-blue-600 hover:underline">立即登录</router-link>
        </p>
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useVuelidate } from '@vuelidate/core';
import { required, email, minLength, sameAs } from '@vuelidate/validators';
import { useUserStore } from '../stores/user';

const router = useRouter();
const userStore = useUserStore();

const formData = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
});

const errorMessage = ref('');
const isSubmitting = ref(false);
// 修改验证规则
const rules = {
  username: { required },
  email: { required, email },
  password: { required, minLength: minLength(6) },
  confirmPassword: { required }
};

const v$ = useVuelidate(rules, formData);

// 修改表单提交逻辑
const submitForm = async () => {
  errorMessage.value = '';
  
  // 强制重新验证
  v$.value.$reset();
  const isFormValid = await v$.value.$validate();
  
  if (!isFormValid) {
    // 添加调试日志
    console.log('验证错误:', v$.value.$errors);
    return;
  }
  
  isSubmitting.value = true;
  
  try {
    // 准备注册数据
    const registerData = {
      username: formData.username,
      email: formData.email,
      password: formData.password
    };
    
    await userStore.register(registerData);
    
    // 注册成功后跳转到登录页
    router.push('/login');
  } catch (error) {
    if (error.response && error.response.data) {
      errorMessage.value = error.response.data.error || '注册失败，请稍后再试';
    } else {
      errorMessage.value = '注册失败，请稍后再试';
    }
  } finally {
    isSubmitting.value = false;
  }
};
</script>