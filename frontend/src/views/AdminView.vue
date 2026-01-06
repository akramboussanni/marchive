<template>
  <div class="admin-view">
    <AppHeader />
    <main class="content">
      <h1 class="page-title">Admin Panel</h1>
      <p class="page-subtitle">Manage users, invites, and system settings</p>
      
      <div class="admin-grid">
        <div class="admin-card">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
          </div>
          <h3 class="card-title">User Management</h3>
          <p class="card-description">View, edit, and manage user accounts</p>
          <button class="card-button">Coming Soon</button>
        </div>

        <div class="admin-card">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
          </div>
          <h3 class="card-title">Invite System</h3>
          <p class="card-description">Create and manage invite codes</p>
          <button class="card-button">Coming Soon</button>
        </div>

        <div class="admin-card">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="3"></circle>
              <path d="M12 1v6m0 6v6m-9-9h6m6 0h6"></path>
            </svg>
          </div>
          <h3 class="card-title">System Stats</h3>
          <p class="card-description">View system analytics and statistics</p>
          <button class="card-button">Coming Soon</button>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
</script>

<style scoped>
.admin-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #0a0f1e 0%, #1a1f3a 100%);
}

.content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 3rem 2rem;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #fff 0%, #94a3b8 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: fadeInUp 0.6s ease;
}

.page-subtitle {
  font-size: 1.125rem;
  color: #64748b;
  margin-bottom: 3rem;
  animation: fadeInUp 0.6s ease 0.1s both;
}

.admin-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.admin-card {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.8) 0%, rgba(30, 41, 59, 0.6) 100%);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 16px;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  transition: all 0.3s ease;
  backdrop-filter: blur(12px);
  animation: fadeInUp 0.6s ease calc(0.2s + var(--delay, 0s)) both;
}

.admin-card:nth-child(1) { --delay: 0s; }
.admin-card:nth-child(2) { --delay: 0.1s; }
.admin-card:nth-child(3) { --delay: 0.2s; }

.admin-card:hover {
  transform: translateY(-4px);
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 12px 24px rgba(59, 130, 246, 0.15);
}

.card-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-icon svg {
  width: 28px;
  height: 28px;
  color: white;
}

.card-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #e2e8f0;
}

.card-description {
  color: #94a3b8;
  flex: 1;
}

.card-button {
  padding: 0.75rem 1.5rem;
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.card-button:hover {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.5);
  transform: translateY(-2px);
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
  .page-title {
    font-size: 1.875rem;
  }

  .admin-grid {
    grid-template-columns: 1fr;
  }
}
</style>
