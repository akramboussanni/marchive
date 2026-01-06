<template>
  <div class="book-card">
    <div class="book-cover" @click="$emit('view', book)">
      <img 
        v-if="book.cover_url || book.cover_data" 
        :src="book.cover_data || book.cover_url" 
        :alt="book.title"
        @error="handleImageError"
      />
      <div v-else class="no-cover">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
          <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
        </svg>
      </div>
      <div v-if="book.is_ghost" class="ghost-badge" title="Ghost mode - only visible to requester and admins">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 2a7 7 0 0 0-7 7v3l-2 2v5h18v-5l-2-2V9a7 7 0 0 0-7-7z"></path>
          <path d="M8.5 19a4 4 0 0 0 7 0"></path>
          <circle cx="12" cy="11" r="1" fill="currentColor"></circle>
        </svg>
      </div>
      <div v-if="book.status === 'processing'" class="status-badge processing">
        <svg class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"></circle>
        </svg>
        <span>Processing...</span>
      </div>
      <div v-else-if="book.status === 'error'" class="status-badge failed">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
        <span>Failed</span>
      </div>
    </div>

    <div class="book-info">
      <h3 class="book-title" :title="book.title">{{ book.title }}</h3>
      <p class="book-author" :title="book.authors">{{ book.authors || 'Unknown Author' }}</p>
      <div class="book-meta">
        <span class="book-format">{{ book.format }}</span>
        <span class="book-size">{{ book.size }}</span>
      </div>
    </div>

    <div class="book-actions">
      <template v-if="actionMode === 'add'">
        <button 
          class="action-btn add-btn full-width"
          @click="$emit('addToLibrary', book)"
          :disabled="isDownloading || isAvailable"
          :title="isAvailable ? 'Already in library' : 'Add to Marchive'"
        >
          <svg v-if="!isDownloading && !isAvailable" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
            <polyline points="7 10 12 15 17 10"></polyline>
            <line x1="12" y1="15" x2="12" y2="3"></line>
          </svg>
          <svg v-else-if="isDownloading" class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <span v-if="!isDownloading && !isAvailable">Add to Marchive</span>
          <span v-else-if="isDownloading">Adding...</span>
          <span v-else>Already Added</span>
        </button>
      </template>
      
      <template v-else>
        <button 
          class="action-btn full-width open-btn"
          @click="$emit('open', book)"
          :disabled="book.status !== 'ready'"
          title="Open book"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
            <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
          </svg>
          <span>Open</span>
        </button>

        <div class="action-row">
          <button 
            class="action-btn primary-btn download-btn"
            @click="$emit('download', book)"
            :disabled="isDownloading || book.status !== 'ready'"
            :title="book.status === 'processing' ? 'Book is being downloaded to server...' : book.status === 'error' ? 'Download failed' : 'Download book'"
          >
            <svg v-if="!isDownloading && book.status !== 'processing'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
              <polyline points="7 10 12 15 17 10"></polyline>
              <line x1="12" y1="15" x2="12" y2="3"></line>
            </svg>
            <svg v-else class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
            </svg>
            <span>Download</span>
          </button>

          <div class="action-menu">
            <button 
              class="action-btn menu-btn"
              @click="showMenu = !showMenu"
              title="More options"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="1" fill="currentColor"></circle>
                <circle cx="19" cy="12" r="1" fill="currentColor"></circle>
                <circle cx="5" cy="12" r="1" fill="currentColor"></circle>
              </svg>
            </button>

            <div v-if="showMenu" class="menu-dropdown" @click.stop>
              <button 
                class="menu-item"
                @click="copyUrl(); showMenu = false"
                :title="copied ? 'Copied!' : 'Copy download link'"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                </svg>
                <span>{{ copied ? 'Copied!' : 'Copy Link' }}</span>
              </button>

              <button 
                class="menu-item"
                :class="{ favorited: isFavorited }"
                @click="$emit('toggleFavorite', book); showMenu = false"
                :title="isFavorited ? 'Remove from favorites' : 'Add to favorites'"
              >
                <svg viewBox="0 0 24 24" :fill="isFavorited ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                  <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path>
                </svg>
                <span>{{ isFavorited ? 'Unfavorite' : 'Favorite' }}</span>
              </button>

              <template v-if="showAdminControls">
                <div class="menu-divider"></div>
                
                <button 
                  class="menu-item"
                  @click="$emit('edit', book); showMenu = false"
                  title="Edit metadata"
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                  </svg>
                  <span>Edit Data</span>
                </button>

                <button 
                  class="menu-item danger"
                  @click="$emit('delete', book); showMenu = false"
                  title="Delete book"
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                  <span>Delete</span>
                </button>
              </template>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import type { Book } from '@/types/book'

interface Props {
  book: Book
  isFavorited?: boolean
  isDownloading?: boolean
  actionMode?: 'add' | 'full'
  showAdminControls?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showAdminControls: false
})

defineEmits<{
  view: [book: Book]
  toggleFavorite: [book: Book]
  download: [book: Book]
  open: [book: Book]
  addToLibrary: [book: Book]
  edit: [book: Book]
  delete: [book: Book]
}>()

const authStore = useAuthStore()
const copied = ref(false)
const imageError = ref(false)
const showMenu = ref(false)

const isAvailable = computed(() => props.book.status === 'ready')

const handleImageError = () => {
  imageError.value = true
}

const copyUrl = async () => {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:9520'
  const url = `${apiUrl}/api/books/${props.book.hash}/download`
  try {
    await navigator.clipboard.writeText(url)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy URL:', err)
  }
}

const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.action-menu')) {
    showMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.book-card {
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.book-card:hover {
  transform: translateY(-4px);
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 8px 24px rgba(59, 130, 246, 0.15);
}

.book-cover {
  width: 100%;
  aspect-ratio: 2/3;
  overflow: hidden;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(147, 51, 234, 0.1) 100%);
  cursor: pointer;
  position: relative;
}

.book-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.book-card:hover .book-cover img {
  transform: scale(1.05);
}

.status-badge {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.375rem 0.625rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  backdrop-filter: blur(8px);
  z-index: 10;
}

.status-badge svg {
  width: 14px;
  height: 14px;
}

.status-badge.processing {
  background: rgba(59, 130, 246, 0.9);
  color: white;
  border: 1px solid rgba(59, 130, 246, 0.4);
}

.status-badge.available {
  background: rgba(16, 185, 129, 0.9);
  color: white;
  border: 1px solid rgba(16, 185, 129, 0.4);
}

.status-badge.failed {
  background: rgba(239, 68, 68, 0.9);
  color: white;
  border: 1px solid rgba(239, 68, 68, 0.4);
}

.status-badge .spinner {
  animation: spin 1s linear infinite;
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
  width: 64px;
  height: 64px;
}

.book-info {
  padding: 1rem;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.book-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.3;
  min-height: 2.6em;
}

.book-author {
  font-size: 0.8rem;
  color: #94a3b8;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.book-meta {
  display: flex;
  gap: 0.5rem;
  font-size: 0.75rem;
  color: #64748b;
  margin-top: auto;
}

.book-format,
.book-size {
  padding: 0.25rem 0.5rem;
  background: rgba(59, 130, 246, 0.1);
  border-radius: 4px;
  font-weight: 500;
}

.book-actions {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.action-row {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.action-btn {
  padding: 0.625rem 1rem;
  border: 1px solid rgba(59, 130, 246, 0.3);
  background: rgba(59, 130, 246, 0.1);
  color: #60a5fa;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.action-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.5);
  transform: translateY(-1px);
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn svg {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.primary-btn {
  flex: 1;
}

.menu-btn {
  width: 40px;
  height: 40px;
  padding: 0.5rem;
  flex-shrink: 0;
}

.action-menu {
  position: relative;
}

.menu-dropdown {
  position: absolute;
  right: 0;
  bottom: calc(100% + 0.5rem);
  background: rgba(15, 23, 42, 0.95);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 8px;
  min-width: 160px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(8px);
  z-index: 100;
  overflow: hidden;
}

.menu-item {
  width: 100%;
  padding: 0.75rem 1rem;
  border: none;
  background: transparent;
  color: #94a3b8;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  text-align: left;
}

.menu-item:hover {
  background: rgba(59, 130, 246, 0.1);
  color: #60a5fa;
}

.menu-item.danger:hover {
  background: rgba(239, 68, 68, 0.1);
  color: #f87171;
}

.menu-item.favorited {
  color: #f472b6;
}

.menu-item svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.menu-divider {
  height: 1px;
  background: rgba(59, 130, 246, 0.2);
  margin: 0.25rem 0;
}

.action-btn.full-width {
  flex: 1;
}

.add-btn {
  font-size: 0.9rem;
  font-weight: 500;
}

.add-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
}

.add-btn:disabled {
  opacity: 0.6;
}

.favorite-btn.favorited {
  color: #ec4899;
}

.favorite-btn.favorited:hover {
  color: #db2777;
}

.spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .book-title {
    font-size: 0.85rem;
  }

  .menu-dropdown {
    right: auto;
    left: 0;
  }

  .book-author {
    font-size: 0.75rem;
  }

  .book-meta {
    font-size: 0.7rem;
  }

  .action-btn {
    padding: 0.65rem;
  }

  .action-btn svg {
    width: 16px;
    height: 16px;
  }
}
</style>
