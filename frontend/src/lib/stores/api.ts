import { auth } from './auth';

const API_BASE = '/api';

interface ApiError {
	error: string;
}

class ApiClient {
	private isRefreshing = false;
	private refreshPromise: Promise<boolean> | null = null;
	private failedRefreshCount = 0;
	private readonly MAX_REFRESH_ATTEMPTS = 3;

	async request(endpoint: string, options: RequestInit = {}): Promise<Response> {
		const url = `${API_BASE}${endpoint}`;
		
		const config: RequestInit = {
			headers: {
				'Content-Type': 'application/json',
				...options.headers,
			},
			credentials: 'include',
			...options,
		};

		let response = await fetch(url, config);

		// If we get a 401 and this isn't a refresh request, try to refresh the token
		if (response.status === 401 && endpoint !== '/auth/refresh' && !this.isRefreshing && this.failedRefreshCount < this.MAX_REFRESH_ATTEMPTS) {
			console.log(`API request to ${endpoint} failed with 401, attempting token refresh...`);
			
			// Use a single refresh promise to avoid multiple simultaneous refresh attempts
			if (!this.refreshPromise) {
				this.isRefreshing = true;
				console.log('Starting token refresh...');
				this.refreshPromise = auth.handleUnauthorized();
			} else {
				console.log('Token refresh already in progress, waiting...');
			}
			
			const refreshSuccess = await this.refreshPromise;
			
			if (refreshSuccess) {
				console.log('Token refresh successful, retrying original request...');
				// Reset failed count on successful refresh
				this.failedRefreshCount = 0;
				// Retry the original request with the new token
				response = await fetch(url, config);
			} else {
				// Increment failed count
				this.failedRefreshCount++;
				console.log(`Token refresh failed. Failed attempts: ${this.failedRefreshCount}/${this.MAX_REFRESH_ATTEMPTS}`);
			}
			
			// Reset refresh state
			this.isRefreshing = false;
			this.refreshPromise = null;
		} else if (response.status === 401) {
			if (endpoint === '/auth/refresh') {
				console.log('Refresh endpoint returned 401, this is expected for invalid refresh tokens');
			} else if (this.isRefreshing) {
				console.log('Request returned 401 but refresh is already in progress');
			} else if (this.failedRefreshCount >= this.MAX_REFRESH_ATTEMPTS) {
				console.log('Request returned 401 but max refresh attempts reached');
			}
		}

		return response;
	}

	async get(endpoint: string): Promise<Response> {
		return this.request(endpoint, { method: 'GET' });
	}

	async post(endpoint: string, data: any): Promise<Response> {
		return this.request(endpoint, {
			method: 'POST',
			body: JSON.stringify(data),
		});
	}

	async put(endpoint: string, data: any): Promise<Response> {
		return this.request(endpoint, {
			method: 'PUT',
			body: JSON.stringify(data),
		});
	}

	async delete(endpoint: string): Promise<Response> {
		return this.request(endpoint, { method: 'DELETE' });
	}

	async handleResponse<T>(response: Response): Promise<T> {
		if (!response.ok) {
			const error = await response.json() as ApiError;
			const errorWithStatus = new Error(error.error || 'An error occurred');
			(errorWithStatus as any).status = response.status;
			throw errorWithStatus;
		}
		return response.json();
	}

	// Reset failed refresh count (useful after successful login)
	resetFailedRefreshCount(): void {
		console.log('Resetting failed refresh count');
		this.failedRefreshCount = 0;
	}
}

export const api = new ApiClient();
