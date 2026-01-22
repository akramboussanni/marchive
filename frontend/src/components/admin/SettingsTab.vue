<template>
  <div class="settings-tab">
    <div class="tab-header">
      <h2>Application Settings</h2>
      <p>Configure system-wide settings</p>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading settings...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <p>{{ error }}</p>
      <button @click="loadSettings" class="btn-primary">Retry</button>
    </div>

    <div v-else class="settings-list">
      <div class="setting-card">
        <div class="setting-info">
          <h3>Anonymous Access</h3>
          <p>Allow unauthenticated users to search for new books and request downloads. When disabled, anonymous users can only browse the existing library.</p>
        </div>
        <div class="setting-control">
          <label class="toggle-switch">
            <input 
              type="checkbox" 
              :checked="settings.anonymous_access_enabled === 'true'" 
              @change="toggleAnonymousAccess"
              :disabled="saving"
            />
            <span class="slider"></span>
          </label>
          <span class="toggle-label">{{ settings.anonymous_access_enabled === 'true' ? 'Enabled' : 'Disabled' }}</span>
        </div>
      </div>
    </div>

    <div v-if="saveMessage" :class="['save-message', saveMessageType]">
      {{ saveMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import apiClient from '@/api/client'

const loading = ref(true)
const error = ref('')
const saving = ref(false)
const saveMessage = ref('')
const saveMessageType = ref<'success' | 'error'>('success')
const settings = ref<Record<string, string>>({
  anonymous_access_enabled: 'false'
})

const loadSettings = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await apiClient.get('/admin/settings')
    settings.value = response.data.settings || {}
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to load settings'
  } finally {
    loading.value = false
  }
}

const toggleAnonymousAccess = async () => {
  const newValue = settings.value.anonymous_access_enabled === 'true' ? 'false' : 'true'
  saving.value = true
  saveMessage.value = ''
  
  try {
    await apiClient.post('/admin/settings', {
      key: 'anonymous_access_enabled',
      value: newValue
    })
    settings.value.anonymous_access_enabled = newValue
    saveMessage.value = 'Setting saved successfully'
    saveMessageType.value = 'success'
  } catch (err: any) {
    saveMessage.value = err.response?.data?.message || 'Failed to save setting'
    saveMessageType.value = 'error'
  } finally {
    saving.value = false
    setTimeout(() => { saveMessage.value = '' }, 3000)
  }
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.settings-tab {
  padding: 0;
}

.tab-header {
  margin-bottom: 2rem;
}

.tab-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.tab-header p {
  color: var(--text-secondary);
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  gap: 1rem;
  color: var(--text-secondary);
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--surface-alt);
  border-top-color: var(--accent-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.settings-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.setting-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  background: var(--surface-alt);
  border-radius: 12px;
  gap: 2rem;
}

.setting-info h3 {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.setting-info p {
  color: var(--text-secondary);
  font-size: 0.9rem;
  line-height: 1.5;
  max-width: 500px;
}

.setting-control {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.toggle-label {
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--text-secondary);
  min-width: 70px;
}

/* Toggle Switch */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 52px;
  height: 28px;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--surface-overlay);
  transition: 0.3s;
  border-radius: 28px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 22px;
  width: 22px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: var(--accent-primary);
}

input:checked + .slider:before {
  transform: translateX(24px);
}

input:disabled + .slider {
  opacity: 0.5;
  cursor: not-allowed;
}

.save-message {
  margin-top: 1rem;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  font-size: 0.9rem;
  text-align: center;
}

.save-message.success {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.save-message.error {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.btn-primary {
  padding: 0.75rem 1.5rem;
  background: var(--accent-primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn-primary:hover {
  opacity: 0.9;
}
</style>
