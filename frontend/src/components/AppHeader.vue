<template>
  <header class="app-header">
    <div class="header-container">
      <router-link to="/" class="logo">
        <span class="logo-text">Marchive</span>
      </router-link>

      <div class="search-container">
        <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"></circle>
          <path d="m21 21-4.35-4.35"></path>
        </svg>
        <input
          v-if="authStore.isAuthenticated"
          type="text"
          placeholder="Search books..."
          class="search-input"
          v-model="searchQuery"
        />
        <router-link v-else to="/login" class="search-input-placeholder">
          <span>Search books (login required)</span>
        </router-link>
      </div>

      <nav class="nav-links">
        <template v-if="authStore.isAuthenticated">
          <router-link to="/" class="nav-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
              <polyline points="9 22 9 12 15 12 15 22"></polyline>
            </svg>
            <span>Home</span>
          </router-link>

          <router-link to="/profile" class="nav-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
            <span>Profile</span>
          </router-link>

          <router-link v-if="authStore.isAdmin" to="/admin" class="nav-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="9" y1="9" x2="15" y2="9"></line>
              <line x1="9" y1="15" x2="15" y2="15"></line>
            </svg>
            <span>Admin</span>
          </router-link>

          <button @click="handleLogout" class="nav-link logout-btn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
            <span>Logout</span>
          </button>
        </template>

        <template v-else>
          <router-link to="/" class="nav-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
              <polyline points="9 22 9 12 15 12 15 22"></polyline>
            </svg>
            <span>Home</span>
          </router-link>

          <router-link to="/login" class="nav-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path>
              <polyline points="10 17 15 12 10 7"></polyline>
              <line x1="15" y1="12" x2="3" y2="12"></line>
            </svg>
            <span>Login</span>
          </router-link>
        </template>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const searchQuery = ref('')

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-header {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  border-bottom: 1px solid rgba(59, 130, 246, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
  backdrop-filter: blur(12px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3);
}

.header-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  gap: 2rem;
}

.logo {
  text-decoration: none;
  transition: transform 0.2s ease;
  flex-shrink: 0;
}

.logo:hover {
  transform: scale(1.05);
}

.logo-text {
  font-size: 1.5rem;
  font-weight: 700;
  background: linear-gradient(135deg, #3b82f6 0%, #60a5fa 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: -0.02em;
}

.search-container {
  flex: 1;
  max-width: 600px;
  position: relative;
  margin: 0 auto;
}

.search-icon {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  width: 20px;
  height: 20px;
  color: #64748b;
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 3rem;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
  color: #e2e8f0;
  font-size: 0.95rem;
  transition: all 0.3s ease;
  outline: none;
}

.search-input::placeholder {
  color: #64748b;
}

.search-input:focus {
  background: rgba(15, 23, 42, 0.8);
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.search-input-placeholder {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 3rem;
  background: rgba(15, 23, 42, 0.3);
  border: 1px solid rgba(59, 130, 246, 0.1);
  border-radius: 12px;
  color: #64748b;
  font-size: 0.95rem;
  text-decoration: none;
  display: flex;
  align-items: center;
  transition: all 0.3s ease;
  cursor: pointer;
}

.search-input-placeholder:hover {
  background: rgba(15, 23, 42, 0.5);
  border-color: rgba(59, 130, 246, 0.2);
}

.nav-links {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  margin-left: auto;
  flex-shrink: 0;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  color: #94a3b8;
  text-decoration: none;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s ease;
  border: 1px solid transparent;
  cursor: pointer;
  background: transparent;
}

.nav-link svg {
  width: 18px;
  height: 18px;
}

.nav-link:hover {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.2);
  transform: translateY(-1px);
}

.nav-link.router-link-active {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.15);
  border-color: rgba(59, 130, 246, 0.3);
}

.logout-btn {
  color: #f87171;
}

.logout-btn:hover {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
}

@media (max-width: 768px) {
  .header-container {
    flex-wrap: wrap;
    gap: 1rem;
    padding: 0.875rem 1rem;
  }

  .logo-text {
    font-size: 1.25rem;
  }

  .search-container {
    order: 3;
    flex-basis: 100%;
  }

  .search-input,
  .search-input-placeholder {
    font-size: 0.9rem;
    padding: 0.625rem 1rem 0.625rem 2.75rem;
  }

  .search-icon {
    width: 18px;
    height: 18px;
    left: 0.875rem;
  }

  .nav-links {
    gap: 0.25rem;
  }

  .nav-link span {
    display: none;
  }

  .nav-link {
    padding: 0.625rem;
  }

  .nav-link svg {
    width: 20px;
    height: 20px;
  }
}

@media (max-width: 480px) {
  .header-container {
    padding: 0.75rem 0.875rem;
  }

  .logo-text {
    font-size: 1.125rem;
  }

  .search-input,
  .search-input-placeholder {
    font-size: 0.875rem;
    padding: 0.5rem 0.875rem 0.5rem 2.5rem;
  }

  .search-icon {
    width: 16px;
    height: 16px;
    left: 0.75rem;
  }

  .nav-link {
    padding: 0.5rem;
  }

  .nav-link svg {
    width: 18px;
    height: 18px;
  }
}
</style>
