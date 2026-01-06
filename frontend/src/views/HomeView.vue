<template>
  <div class="home-view">
    <main class="content">
      <div class="welcome-section">
        <h1 class="welcome-title">Library</h1>
        <p class="welcome-subtitle">Browse and discover books</p>
      </div>

      <section class="favorites-section">
        <div class="section-header">
          <h2>Favorites</h2>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
          </svg>
        </div>
        <div class="books-grid">
          <p class="empty-state">No favorites yet</p>
        </div>
      </section>

      <section class="library-section">
        <div class="section-header">
          <h2>All Books</h2>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
            <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
          </svg>
        </div>
        <div class="books-grid">
          <p class="empty-state">No books in library</p>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>

<style scoped>
.home-view {
  min-height: calc(100vh - 73px);
}

.content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.welcome-section {
  margin-bottom: 3rem;
  animation: fadeInUp 0.6s ease;
}

.welcome-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #fff 0%, #94a3b8 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: -0.02em;
}

.welcome-subtitle {
  font-size: 1.125rem;
  color: #64748b;
}

.favorites-section,
.library-section {
  margin-bottom: 3rem;
  animation: fadeInUp 0.6s ease 0.2s both;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid rgba(59, 130, 246, 0.2);
}

.section-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #e2e8f0;
}

.section-header svg {
  width: 24px;
  height: 24px;
  color: #3b82f6;
}

.books-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1.5rem;
}

.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 3rem;
  color: #64748b;
  font-size: 1rem;
  background: rgba(15, 23, 42, 0.4);
  border: 1px dashed rgba(59, 130, 246, 0.2);
  border-radius: 12px;
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

@media (max-width: 1024px) {
  .content {
    padding: 2rem 1.25rem;
  }
  
  .books-grid {
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 1.25rem;
  }
}

@media (max-width: 768px) {
  .content {
    padding: 1.5rem 1rem;
  }

  .welcome-title {
    font-size: 2rem;
  }

  .welcome-subtitle {
    font-size: 1rem;
  }

  .section-header h2 {
    font-size: 1.25rem;
  }

  .books-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 1rem;
  }
}

@media (max-width: 480px) {
  .content {
    padding: 1rem 0.75rem;
  }

  .welcome-section {
    margin-bottom: 2rem;
  }

  .welcome-title {
    font-size: 1.75rem;
  }

  .section-header {
    margin-bottom: 1rem;
  }

  .books-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 0.875rem;
  }

  .empty-state {
    padding: 2rem 1rem;
    font-size: 0.9rem;
  }
}
</style>
