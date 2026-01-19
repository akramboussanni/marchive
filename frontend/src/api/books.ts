import apiClient from './client'
import type { BookListResponse, FavoritesResponse, ToggleFavoriteRequest, ToggleFavoriteResponse, SearchRequest, SearchResponse, DownloadRequest, BookDetailResponse } from '@/types/book'

export const booksApi = {
  // Search for books
  async searchBooks(query: string, limit = 20, offset = 0, searchType: 'all' | 'downloaded' | 'missing' = 'all'): Promise<SearchResponse> {
    const response = await apiClient.post('/books/search', { query, limit, offset, search_type: searchType } as SearchRequest)
    return response.data
  },

  // Get all books with pagination
  async getBooks(limit = 20, offset = 0): Promise<BookListResponse> {
    const response = await apiClient.get(`/books/explore?limit=${limit}&offset=${offset}`)
    return response.data
  },

  // Get user's favorites
  async getFavorites(limit = 20, offset = 0): Promise<FavoritesResponse> {
    const response = await apiClient.get(`/books/favorites?limit=${limit}&offset=${offset}`)
    return response.data
  },

  // Toggle favorite status
  async toggleFavorite(bookHash: string): Promise<ToggleFavoriteResponse> {
    const response = await apiClient.post('/books/favorite', { book_hash: bookHash } as ToggleFavoriteRequest)
    return response.data
  },

  // Download book
  async downloadBook(bookData: DownloadRequest) {
    const response = await apiClient.post('/books/download', bookData)
    return response.data
  },

  // Get book detail
  async getBookDetail(hash: string): Promise<BookDetailResponse> {
    const response = await apiClient.get(`/books/${hash}`)
    return response.data
  },

  // Admin: Update ghost mode
  async updateGhostMode(bookHash: string, isGhost: boolean) {
    const response = await apiClient.post('/books/ghost-mode', { book_hash: bookHash, is_ghost: isGhost })
    return response.data
  },

  // Admin: Delete book
  async deleteBook(bookHash: string) {
    const response = await apiClient.post('/books/delete', { book_hash: bookHash })
    return response.data
  },

  // Admin: Update book metadata
  async updateBookMetadata(bookHash: string, title: string, authors: string, publisher: string) {
    const response = await apiClient.post('/books/metadata', {
      book_hash: bookHash,
      title,
      authors,
      publisher
    })
    return response.data
  },

  // Upload book
  async uploadBook(formData: FormData) {
    const response = await apiClient.post('/books/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  },

  // Update book cover
  async updateBookCover(bookHash: string, formData: FormData) {
    const response = await apiClient.put(`/books/${bookHash}/cover`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  },

  // Admin: Restore books from downloads directory
  async restoreBooks() {
    const response = await apiClient.post('/books/restore')
    return response.data
  }
}

