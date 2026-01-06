import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Initialize auth state on app start
const authStore = useAuthStore()
authStore.initialize().catch(() => {
  // Silently fail if user is not logged in
})

app.mount('#app')
