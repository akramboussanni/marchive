<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="modal-overlay" @click="closeModal">
        <div class="modal-container" @click.stop>
          <div class="modal-header">
            <h2>Search Books</h2>
            <button @click="closeModal" class="close-btn" title="Close">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <div class="modal-body">
            <SearchBar 
              v-model="searchQuery"
              placeholder="Search for books by title, author, or ISBN..."
              :searching="searching"
              @search="handleSearch"
              @clear="clearSearch"
            />

            <div v-if="searchError" class="search-error">
              {{ searchError }}
            </div>

            <div v-if="searching" class="loading-state">
              <div class="spinner"></div>
              <p>Searching...</p>
            </div>

            <div v-else-if="currentSearchQuery && searchResults.length === 0 && !searchError" class="empty-state">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"></circle>
                <path d="m21 21-4.35-4.35"></path>
              </svg>
              <p>No results found for "{{ currentSearchQuery }}"</p>
            </div>

            <div v-else-if="searchResults.length > 0" class="results-container">
              <div class="results-header">
                <span class="result-count">
                  {{ searchPagination.total }} {{ searchPagination.total === 1 ? 'result' : 'results' }} for "{{ currentSearchQuery }}"
                </span>
              </div>

              <div class="books-grid">
                <BookCard 
                  v-for="book in searchResults"
                  :key="book.hash"
                  :book="book"
                  :is-available="book.status === 'available'"
                  :is-downloading="downloadingBooks.has(book.hash)"
                  action-mode="add"
                  @add-to-library="handleAddToLibrary"
                />
              </div>

              <div v-if="searchPagination.has_next" class="load-more">
                <button @click="loadMoreSearchResults" :disabled="loadingMore" class="load-more-btn">
                  <span v-if="!loadingMore">Load More</span>
                  <div v-else class="spinner-small"></div>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Ghost Mode Dialog -->
    <Transition name="modal">
      <div v-if="showGhostModeDialog" class="modal-overlay" @click="cancelAddToLibrary">
        <div class="ghost-dialog" @click.stop>
          <h3>Add to Marchive</h3>
          <p class="dialog-description">Choose how to add this book to your library:</p>
          
          <div class="ghost-mode-option">
            <label class="checkbox-label">
              <input 
                type="checkbox" 
                v-model="ghostModeEnabled"
                class="checkbox-input"
              />
              <span class="checkbox-text">
                <svg class="ghost-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 2a7 7 0 0 0-7 7v3l-2 2v5h18v-5l-2-2V9a7 7 0 0 0-7-7z"></path>
                  <path d="M8.5 19a4 4 0 0 0 7 0"></path>
                  <circle cx="12" cy="11" r="1" fill="currentColor"></circle>
                </svg>
                <strong>Ghost Mode</strong>
              </span>
            </label>
            <p class="ghost-description">
              Book will only be visible to you and admins. Others won't see it in the public library.
            </p>
          </div>

          <div v-if="addBookError" class="error-message">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="15" y1="9" x2="9" y2="15"></line>
              <line x1="9" y1="9" x2="15" y2="15"></line>
            </svg>
            <span>{{ addBookError }}</span>
          </div>

          <div class="dialog-actions">
            <button class="dialog-btn cancel-btn" @click="cancelAddToLibrary">
              Cancel
            </button>
            <button class="dialog-btn confirm-btn" @click="confirmAddToLibrary(ghostModeEnabled)">
              Add Book
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { booksApi } from '@/api/books'
import { useAuthStore } from '@/stores/auth'
import BookCard from './BookCard.vue'
import SearchBar from './SearchBar.vue'
import type { Book, Pagination } from '@/types/book'

interface Props {
  isOpen: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
  bookAdded: []
}>()

const authStore = useAuthStore()

const searchQuery = ref('')
const currentSearchQuery = ref('')
const searchResults = ref<Book[]>([])
const searching = ref(false)
const searchError = ref('')
const loadingMore = ref(false)
const downloadingBooks = ref(new Set<string>())
const showGhostModeDialog = ref(false)
const selectedBook = ref<Book | null>(null)
const ghostModeEnabled = ref(false)
const addBookError = ref('')

const searchPagination = ref<Pagination>({
  limit: 20,
  offset: 0,
  total: 0,
  has_next: false
})

const handleAddToLibrary = async (book: Book) => {
  selectedBook.value = book
  showGhostModeDialog.value = true
}

const confirmAddToLibrary = async (isGhost: boolean) => {
  if (!selectedBook.value) return

  try {
    addBookError.value = ''
    downloadingBooks.value.add(selectedBook.value.hash)
    showGhostModeDialog.value = false
    
    const response = await booksApi.downloadBook({
      hash: selectedBook.value.hash,
      title: selectedBook.value.title,
      authors: selectedBook.value.authors,
      publisher: selectedBook.value.publisher,
      language: selectedBook.value.language,
      format: selectedBook.value.format,
      size: selectedBook.value.size,
      cover_url: selectedBook.value.cover_url,
      cover_data: selectedBook.value.cover_data,
      is_ghost: isGhost
    })
    
    // Update book status in results based on response
    const bookIndex = searchResults.value.findIndex(b => b.hash === selectedBook.value!.hash)
    if (bookIndex !== -1 && searchResults.value[bookIndex]) {
      searchResults.value[bookIndex]!.status = response.status || 'processing'
    }
    
    // Close the modal and navigate to home
    closeModal()
    emit('bookAdded')
  } catch (error: any) {
    console.error('Failed to add book to library:', error)
    addBookError.value = error.response?.data?.message || error.message || 'Failed to add book to library. Please try again.'
    showGhostModeDialog.value = true
  } finally {
    if (selectedBook.value) {
      downloadingBooks.value.delete(selectedBook.value.hash)
    }
  }
}

const cancelAddToLibrary = () => {
  showGhostModeDialog.value = false
  selectedBook.value = null
  ghostModeEnabled.value = false
  addBookError.value = ''
}

const handleSearch = async (query?: string) => {
  const searchTerm = (query || searchQuery.value).trim()
  if (!searchTerm || searching.value) return

  try {
    searching.value = true
    searchError.value = ''
    
    const response = await booksApi.searchBooks(searchTerm, 20, 0)
    searchResults.value = response.books
    searchPagination.value = response.pagination
    currentSearchQuery.value = searchTerm
  } catch (error: any) {
    console.error('Failed to search books:', error)
    searchError.value = error.response?.data?.message || 'Failed to search books. Please try again.'
    searchResults.value = []
  } finally {
    searching.value = false
  }
}

const clearSearch = () => {
  searchQuery.value = ''
  searchResults.value = []
  currentSearchQuery.value = ''
  searchError.value = ''
  searchPagination.value = {
    limit: 20,
    offset: 0,
    total: 0,
    has_next: false
  }
}

const loadMoreSearchResults = async () => {
  if (loadingMore.value || !searchPagination.value.has_next || !currentSearchQuery.value) return

  try {
    loadingMore.value = true
    const newOffset = searchPagination.value.offset + searchPagination.value.limit
    const response = await booksApi.searchBooks(currentSearchQuery.value, 20, newOffset)
    searchResults.value.push(...response.books)
    searchPagination.value = response.pagination
  } catch (error) {
    console.error('Failed to load more search results:', error)
  } finally {
    loadingMore.value = false
  }
}

const closeModal = () => {
  emit('close')
}

// Close on Escape key
watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        closeModal()
      }
    }
    document.addEventListener('keydown', handleEscape)
    return () => document.removeEventListener('keydown', handleEscape)
  }
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
  overflow-y: auto;
}

.modal-container {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 16px;
  max-width: 1200px;
  width: 100%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  margin: auto;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid rgba(59, 130, 246, 0.2);
}

.modal-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0;
}

.close-btn {
  background: transparent;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.close-btn svg {
  width: 24px;
  height: 24px;
}

.modal-body {
  padding: 2rem;
  overflow-y: auto;
  flex: 1;
}

.search-error {
  padding: 0.75rem 1rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  color: #fca5a5;
  font-size: 0.9rem;
  margin-top: 1rem;
}

.loading-state {
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

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 4rem 2rem;
  color: #64748b;
  text-align: center;
}

.empty-state svg {
  width: 64px;
  height: 64px;
  color: #475569;
}

.results-container {
  margin-top: 1.5rem;
}

.results-header {
  margin-bottom: 1.5rem;
}

.result-count {
  font-size: 0.9rem;
  color: #64748b;
}

.books-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1.5rem;
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

.ghost-dialog {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  border: 1px solid rgba(147, 51, 234, 0.3);
  border-radius: 16px;
  padding: 2rem;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.ghost-dialog h3 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0 0 0.5rem 0;
}

.dialog-description {
  color: #94a3b8;
  margin: 0 0 1.5rem 0;
  font-size: 0.95rem;
}

.ghost-mode-option {
  background: rgba(147, 51, 234, 0.1);
  border: 1px solid rgba(147, 51, 234, 0.2);
  border-radius: 12px;
  padding: 1.25rem;
  margin-bottom: 1.5rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  margin-bottom: 0.75rem;
}

.checkbox-input {
  width: 20px;
  height: 20px;
  cursor: pointer;
  accent-color: #9333ea;
}

.checkbox-text {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #e2e8f0;
  font-size: 1rem;
}

.ghost-icon {
  width: 20px;
  height: 20px;
  color: #9333ea;
}

.ghost-description {
  color: #94a3b8;
  font-size: 0.875rem;
  margin: 0;
  line-height: 1.5;
  padding-left: 2rem;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  color: #f87171;
  font-size: 0.875rem;
  animation: slideDown 0.3s ease;
}

.error-message svg {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dialog-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.dialog-btn {
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.cancel-btn {
  background: rgba(100, 116, 139, 0.1);
  color: #94a3b8;
  border: 1px solid rgba(100, 116, 139, 0.3);
}

.cancel-btn:hover {
  background: rgba(100, 116, 139, 0.2);
  color: #cbd5e1;
}

.confirm-btn {
  background: rgba(147, 51, 234, 0.2);
  color: #c084fc;
  border: 1px solid rgba(147, 51, 234, 0.3);
}

.confirm-btn:hover {
  background: rgba(147, 51, 234, 0.3);
  color: #e9d5ff;
  border-color: rgba(147, 51, 234, 0.5);
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
  transition: transform 0.3s ease;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.9);
}

@media (max-width: 1024px) {
  .books-grid {
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 1.25rem;
  }
}

@media (max-width: 768px) {
  .modal-container {
    max-height: 95vh;
  }

  .modal-header {
    padding: 1.25rem 1.5rem;
  }

  .modal-header h2 {
    font-size: 1.25rem;
  }

  .modal-body {
    padding: 1.5rem;
  }

  .books-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 1rem;
  }
}

@media (max-width: 480px) {
  .modal-overlay {
    padding: 0;
  }

  .modal-container {
    max-height: 100vh;
    border-radius: 0;
  }

  .modal-header {
    padding: 1rem 1.25rem;
  }

  .modal-body {
    padding: 1rem;
  }

  .books-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 0.875rem;
  }
}
</style>
