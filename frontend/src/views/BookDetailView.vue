<template>
  <div class="book-detail-view">
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading book details...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="15" y1="9" x2="9" y2="15"></line>
        <line x1="9" y1="9" x2="15" y2="15"></line>
      </svg>
      <h2>{{ error }}</h2>
      <button @click="$router.push('/')" class="back-btn">
        Go Back to Library
      </button>
    </div>

    <div v-else-if="book" class="book-detail-content">
      <button @click="$router.push('/')" class="back-link">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"></line>
          <polyline points="12 19 5 12 12 5"></polyline>
        </svg>
        Back to Library
      </button>

      <div class="book-detail-grid">
        <div class="book-cover-section">
          <div class="book-cover-large">
            <img 
              v-if="book.cover_url || book.cover_data" 
              :src="book.cover_data || book.cover_url" 
              :alt="book.title"
            />
            <div v-else class="no-cover">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
                <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
              </svg>
            </div>
          </div>

          <div class="book-actions">
            <button 
              v-if="book.status === 'ready'"
              @click="handleDownload"
              class="primary-action-btn"
              :disabled="downloading"
            >
              <svg v-if="!downloading" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="7 10 12 15 17 10"></polyline>
                <line x1="12" y1="15" x2="12" y2="3"></line>
              </svg>
              <svg v-else class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
              </svg>
              <span v-if="!downloading">Download Book</span>
              <span v-else>Downloading...</span>
            </button>

            <div v-else-if="book.status === 'processing'" class="status-info processing">
              <svg class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
              </svg>
              <span>Book is being processed...</span>
            </div>

            <div v-else-if="book.status === 'error'" class="status-info error">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="15" y1="9" x2="9" y2="15"></line>
                <line x1="9" y1="9" x2="15" y2="15"></line>
              </svg>
              <span>Download failed</span>
            </div>

            <button 
              @click="handleToggleFavorite"
              class="secondary-action-btn"
              :class="{ favorited: isFavorited }"
            >
              <svg viewBox="0 0 24 24" :fill="isFavorited ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
              </svg>
              <span>{{ isFavorited ? 'Remove from Favorites' : 'Add to Favorites' }}</span>
            </button>
          </div>
        </div>

        <div class="book-info-section">
          <div class="book-header">
            <div>
              <h1 class="book-title">{{ book.title }}</h1>
              <p class="book-authors">{{ book.authors || 'Unknown Author' }}</p>
            </div>
            <div v-if="book.is_ghost" class="ghost-badge-large" title="Ghost mode - only visible to requester and admins">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2a7 7 0 0 0-7 7v3l-2 2v5h18v-5l-2-2V9a7 7 0 0 0-7-7z"></path>
                <path d="M8.5 19a4 4 0 0 0 7 0"></path>
                <circle cx="12" cy="11" r="1" fill="currentColor"></circle>
              </svg>
              <span>Ghost Mode</span>
            </div>
          </div>

          <div class="book-metadata">
            <div class="metadata-item" v-if="book.publisher">
              <span class="metadata-label">Publisher</span>
              <span class="metadata-value">{{ book.publisher }}</span>
            </div>

            <div class="metadata-item" v-if="book.language">
              <span class="metadata-label">Language</span>
              <span class="metadata-value">{{ book.language }}</span>
            </div>

            <div class="metadata-item">
              <span class="metadata-label">Format</span>
              <span class="metadata-value">{{ book.format }}</span>
            </div>

            <div class="metadata-item">
              <span class="metadata-label">Size</span>
              <span class="metadata-value">{{ book.size }}</span>
            </div>

            <div class="metadata-item">
              <span class="metadata-label">Downloads</span>
              <span class="metadata-value">{{ book.download_count || 0 }}</span>
            </div>

            <div class="metadata-item" v-if="book.created_at">
              <span class="metadata-label">Added</span>
              <span class="metadata-value">{{ formatDate(book.created_at) }}</span>
            </div>

            <div class="metadata-item" v-if="requestedBy">
              <span class="metadata-label">Requested By</span>
              <span class="metadata-value">{{ requestedBy.username }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { booksApi } from '@/api/books'
import type { Book } from '@/types/book'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const book = ref<Book | null>(null)
const requestedBy = ref<{ id: string, username: string, role: string } | null>(null)
const loading = ref(false)
const error = ref('')
const downloading = ref(false)
const isFavorited = ref(false)

const loadBookDetail = async () => {
  const hash = route.params.hash as string
  if (!hash) {
    error.value = 'Book not found'
    return
  }

  try {
    loading.value = true
    const response = await booksApi.getBookDetail(hash)
    book.value = response.book
    requestedBy.value = response.requested_by || null

    // Check if favorited
    if (authStore.isAuthenticated) {
      try {
        const favoritesResponse = await booksApi.getFavorites(1000, 0)
        isFavorited.value = favoritesResponse.books.some(b => b.hash === hash)
      } catch (e) {
        console.error('Failed to load favorites:', e)
      }
    } else {
      const localFavorites = localStorage.getItem('marchive_favorites')
      if (localFavorites) {
        try {
          const hashes = JSON.parse(localFavorites) as string[]
          isFavorited.value = hashes.includes(hash)
        } catch (e) {
          console.error('Failed to parse local favorites:', e)
        }
      }
    }
  } catch (e: any) {
    console.error('Failed to load book detail:', e)
    error.value = e.response?.data?.message || 'Failed to load book details'
  } finally {
    loading.value = false
  }
}

const handleDownload = async () => {
  if (!book.value || book.value.status !== 'ready') return

  try {
    downloading.value = true
    const downloadUrl = `${import.meta.env.VITE_API_URL || 'http://localhost:9520'}/api/books/${book.value.hash}/download`
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = `${book.value.title}.${book.value.format}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (error) {
    console.error('Failed to download book:', error)
  } finally {
    downloading.value = false
  }
}

const handleToggleFavorite = async () => {
  if (!book.value) return

  if (!authStore.isAuthenticated) {
    // Handle local favorites
    const localFavorites = localStorage.getItem('marchive_favorites')
    let hashes: string[] = []
    
    if (localFavorites) {
      try {
        hashes = JSON.parse(localFavorites)
      } catch (error) {
        console.error('Failed to parse local favorites:', error)
      }
    }

    if (isFavorited.value) {
      hashes = hashes.filter(h => h !== book.value!.hash)
      isFavorited.value = false
    } else {
      hashes.push(book.value.hash)
      isFavorited.value = true
    }

    localStorage.setItem('marchive_favorites', JSON.stringify(hashes))
    return
  }

  try {
    const response = await booksApi.toggleFavorite(book.value.hash)
    isFavorited.value = response.is_favorited
  } catch (error) {
    console.error('Failed to toggle favorite:', error)
  }
}

const formatDate = (timestamp: string) => {
  const date = new Date(parseInt(timestamp) * 1000)
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
}

onMounted(() => {
  loadBookDetail()
})
</script>

<style scoped>
.book-detail-view {
  min-height: calc(100vh - 73px);
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  padding: 2rem 1.5rem;
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 4rem 2rem;
  color: #94a3b8;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid rgba(59, 130, 246, 0.1);
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.error-state svg {
  width: 64px;
  height: 64px;
  color: #ef4444;
}

.error-state h2 {
  font-size: 1.5rem;
  color: #e2e8f0;
  margin: 0;
}

.back-btn {
  padding: 0.75rem 2rem;
  background: rgba(59, 130, 246, 0.2);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 8px;
  color: #3b82f6;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: rgba(59, 130, 246, 0.3);
  border-color: rgba(59, 130, 246, 0.5);
}

.book-detail-content {
  max-width: 1200px;
  margin: 0 auto;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 8px;
  color: #3b82f6;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 2rem;
}

.back-link:hover {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.4);
}

.back-link svg {
  width: 20px;
  height: 20px;
}

.book-detail-grid {
  display: grid;
  grid-template-columns: 350px 1fr;
  gap: 3rem;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 16px;
  padding: 2rem;
}

.book-cover-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.book-cover-large {
  width: 100%;
  aspect-ratio: 2/3;
  border-radius: 12px;
  overflow: hidden;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(147, 51, 234, 0.1) 100%);
}

.book-cover-large img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-cover {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(148, 163, 184, 0.5);
}

.no-cover svg {
  width: 96px;
  height: 96px;
}

.book-actions {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.primary-action-btn,
.secondary-action-btn {
  padding: 1rem 1.5rem;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  border: none;
}

.primary-action-btn {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.primary-action-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.3);
  border-color: rgba(59, 130, 246, 0.5);
}

.primary-action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.secondary-action-btn {
  background: rgba(100, 116, 139, 0.1);
  color: #94a3b8;
  border: 1px solid rgba(100, 116, 139, 0.2);
}

.secondary-action-btn:hover {
  background: rgba(100, 116, 139, 0.2);
  color: #cbd5e1;
}

.secondary-action-btn.favorited {
  color: #ec4899;
  border-color: rgba(236, 72, 153, 0.3);
  background: rgba(236, 72, 153, 0.1);
}

.secondary-action-btn.favorited:hover {
  background: rgba(236, 72, 153, 0.2);
  color: #f472b6;
}

.primary-action-btn svg,
.secondary-action-btn svg {
  width: 20px;
  height: 20px;
}

.status-info {
  padding: 1rem 1.5rem;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  font-weight: 500;
}

.status-info.processing {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.status-info.error {
  background: rgba(239, 68, 68, 0.2);
  color: #fca5a5;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.status-info svg {
  width: 20px;
  height: 20px;
}

.book-info-section {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.book-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
}

.book-title {
  font-size: 2rem;
  font-weight: 700;
  color: #e2e8f0;
  margin: 0 0 0.5rem 0;
  line-height: 1.3;
}

.book-authors {
  font-size: 1.25rem;
  color: #94a3b8;
  margin: 0;
}

.ghost-badge-large {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  background: rgba(147, 51, 234, 0.2);
  border: 1px solid rgba(147, 51, 234, 0.3);
  border-radius: 8px;
  color: #c084fc;
  font-weight: 500;
}

.ghost-badge-large svg {
  width: 20px;
  height: 20px;
}

.book-metadata {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

.metadata-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.metadata-label {
  font-size: 0.875rem;
  color: #64748b;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.metadata-value {
  font-size: 1.125rem;
  color: #e2e8f0;
  font-weight: 500;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 1024px) {
  .book-detail-grid {
    grid-template-columns: 1fr;
    gap: 2rem;
  }

  .book-cover-section {
    max-width: 400px;
    margin: 0 auto;
  }
}

@media (max-width: 768px) {
  .book-detail-view {
    padding: 1rem;
  }

  .book-title {
    font-size: 1.5rem;
  }

  .book-authors {
    font-size: 1rem;
  }

  .book-metadata {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .book-detail-grid {
    padding: 1.5rem;
  }
}
</style>
