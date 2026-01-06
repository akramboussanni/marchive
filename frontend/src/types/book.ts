export interface Book {
  hash: string
  title: string
  authors: string
  publisher: string
  language: string
  format: string
  size: string
  cover_url: string
  cover_data: string
  status: string
  download_count?: number
  is_ghost?: boolean
  requested_by?: string
  created_at?: string
  url?: string
}

export interface Pagination {
  limit: number
  offset: number
  total: number
  has_next: boolean
}

export interface BookListResponse {
  books: Book[]
  pagination: Pagination
}

export interface FavoritesResponse {
  books: Book[]
  pagination: Pagination
}

export interface SearchRequest {
  query: string
  limit?: number
  offset?: number
}

export interface SearchResponse {
  books: Book[]
  total: number
  query: string
  pagination: Pagination
}

export interface ToggleFavoriteRequest {
  book_hash: string
}

export interface ToggleFavoriteResponse {
  is_favorited: boolean
  message: string
}

export interface DownloadRequest {
  hash: string
  title: string
  authors?: string
  publisher?: string
  language?: string
  format?: string
  size?: string
  cover_url?: string
  cover_data?: string
  is_ghost?: boolean
}

export interface BookDetailResponse {
  book: Book
  requested_by?: {
    id: string
    username: string
    role: string
  }
}

