<template>
  <div class="max-w-md mx-auto bg-white rounded-lg shadow-md overflow-hidden mt-10">
    <div class="py-4 px-6 bg-blue-600 text-white text-center">
      <h2 class="text-2xl font-bold">用户登录</h2>
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
      
      <div class="mb-6">
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
      
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center">
          <input type="checkbox" id="remember" class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded">
          <label for="remember" class="ml-2 block text-gray-700">记住我</label>
        </div>
        <div>
          <a href="#" class="text-blue-600 hover:underline">忘记密码?</a>
        </div>
      </div>
      
      <div>
        <button 
          type="submit" 
          class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
          :disabled="isSubmitting"
        >
          {{ isSubmitting ? '登录中...' : '登录' }}
        </button>
      </div>
      
      <div class="mt-4 text-center">
        <p class="text-gray-600">
          还没有账号? 
          <router-link to="/register" class="text-blue-600 hover:underline">立即注册</router-link>
        </p>
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useVuelidate } from '@vuelidate/core';
import { required, minLength } from '@vuelidate/validators';
import { useUserStore } from '../stores/user';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const formData = reactive({
  username: '',
  password: ''
});

const rules = {
  username: { required },
  password: { required, minLength: minLength(6) }
};

const v$ = useVuelidate(rules, formData);

const isSubmitting = ref(false);
const errorMessage = ref('');

const submitForm = async () => {
  errorMessage.value = '';
  
  // 表单验证
  const isFormValid = await v$.value.$validate();
  if (!isFormValid) return;
  
  isSubmitting.value = true;
  
  try {
    await userStore.login(formData);
    
    // 登录成功后重定向
    const redirectPath = route.query.redirect || '/';
    router.push(redirectPath);
  } catch (error) {
    if (error.response && error.response.data) {
      errorMessage.value = error.response.data.error || '登录失败，请检查用户名和密码';
    } else {
      errorMessage.value = '登录失败，请稍后再试';
    }
  } finally {
    isSubmitting.value = false;
  }
};
</script>