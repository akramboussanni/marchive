import { writable } from 'svelte/store';
import { api } from './api';

export interface Book {
	hash: string;
	title: string;
	authors: string;
	publisher: string;
	language: string;
	format: string;
	size: string;
	cover_url: string;
	cover_data: string;
	status?: string;
	download_count?: number;
	created_at?: string;
}

export interface DownloadJob {
	id: string;
	status: string;
	progress: number;
	error_msg?: string;
	book_hash: string;
	created_at?: string;
	updated_at?: string;
	// Book metadata
	title?: string;
	authors?: string;
	publisher?: string;
	language?: string;
	format?: string;
	size?: string;
	cover_url?: string;
	cover_data?: string;
}

export interface SearchResponse {
	search_id?: string;
	books: Book[];
	total: number;
	query: string;
	pagination: {
		limit: number;
		offset: number;
		total: number;
		has_next: boolean;
	};
	expires_at?: number;
}

export interface BookListResponse {
	books: Array<Book & { download_count: number }>;
	pagination: {
		limit: number;
		offset: number;
		total: number;
		has_next: boolean;
	};
}

export interface FavoritesResponse {
	books: Array<Book & { download_count: number }>;
	pagination: {
		limit: number;
		offset: number;
		total: number;
		has_next: boolean;
	};
}

export const searchResults = writable<SearchResponse | null>(null);
export const exploredBooks = writable<BookListResponse | null>(null);
export const downloadJobs = writable<DownloadJob[]>([]);
export const isSearching = writable(false);
export const userFavorites = writable<FavoritesResponse | null>(null);
export const isFavoritesLoading = writable(false);

export const books = {
	async search(query: string, limit = 20, offset = 0): Promise<SearchResponse> {
		isSearching.set(true);
		try {
			const response = await api.post('/books/search', { query, limit, offset });
			const data = await api.handleResponse<SearchResponse>(response);
			searchResults.set(data);
			return data;
		} finally {
			isSearching.set(false);
		}
	},

	async explore(limit = 24, offset = 0): Promise<BookListResponse> {
		const response = await api.get(`/books/explore?limit=${limit}&offset=${offset}`);
		const data = await api.handleResponse<BookListResponse>(response);
		exploredBooks.set(data);
		return data;
	},

	async requestDownload(hash: string, title: string): Promise<DownloadJob> {
		const response = await api.post('/books/download', { hash, title });
		const data = await api.handleResponse<{ job_id: string; status: string; message: string }>(response);
		
		const job: DownloadJob = {
			id: data.job_id,
			status: data.status,
			progress: 0,
			book_hash: hash
		};
		
		downloadJobs.update(jobs => [...jobs, job]);
		return job;
	},

	async requestCachedDownload(searchId: string, index: number, book: Book): Promise<DownloadJob> {
		const response = await api.post('/books/download/cached', { 
			search_id: searchId, 
			index: index 
		});
		const data = await api.handleResponse<{ job_id: string; status: string; message: string }>(response);
		
		const job: DownloadJob = {
			id: data.job_id,
			status: data.status,
			progress: 0,
			book_hash: book.hash,
			title: book.title,
			authors: book.authors,
			format: book.format,
			size: book.size,
			cover_url: book.cover_url,
			cover_data: book.cover_data
		};
		
		downloadJobs.update(jobs => [...jobs, job]);
		return job;
	},

	async getJobStatus(jobId: string): Promise<DownloadJob> {
		const response = await api.get(`/books/job/${jobId}`);
		return api.handleResponse<DownloadJob>(response);
	},

	async getUserDownloads(limit = 20, offset = 0): Promise<{ jobs: DownloadJob[]; pagination: { limit: number; offset: number; total: number; has_next: boolean } }> {
		const response = await api.get(`/books/downloads?limit=${limit}&offset=${offset}`);
		return api.handleResponse<{ jobs: DownloadJob[]; pagination: { limit: number; offset: number; total: number; has_next: boolean } }>(response);
	},

	getDownloadUrl(hash: string): string {
		return `/api/books/${hash}/download`;
	},

	async downloadBook(hash: string, title: string): Promise<Blob> {
		const response = await api.get(`/books/${hash}/download`);
		if (!response.ok) {
			throw new Error('Download failed');
		}
		return response.blob();
	},

	updateJobStatus(jobId: string, status: Partial<DownloadJob>): void {
		downloadJobs.update(jobs => 
			jobs.map(job => job.id === jobId ? { ...job, ...status } : job)
		);
	},

	async toggleFavorite(bookHash: string): Promise<{ isFavorited: boolean; message: string }> {
		const response = await api.post('/books/favorite', { book_hash: bookHash });
		const data = await api.handleResponse<{ is_favorited: boolean; message: string }>(response);
		return {
			isFavorited: data.is_favorited,
			message: data.message
		};
	},

	async getUserFavorites(limit = 24, offset = 0): Promise<FavoritesResponse> {
		isFavoritesLoading.set(true);
		try {
			const response = await api.get(`/books/favorites?limit=${limit}&offset=${offset}`);
			const data = await api.handleResponse<FavoritesResponse>(response);
			userFavorites.set(data);
			return data;
		} finally {
			isFavoritesLoading.set(false);
		}
	}
};
