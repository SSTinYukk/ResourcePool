import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// PrimeVue组件和样式
import PrimeVue from 'primevue/config'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import 'primevue/resources/themes/lara-light-blue/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'

// 自定义样式
import './assets/css/tailwind.css'
import './assets/css/main.css'

// 全局UI组件
import AppButton from './components/ui/AppButton.vue'
import AppCard from './components/ui/AppCard.vue'
import AppInput from './components/ui/AppInput.vue'
import AppForm from './components/ui/AppForm.vue'
import AppSelect from './components/ui/AppSelect.vue'
import AppTable from './components/ui/AppTable.vue'
import AppMessage from './components/ui/AppMessage.vue'
import AppPagination from './components/ui/AppPagination.vue'
import AppTag from './components/ui/AppTag.vue'
import AppModal from './components/ui/AppModal.vue'
import AppSkeleton from './components/ui/AppSkeleton.vue'

const app = createApp(App)

// 注册Pinia和路由
app.use(createPinia())
app.use(router)

// 注册PrimeVue及其服务
app.use(PrimeVue, { ripple: true })
app.use(ToastService)
app.use(ConfirmationService)

// 注册全局UI组件
app.component('AppButton', AppButton)
app.component('AppCard', AppCard)
app.component('AppInput', AppInput)
app.component('AppForm', AppForm)
app.component('AppSelect', AppSelect)
app.component('AppTable', AppTable)
app.component('AppMessage', AppMessage)
app.component('AppPagination', AppPagination)
app.component('AppTag', AppTag)
app.component('AppModal', AppModal)
app.component('AppSkeleton', AppSkeleton)

app.mount('#app')