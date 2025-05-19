import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/Register.vue'),
      meta: { requiresGuest: true }
    },
    // 资源中心
    {
      path: '/resources',
      name: 'resources',
      component: () => import('../views/resources/ResourceList.vue')
    },
    {
      path: '/resources/:id',
      name: 'resourceDetail',
      component: () => import('../views/resources/ResourceDetail.vue')
    },
    {
      path: '/resources/upload',
      name: 'resourceUpload',
      component: () => import('../views/resources/ResourceUpload.vue'),
      meta: { requiresAuth: true }
    },
    // 论坛
    {
      path: '/forum',
      name: 'forum',
      component: () => import('../views/forum/ForumList.vue')
    },
    {
      path: '/forum/categories/:id',
      name: 'forum-category',
      component: () => import('../views/forum/ForumList.vue')
    },
    {
      path: '/forum/topics',
      name: 'forum-topics',
      component: () => import('../views/forum/TopicList.vue')
    },
    {
      path: '/forum/topic/:id',
      name: 'topic-detail',
      component: () => import('../views/forum/TopicDetail.vue')
    },
    {
      path: '/forum/create-topic',
      name: 'create-topic',
      component: () => import('../views/forum/CreateTopic.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/forum/create',
      name: 'forum-create',
      component: () => import('../views/forum/CreateTopic.vue'),
      meta: { requiresAuth: true }
    },

    // 用户中心
    {
      path: '/user/profile',
      name: 'user-profile',
      component: () => import('../views/user/UserProfile.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/user/resources',
      name: 'user-resources',
      component: () => import('../views/user/UserResources.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/user/favorites',
      name: 'user-favorites',
      component: () => import('../views/user/Favorites.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/user/points',
      name: 'user-points',
      component: () => import('../views/user/UserPoints.vue'),
      meta: { requiresAuth: true }
    },
    // AI助理
    {
      path: '/chat',
      name: 'chat-list',
      component: () => import('../views/chat/ChatList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/chat/:id',
      name: 'chat-detail',
      component: () => import('../views/chat/ChatDetail.vue'),
      meta: { requiresAuth: true }
    },
    // 管理员路由
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/admin/AdminDashboard.vue'),
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    // 404页面
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('../views/NotFound.vue')
    }
  ]
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 检查是否需要登录权限
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }
  
  // 检查是否需要管理员权限
  if (to.meta.requiresAdmin && userStore.user?.role !== 'admin') {
    next({ name: 'home' })
    return
  }
  
  // 已登录用户不能访问登录/注册页面
  if (to.meta.requiresGuest && userStore.isLoggedIn) {
    next({ name: 'home' })
    return
  }
  
  next()
})

export default router