<template>
  <div class="admin-view">
    <main class="content">
      <div class="header">
        <div>
          <h1 class="page-title">Admin Panel</h1>
          <p class="page-subtitle">Manage users, invites, and system settings</p>
        </div>
        <div class="tabs">
          <button 
            :class="['tab', { active: activeTab === 'stats' }]"
            @click="activeTab = 'stats'"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="20" x2="18" y2="10"></line>
              <line x1="12" y1="20" x2="12" y2="4"></line>
              <line x1="6" y1="20" x2="6" y2="14"></line>
            </svg>
            Stats
          </button>
          <button 
            :class="['tab', { active: activeTab === 'users' }]"
            @click="activeTab = 'users'"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
            Users
          </button>
          <button 
            :class="['tab', { active: activeTab === 'invites' }]"
            @click="activeTab = 'invites'"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
            Invites
          </button>
        </div>
      </div>

      <div class="tab-content">
        <StatsTab v-if="activeTab === 'stats'" />
        
        <UsersTab
          v-if="activeTab === 'users'"
          :users="users"
          :loading="loadingUsers"
          :error="usersError"
          :pagination="usersPagination"
          @create-user="showCreateUserModal = true"
          @edit-user="editUser"
          @delete-user="deleteUser"
          @load-users="loadUsers"
        />
        
        <InvitesTab
          v-if="activeTab === 'invites'"
          :invites="invites"
          :invite-data="inviteData"
          :loading="loadingInvites"
          :error="invitesError"
          @create-invite="createInvite"
          @revoke-invite="revokeInvite"
        />
      </div>
    </main>

    <UserModal
      v-if="showCreateUserModal"
      @close="showCreateUserModal = false"
      @submit="handleCreateUser"
    />

    <UserModal
      v-if="showEditUserModal"
      :user="editingUser"
      :is-edit="true"
      @close="showEditUserModal = false"
      @submit="handleUpdateUser"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import apiClient from '@/api/client'
import StatsTab from '@/components/admin/StatsTab.vue'
import UsersTab from '@/components/admin/UsersTab.vue'
import InvitesTab from '@/components/admin/InvitesTab.vue'
import UserModal from '@/components/admin/UserModal.vue'

const activeTab = ref('stats')

// Users
const users = ref<any[]>([])
const loadingUsers = ref(false)
const usersError = ref('')
const usersPagination = ref<any>(null)
const showCreateUserModal = ref(false)
const showEditUserModal = ref(false)
const editingUser = ref<any>(null)

// Invites
const invites = ref<any[]>([])
const inviteData = ref<any>(null)
const loadingInvites = ref(false)
const invitesError = ref('')

onMounted(() => {
  loadUsers()
  loadInvites()
})

async function loadUsers(offset = 0) {
  loadingUsers.value = true
  usersError.value = ''
  try {
    const response = await apiClient.get(`/admin/users?limit=20&offset=${offset}`)
    users.value = response.data.users || []
    usersPagination.value = response.data.pagination
  } catch (err: any) {
    usersError.value = err.response?.data?.message || 'Failed to load users'
  } finally {
    loadingUsers.value = false
  }
}

async function loadInvites() {
  loadingInvites.value = true
  invitesError.value = ''
  try {
    const response = await apiClient.get('/invites')
    inviteData.value = response.data
    invites.value = response.data.invites || []
  } catch (err: any) {
    invitesError.value = err.response?.data?.message || 'Failed to load invites'
  } finally {
    loadingInvites.value = false
  }
}

async function handleCreateUser(data: any) {
  try {
    await apiClient.post('/admin/users', data)
    showCreateUserModal.value = false
    loadUsers()
  } catch (err: any) {
    alert(err.response?.data?.message || 'Failed to create user')
  }
}

function editUser(user: any) {
  editingUser.value = { ...user }
  showEditUserModal.value = true
}

async function handleUpdateUser(data: any) {
  try {
    await apiClient.put(`/admin/users/${editingUser.value.id}`, {
      username: data.username,
      role: data.role
    })
    showEditUserModal.value = false
    editingUser.value = null
    loadUsers()
  } catch (err: any) {
    alert(err.response?.data?.message || 'Failed to update user')
  }
}

async function deleteUser(user: any) {
  if (!confirm(`Are you sure you want to delete user "${user.username}"?`)) return
  
  try {
    await apiClient.delete(`/admin/users/${user.id}`)
    loadUsers()
  } catch (err: any) {
    alert(err.response?.data?.message || 'Failed to delete user')
  }
}

async function createInvite() {
  try {
    const response = await apiClient.post('/invites')
    loadInvites()
    // Invite will be shown in the list with copy button
  } catch (err: any) {
    alert(err.response?.data?.message || 'Failed to create invite')
  }
}

async function revokeInvite(token: string) {
  if (!confirm('Are you sure you want to revoke this invite?')) return
  
  try {
    await apiClient.post(`/invites/${token}/revoke`)
    loadInvites()
  } catch (err: any) {
    alert(err.response?.data?.message || 'Failed to revoke invite')
  }
}
</script>

<style scoped>
.admin-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0f1e 0%, #1a1f3a 100%);
}

.content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #fff 0%, #94a3b8 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.page-subtitle {
  font-size: 1.125rem;
  color: #64748b;
}

.tabs {
  display: flex;
  gap: 0.5rem;
  background: rgba(15, 23, 42, 0.6);
  padding: 0.5rem;
  border-radius: 12px;
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.tab {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: transparent;
  color: #94a3b8;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab svg {
  width: 18px;
  height: 18px;
}

.tab:hover {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.tab.active {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
}

.tab-content {
  animation: fadeIn 0.3s ease;
}

/* Stats Grid */
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

/* Section Header */
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

.secondary-button {
  padding: 0.75rem 1.5rem;
  background: rgba(100, 116, 139, 0.1);
  color: #94a3b8;
  border: 1px solid rgba(100, 116, 139, 0.3);
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.secondary-button:hover {
  background: rgba(100, 116, 139, 0.2);
}

/* Users Table */
.users-table {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.8) 0%, rgba(30, 41, 59, 0.6) 100%);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 16px;
  overflow: hidden;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead tr {
  background: rgba(59, 130, 246, 0.1);
}

th {
  padding: 1rem;
  text-align: left;
  color: #e2e8f0;
  font-weight: 600;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

td {
  padding: 1rem;
  color: #94a3b8;
  border-top: 1px solid rgba(59, 130, 246, 0.1);
}

tbody tr:hover {
  background: rgba(59, 130, 246, 0.05);
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 0.875rem;
}

.role-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.role-badge.admin {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.role-badge.user {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
}

.icon-button {
  width: 32px;
  height: 32px;
  padding: 0;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.icon-button svg {
  width: 16px;
  height: 16px;
  color: #3b82f6;
}

.icon-button:hover {
  background: rgba(59, 130, 246, 0.2);
  transform: translateY(-2px);
}

.icon-button.danger {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
}

.icon-button.danger svg {
  color: #ef4444;
}

.icon-button.danger:hover {
  background: rgba(239, 68, 68, 0.2);
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: rgba(59, 130, 246, 0.05);
  color: #94a3b8;
}

.pagination button {
  padding: 0.5rem 1rem;
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s ease;
}

.pagination button:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.2);
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Invites */
.invite-stats {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(37, 99, 235, 0.05) 100%);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
  padding: 1rem 1.5rem;
  margin-bottom: 2rem;
  color: #e2e8f0;
}

.invite-stats strong {
  color: #3b82f6;
  font-size: 1.25rem;
}

.invites-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
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

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.modal {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 16px;
  padding: 2rem;
  max-width: 500px;
  width: 90%;
  animation: slideUp 0.3s ease;
}

.modal h2 {
  color: #e2e8f0;
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  color: #94a3b8;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 8px;
  color: #e2e8f0;
  font-size: 1rem;
  transition: all 0.2s ease;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

/* Utility */
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

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
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
  .header {
    flex-direction: column;
    gap: 1rem;
  }

  .tabs {
    width: 100%;
  }

  .tab {
    flex: 1;
    justify-content: center;
    padding: 0.75rem 1rem;
  }

  .tab span {
    display: none;
  }

  .users-table {
    overflow-x: auto;
  }

  table {
    min-width: 600px;
  }
}
</style>
