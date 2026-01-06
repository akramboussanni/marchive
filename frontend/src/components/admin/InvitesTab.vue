<template>
  <div class="invites-tab">
    <div class="section-header">
      <h2 class="section-title">Invite Management</h2>
      <button class="primary-button" @click="$emit('create-invite')">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        Create Invite
      </button>
    </div>

    <div v-if="loading" class="loading">Loading invites...</div>
    <div v-else-if="error" class="error-message">{{ error }}</div>
    <div v-else class="invites-grid">
      <div v-for="invite in invites.filter(i => !i.revoked_at)" :key="invite.token" class="invite-card">
        <div class="invite-header">
          <span :class="['invite-status', { used: invite.used_at }]">
            {{ invite.used_at ? 'Used' : 'Active' }}
          </span>
          <span class="invite-date">{{ formatDate(invite.created_at) }}</span>
        </div>
        
        <div v-if="!invite.used_at && !invite.revoked_at" class="invite-link-section">
          <label class="invite-label">Invite URL</label>
          <div class="invite-url-container">
            <input 
              :value="getInviteUrl(invite.token)" 
              readonly 
              class="invite-url-input"
              :ref="(el) => { if (el) urlInputs[invite.token] = el as HTMLInputElement }"
            />
            <button 
              class="copy-button" 
              @click="copyToClipboard(invite.token)"
              :class="{ copied: copiedToken === invite.token }"
            >
              <svg v-if="copiedToken !== invite.token" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
              {{ copiedToken === invite.token ? 'Copied!' : 'Copy' }}
            </button>
          </div>
        </div>

        <div v-if="invite.invitee_username" class="invite-info">
          <div class="info-row">
            <span class="info-label">Used by:</span>
            <strong>{{ invite.invitee_username }}</strong>
          </div>
          <div v-if="invite.used_at" class="info-row">
            <span class="info-label">Used on:</span>
            <span>{{ formatDate(invite.used_at) }}</span>
          </div>
        </div>
        
        <button 
          v-if="!invite.used_at && !invite.revoked_at"
          class="revoke-button"
          @click="$emit('revoke-invite', invite.token)"
        >
          Revoke Invite
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  invites: any[]
  inviteData: any
  loading: boolean
  error: string
}>()

defineEmits<{
  'create-invite': []
  'revoke-invite': [token: string]
}>()

const urlInputs = ref<Record<string, HTMLInputElement | null>>({})
const copiedToken = ref<string | null>(null)

function getInviteUrl(token: string): string {
  const baseUrl = window.location.origin
  return `${baseUrl}/register?token=${token}`
}

function copyToClipboard(token: string) {
  const url = getInviteUrl(token)
  const input = urlInputs.value[token]
  
  if (input) {
    input.select()
    input.setSelectionRange(0, 99999) // For mobile devices
  }
  
  navigator.clipboard.writeText(url).then(() => {
    copiedToken.value = token
    setTimeout(() => {
      copiedToken.value = null
    }, 2000)
  })
}

function formatDate(timestamp?: string | number) {
  if (!timestamp) return 'N/A'
  const date = new Date(typeof timestamp === 'string' ? parseInt(timestamp) * 1000 : timestamp * 1000)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}
</script>

<style scoped>
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #e2e8f0;
}

.primary-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.primary-button svg {
  width: 18px;
  height: 18px;
}

.primary-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(59, 130, 246, 0.3);
}

.invites-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(500px, 1fr));
  gap: 1.5rem;
}

.invite-card {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.8) 0%, rgba(30, 41, 59, 0.6) 100%);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
  padding: 1.5rem;
  transition: all 0.3s ease;
}

.invite-card:hover {
  transform: translateY(-4px);
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 8px 16px rgba(59, 130, 246, 0.15);
}

.invite-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.invite-status {
  padding: 0.25rem 0.75rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.invite-status.used {
  background: rgba(100, 116, 139, 0.1);
  color: #94a3b8;
  border-color: rgba(100, 116, 139, 0.2);
}

.invite-status.revoked {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.2);
}

.invite-date {
  color: #64748b;
  font-size: 0.875rem;
}

.invite-link-section {
  margin-bottom: 1rem;
}

.invite-label {
  display: block;
  color: #94a3b8;
  font-size: 0.875rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.invite-url-container {
  display: flex;
  gap: 0.5rem;
}

.invite-url-input {
  flex: 1;
  font-family: 'Courier New', monospace;
  background: rgba(15, 23, 42, 0.8);
  padding: 0.75rem;
  border-radius: 8px;
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.2);
  font-size: 0.875rem;
  outline: none;
}

.invite-url-input:focus {
  border-color: rgba(59, 130, 246, 0.4);
}

.copy-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.copy-button svg {
  width: 16px;
  height: 16px;
}

.copy-button:hover {
  background: rgba(59, 130, 246, 0.2);
}

.copy-button.copied {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  border-color: rgba(16, 185, 129, 0.2);
}

.invite-token {
  font-family: 'Courier New', monospace;
  background: rgba(59, 130, 246, 0.1);
  padding: 0.75rem;
  border-radius: 8px;
  color: #3b82f6;
  word-break: break-all;
  margin-bottom: 1rem;
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.invite-info {
  color: #94a3b8;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
}

.info-label {
  color: #64748b;
}

.invite-info strong {
  color: #e2e8f0;
}

.revoke-button {
  width: 100%;
  padding: 0.5rem;
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.revoke-button:hover {
  background: rgba(239, 68, 68, 0.2);
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
