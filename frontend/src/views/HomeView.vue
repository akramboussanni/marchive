<template>
  <div class="home-view">
    <main class="content">
      <div class="welcome-section">
        <h1 class="welcome-title">Library</h1>
        <p class="welcome-subtitle">Browse and discover books</p>
      </div>

      <section v-if="favorites.length > 0" class="favorites-section">
        <div class="section-header">
          <h2>Favorites</h2>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
          </svg>
        </div>
        <div class="books-grid">
          <BookCard 
            v-for="book in favorites"
            :key="book.hash"
            :book="book"
            :is-favorited="true"
            :is-downloading="downloadingBooks.has(book.hash)"
            @toggle-favorite="handleToggleFavorite"
            @download="handleDownload"
            @open="handleOpen"
            @read="handleRead"
          />
        </div>
      </section>

      <section class="library-section">
        <div class="section-header">
          <h2>All Books</h2>
          <div class="header-actions">
            <button 
              class="refresh-btn" 
              @click="refreshLibrary" 
              :disabled="loadingBooks"
              title="Refresh library"
            >
              <svg :class="{ spinning: loadingBooks }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="23 4 23 10 17 10"></polyline>
                <polyline points="1 20 1 14 7 14"></polyline>
                <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
              </svg>
            </button>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path>
              <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path>
            </svg>
          </div>
        </div>
        <div v-if="loadingBooks" class="loading-state">
          <div class="spinner"></div>
          <p>Loading books...</p>
        </div>
        <div v-else-if="books.length > 0" class="books-grid">
          <BookCard 
            v-for="book in books"
            :key="book.hash"
            :book="book"
            :is-favorited="isBookFavorited(book.hash)"
            :is-downloading="downloadingBooks.has(book.hash)"
            :show-admin-controls="authStore.user?.role === 'admin'"
            @toggle-favorite="handleToggleFavorite"
            @download="handleDownload"
            @open="handleOpen"
            @read="handleRead"
            @view="handleViewBook"
            @edit="handleEditBook"
            @delete="handleDeleteBook"
          />
        </div>
        <div v-else class="books-grid">
          <p class="empty-state">No books in library</p>
        </div>
        
        <div v-if="booksPagination.has_next" class="load-more">
          <button @click="loadMoreBooks" :disabled="loadingMore" class="load-more-btn">
            <span v-if="!loadingMore">Load More</span>
            <div v-else class="spinner-small"></div>
          </button>
        </div>
      </section>
    </main>

    <!-- Dialogs -->
    <ConfirmDialog
      :is-open="deleteDialog.isOpen"
      :title="deleteDialog.title"
      :message="deleteDialog.message"
      confirm-text="Delete"
      cancel-text="Cancel"
      variant="danger"
      @confirm="confirmDelete"
      @cancel="cancelDelete"
    />

    <EditMetadataDialog
      :is-open="editDialog.isOpen"
      :title="editDialog.title"
      :authors="editDialog.authors"
      :publisher="editDialog.publisher"
      @save="saveMetadata"
      @cancel="cancelEdit"
    />

    <NotificationToast ref="toastRef" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { booksApi } from '@/api/books'
import BookCard from '@/components/BookCard.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import EditMetadataDialog from '@/components/EditMetadataDialog.vue'
import NotificationToast from '@/components/NotificationToast.vue'
import type { Book, Pagination } from '@/types/book'

const router = useRouter()
const authStore = useAuthStore()

const books = ref<Book[]>([])
const loadingBooks = ref(false)
const loadingMore = ref(false)
const downloadingBooks = ref(new Set<string>())
let refreshInterval: ReturnType<typeof setInterval> | null = null

const toastRef = ref<InstanceType<typeof NotificationToast> | null>(null)

const deleteDialog = ref({
  isOpen: false,
  title: '',
  message: '',
  bookHash: ''
})

const editDialog = ref({
  isOpen: false,
  title: '',
  authors: '',
  publisher: '',
  bookHash: ''
})

const booksPagination = ref<Pagination>({
  limit: 20,
  offset: 0,
  total: 0,
  has_next: false
})

const favoriteHashes = ref(new Set<string>())

// Compute favorites from books list for instant updates
const favorites = computed(() => {
  return books.value.filter(book => favoriteHashes.value.has(book.hash))
})

const isBookFavorited = (hash: string) => {
  return favoriteHashes.value.has(hash)
}

const loadBooks = async () => {
  try {
    loadingBooks.value = true
    const response = await booksApi.getBooks(20, 0)
    books.value = response.books
    booksPagination.value = response.pagination
  } catch (error) {
    console.error('Failed to load books:', error)
  } finally {
    loadingBooks.value = false
  }
}

const loadFavorites = async () => {
  if (!authStore.isAuthenticated) {
    // Load from local storage
    const localFavorites = localStorage.getItem('marchive_favorites')
    if (localFavorites) {
      try {
        const hashes = JSON.parse(localFavorites) as string[]
        favoriteHashes.value = new Set(hashes)
      } catch (error) {
        console.error('Failed to parse local favorites:', error)
      }
    }
    return
  }

  try {
    const response = await booksApi.getFavorites(20, 0)
    favoriteHashes.value = new Set(response.books.map(book => book.hash))
  } catch (error) {
    console.error('Failed to load favorites:', error)
  }
}

const loadMoreBooks = async () => {
  if (loadingMore.value || !booksPagination.value.has_next) return

  try {
    loadingMore.value = true
    const newOffset = booksPagination.value.offset + booksPagination.value.limit
    const response = await booksApi.getBooks(20, newOffset)
    books.value.push(...response.books)
    booksPagination.value = response.pagination
  } catch (error) {
    console.error('Failed to load more books:', error)
  } finally {
    loadingMore.value = false
  }
}

const handleToggleFavorite = async (book: Book) => {
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

    if (favoriteHashes.value.has(book.hash)) {
      hashes = hashes.filter(h => h !== book.hash)
      favoriteHashes.value.delete(book.hash)
    } else {
      hashes.push(book.hash)
      favoriteHashes.value.add(book.hash)
    }

    localStorage.setItem('marchive_favorites', JSON.stringify(hashes))
    return
  }

  try {
    const response = await booksApi.toggleFavorite(book.hash)
    
    if (response.is_favorited) {
      favoriteHashes.value.add(book.hash)
    } else {
      favoriteHashes.value.delete(book.hash)
    }
  } catch (error) {
    console.error('Failed to toggle favorite:', error)
  }
}

const handleDownload = async (book: Book) => {
  // Check if book is ready for download
  if (book.status !== 'ready') {
    console.log('Book is not ready yet, status:', book.status)
    // TODO: Show a message to user that book is still being processed
    return
  }

  try {
    downloadingBooks.value.add(book.hash)
    // Trigger actual file download
    const downloadUrl = `${import.meta.env.VITE_API_URL || window.location.origin}/api/books/${book.hash}/download`
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = `${book.title}.${book.format}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (error) {
    console.error('Failed to download book:', error)
  } finally {
    downloadingBooks.value.delete(book.hash)
  }
}

const handleOpen = (book: Book) => {
  // This will be implemented later
  console.log('Open book:', book.title)
}

const handleRead = (book: Book) => {
  router.push(`/read/${book.hash}`)
}

const handleViewBook = (book: Book) => {
  // Navigate to book detail page
  router.push(`/book/${book.hash}`)
}

const handleEditBook = (book: Book) => {
  editDialog.value = {
    isOpen: true,
    title: book.title,
    authors: book.authors,
    publisher: book.publisher,
    bookHash: book.hash
  }
}

const saveMetadata = async (data: { title: string; authors: string; publisher: string }) => {
  const bookHash = editDialog.value.bookHash
  editDialog.value.isOpen = false

  try {
    await booksApi.updateBookMetadata(bookHash, data.title, data.authors, data.publisher)
    
    // Update local book data
    const bookIndex = books.value.findIndex(b => b.hash === bookHash)
    if (bookIndex !== -1 && books.value[bookIndex]) {
      books.value[bookIndex]!.title = data.title
      books.value[bookIndex]!.authors = data.authors
      books.value[bookIndex]!.publisher = data.publisher
    }
    
    toastRef.value?.addToast('Book metadata updated successfully!', 'success')
  } catch (error) {
    console.error('Failed to update book metadata:', error)
    toastRef.value?.addToast('Failed to update book metadata', 'error')
  }
}

const cancelEdit = () => {
  editDialog.value.isOpen = false
}

const handleDeleteBook = (book: Book) => {
  deleteDialog.value = {
    isOpen: true,
    title: 'Delete Book',
    message: `Are you sure you want to delete "${book.title}"? This action cannot be undone.`,
    bookHash: book.hash
  }
}

const confirmDelete = async () => {
  const bookHash = deleteDialog.value.bookHash
  deleteDialog.value.isOpen = false

  try {
    await booksApi.deleteBook(bookHash)
    
    // Remove from local state
    books.value = books.value.filter(b => b.hash !== bookHash)
    favoriteHashes.value.delete(bookHash)
    
    toastRef.value?.addToast('Book deleted successfully!', 'success')
  } catch (error) {
    console.error('Failed to delete book:', error)
    toastRef.value?.addToast('Failed to delete book', 'error')
  }
}

const cancelDelete = () => {
  deleteDialog.value.isOpen = false
}

const refreshLibrary = async () => {
  await Promise.all([loadBooks(), loadFavorites()])
}

const startAutoRefresh = () => {
  // Check if any books are processing
  const hasProcessingBooks = books.value.some(book => book.status === 'processing')
  
  if (hasProcessingBooks && !refreshInterval) {
    // Refresh every 5 seconds if there are processing books
    refreshInterval = setInterval(() => {
      loadBooks()
      // Stop auto-refresh if no more processing books
      if (!books.value.some(book => book.status === 'processing')) {
        stopAutoRefresh()
      }
    }, 5000)
  }
  
  // Remove failed books after 10 seconds
  const failedBooks = books.value.filter(book => book.status === 'error')
  failedBooks.forEach(book => {
    setTimeout(() => {
      books.value = books.value.filter(b => b.hash !== book.hash)
    }, 10000)
  })
}

const stopAutoRefresh = () => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
}

onMounted(() => {
  loadBooks().then(() => startAutoRefresh())
  loadFavorites()
})

onUnmounted(() => {
  stopAutoRefresh()
})
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
  justify-content: space-between;
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-left: auto;
}

.refresh-btn {
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 8px;
  padding: 0.5rem;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #3b82f6;
}

.refresh-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.2);
  border-color: #3b82f6;
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.refresh-btn svg {
  width: 20px;
  height: 20px;
  transition: transform 0.3s ease;
}

.refresh-btn svg.spinning {
  animation: spin 1s linear infinite;
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

.loading-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 3rem;
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

.load-more {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}

.load-more-btn {
  padding: 0.75rem 2rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 8px;
  color: #3b82f6;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.load-more-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.5);
  transform: translateY(-2px);
}

.load-more-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner-small {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(59, 130, 246, 0.3);
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
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
