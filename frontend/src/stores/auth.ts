import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import type { User, LoginRequest, RegisterRequest } from '@/types/user'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const loading = ref(false)
  const initializing = ref(false)
  const initialized = ref(false)

  const isAuthenticated = computed(() => user.value !== null)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(credentials: LoginRequest) {
    loading.value = true
    try {
      await authApi.login(credentials)
      await fetchUser()
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    loading.value = true
    try {
      await authApi.logout()
      clearAuth()
    } finally {
      loading.value = false
    }
  }

  async function register(data: RegisterRequest) {
    loading.value = true
    try {
      await authApi.register(data)
      await fetchUser()
    } finally {
      loading.value = false
    }
  }

  async function fetchUser() {
    loading.value = true
    try {
      user.value = await authApi.getMe()
    } catch (error) {
      user.value = null
      throw error
    } finally {
      loading.value = false
    }
  }

  async function initialize() {
    if (initialized.value) return
    
    initialized.value = true
    initializing.value = true
    try {
      user.value = await authApi.getMe()
    } catch (error) {
      user.value = null
    } finally {
      initializing.value = false
    }
  }

  function clearAuth() {
    user.value = null
  }

  return {
    user,
    loading,
    initializing,
    initialized,
    isAuthenticated,
    isAdmin,
    login,
    logout,
    register,
    fetchUser,
    initialize,
    clearAuth
  }
})
