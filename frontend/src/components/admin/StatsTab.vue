<template>
  <div class="stats-tab">
    <div v-if="loading" class="loading">Loading stats...</div>
    <div v-else-if="error" class="error-message">{{ error }}</div>
    <div v-else class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon users">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats?.total_users || 0 }}</div>
          <div class="stat-label">Total Users</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon books">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
            <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
          </svg>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats?.total_books || 0 }}</div>
          <div class="stat-label">Total Books</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon downloads">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats?.total_downloads || 0 }}</div>
          <div class="stat-label">Total Downloads</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon active">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 12h-4l-3 9L9 3l-3 9H2"></path>
          </svg>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats?.active_users || 0 }}</div>
          <div class="stat-label">Active (24h)</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import apiClient from '@/api/client'

const loading = ref(false)
const error = ref('')
const stats = ref<any>(null)

onMounted(() => {
  loadStats()
})

async function loadStats() {
  loading.value = true
  error.value = ''
  try {
    const response = await apiClient.get('/admin/stats')
    stats.value = response.data
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to load stats'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-top: 2rem;
}

.stat-card {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.8) 0%, rgba(30, 41, 59, 0.6) 100%);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 16px;
  padding: 1.5rem;
  display: flex;
  gap: 1.5rem;
  align-items: center;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 12px 24px rgba(59, 130, 246, 0.15);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon svg {
  width: 28px;
  height: 28px;
  color: white;
}

.stat-icon.users { background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%); }
.stat-icon.books { background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%); }
.stat-icon.downloads { background: linear-gradient(135deg, #10b981 0%, #059669 100%); }
.stat-icon.active { background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%); }

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #e2e8f0;
  margin-bottom: 0.25rem;
}

.stat-label {
  color: #94a3b8;
  font-size: 0.875rem;
}

.loading, .error-message {
  padding: 2rem;
  text-align: center;
  color: #94a3b8;
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.8) 0%, rgba(30, 41, 59, 0.6) 100%);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
}

.error-message {
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.2);
}
</style>
