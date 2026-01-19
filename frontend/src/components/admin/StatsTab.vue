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

    <!-- Admin Tools Section -->
    <div class="tools-section">
      <h3 class="section-title">Admin Tools</h3>
      <div class="tools-grid">
        <div class="tool-card">
          <div class="tool-icon restore">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
              <path d="M3 3v5h5"></path>
            </svg>
          </div>
          <div class="tool-info">
            <h4>Restore Books</h4>
            <p>Scan downloads folder and add any books not in the database</p>
          </div>
          <button 
            @click="handleRestore" 
            :disabled="restoring"
            class="tool-button"
          >
            {{ restoring ? 'Restoring...' : 'Run Restore' }}
          </button>
        </div>
      </div>

      <!-- Restore Result -->
      <div v-if="restoreResult" class="restore-result" :class="{ error: restoreResult.errors?.length }">
        <p><strong>Restore Complete:</strong></p>
        <p>✓ Restored: {{ restoreResult.restored }} books</p>
        <p>⊘ Skipped: {{ restoreResult.skipped }} books (already in database)</p>
        <div v-if="restoreResult.errors?.length" class="restore-errors">
          <p>✗ Errors:</p>
          <ul>
            <li v-for="(err, i) in restoreResult.errors" :key="i">{{ err }}</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import apiClient from '@/api/client'
import { booksApi } from '@/api/books'

const loading = ref(false)
const error = ref('')
const stats = ref<any>(null)
const restoring = ref(false)
const restoreResult = ref<any>(null)

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

async function handleRestore() {
  if (!confirm('This will scan the downloads folders and add any books not in the database. Continue?')) return
  
  restoring.value = true
  restoreResult.value = null
  
  try {
    const result = await booksApi.restoreBooks()
    restoreResult.value = result
    // Refresh stats after restore
    loadStats()
  } catch (err: any) {
    restoreResult.value = {
      restored: 0,
      skipped: 0,
      errors: [err.response?.data?.message || 'Failed to restore books']
    }
  } finally {
    restoring.value = false
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

/* Tools Section */
.tools-section {
  margin-top: 3rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #e2e8f0;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid rgba(59, 130, 246, 0.2);
}

.tools-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.tool-card {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.8) 0%, rgba(30, 41, 59, 0.6) 100%);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 16px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.tool-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.tool-icon svg {
  width: 24px;
  height: 24px;
  color: white;
}

.tool-icon.restore {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

.tool-info h4 {
  color: #e2e8f0;
  font-size: 1rem;
  margin-bottom: 0.25rem;
}

.tool-info p {
  color: #94a3b8;
  font-size: 0.875rem;
}

.tool-button {
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: auto;
}

.tool-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(16, 185, 129, 0.3);
}

.tool-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.restore-result {
  margin-top: 1.5rem;
  padding: 1rem;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 12px;
  color: #86efac;
}

.restore-result.error {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.3);
  color: #fca5a5;
}

.restore-result p {
  margin-bottom: 0.5rem;
}

.restore-result strong {
  color: #e2e8f0;
}

.restore-errors ul {
  margin-top: 0.5rem;
  padding-left: 1.5rem;
}

.restore-errors li {
  margin-bottom: 0.25rem;
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

