<template>
  <div class="profile-view">
    <main class="content">
      <div class="profile-header">
        <div class="avatar">
          <span>{{ authStore.user?.username?.charAt(0).toUpperCase() }}</span>
        </div>
        <div class="header-info">
          <h1 class="profile-title">{{ authStore.user?.username }}</h1>
          <span class="role-badge">{{ authStore.user?.role }}</span>
        </div>
      </div>

      <div class="profile-card">
        <h2 class="section-title">Account Information</h2>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">Username</span>
            <span class="info-value">{{ authStore.user?.username }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Role</span>
            <span class="info-value">{{ authStore.user?.role }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Member Since</span>
            <span class="info-value">{{ formatDate(authStore.user?.created_at) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">User ID</span>
            <span class="info-value">{{ authStore.user?.id }}</span>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useMeta } from '@/composables/useMeta'

useMeta({
  title: 'Profile',
  description: 'Manage your mArchive profile and settings'
})

const authStore = useAuthStore()

function formatDate(timestamp?: string) {
  if (!timestamp) return ''
  const date = new Date(parseInt(timestamp) * 1000)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>

<style scoped>
.profile-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0f1e 0%, #1a1f3a 100%);
}

.content {
  max-width: 900px;
  margin: 0 auto;
  padding: 3rem 2rem;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 2rem;
  margin-bottom: 3rem;
  animation: fadeInUp 0.6s ease;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 20px;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.5rem;
  font-weight: 700;
  color: white;
  box-shadow: 0 8px 16px rgba(59, 130, 246, 0.3);
}

.header-info {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.profile-title {
  font-size: 2rem;
  font-weight: 700;
  color: #e2e8f0;
  margin: 0;
}

.role-badge {
  display: inline-block;
  padding: 0.375rem 1rem;
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  width: fit-content;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.profile-card {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.8) 0%, rgba(30, 41, 59, 0.6) 100%);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 16px;
  padding: 2rem;
  backdrop-filter: blur(12px);
  animation: fadeInUp 0.6s ease 0.2s both;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #e2e8f0;
  margin-bottom: 1.5rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-label {
  font-size: 0.875rem;
  color: #64748b;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.info-value {
  font-size: 1.125rem;
  color: #e2e8f0;
  font-weight: 500;
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

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .header-info {
    align-items: center;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
