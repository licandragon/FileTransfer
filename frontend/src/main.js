import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { initAuth } from '@/composables/useAuth'
import App from './App.vue'
import router from './router'
import './styles.css'

await initAuth()
const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
