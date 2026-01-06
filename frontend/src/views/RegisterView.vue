<template>
  <div class="register-view">
    <div class="register-card">
      <h1>Register</h1>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="token">Invite Code</label>
          <input
            id="token"
            v-model="token"
            type="text"
            required
            placeholder="Enter your invite code"
          />
        </div>
        <div class="form-group">
          <label for="username">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            required
            autocomplete="username"
          />
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            autocomplete="new-password"
          />
        </div>
        <button type="submit" :disabled="authStore.loading">
          {{ authStore.loading ? 'Registering...' : 'Register' }}
        </button>
        <p v-if="error" class="error">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const token = ref('')
const username = ref('')
const password = ref('')
const error = ref('')

onMounted(() => {
  if (route.query.token) {
    token.value = route.query.token as string
  } else {
    error.value = 'No invite code provided. Registration requires an invite.'
  }
})

async function handleRegister() {
  error.value = ''
  try {
    await authStore.register({
      token: token.value,
      username: username.value,
      password: password.value
    })
    router.push('/')
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Registration failed'
  }
}
</script>

<style scoped>
.register-view {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 73px);
  padding: 2rem;
}

.register-card {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.9) 0%, rgba(30, 41, 59, 0.8) 100%);
  backdrop-filter: blur(12px);
  padding: 3rem;
  border-radius: 16px;
  border: 1px solid rgba(59, 130, 246, 0.2);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5);
  width: 100%;
  max-width: 440px;
  animation: fadeInUp 0.6s ease;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

h1 {
  margin-bottom: 2rem;
  text-align: center;
  color: #e2e8f0;
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #fff 0%, #94a3b8 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.form-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  margin-bottom: 0.625rem;
  font-weight: 500;
  color: #94a3b8;
  font-size: 0.9rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

input {
  width: 100%;
  padding: 0.875rem 1rem;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 10px;
  font-size: 1rem;
  transition: all 0.3s ease;
  color: #e2e8f0;
}

input::placeholder {
  color: #64748b;
}

input:focus {
  outline: none;
  background: rgba(15, 23, 42, 0.8);
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

button {
  width: 100%;
  padding: 1rem;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  margin-top: 1rem;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(59, 130, 246, 0.4);
}

button:active:not(:disabled) {
  transform: translateY(0);
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error {
  color: #f87171;
  margin-top: 1rem;
  text-align: center;
  font-weight: 500;
  padding: 0.75rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 8px;
}

@media (max-width: 480px) {
  .register-view {
    padding: 1rem;
  }

  .register-card {
    padding: 2rem 1.5rem;
  }

  h1 {
    font-size: 1.75rem;
  }
}
</style>
