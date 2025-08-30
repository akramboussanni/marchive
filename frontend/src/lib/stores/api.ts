import { auth } from './auth';

const API_BASE = '/api';

interface ApiError {
	error: string;
}

class ApiClient {
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
		if (response.status === 401 && endpoint !== '/auth/refresh') {
			const refreshSuccess = await auth.handleUnauthorized();
			if (refreshSuccess) {
				// Retry the original request with the new token
				response = await fetch(url, config);
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
}

export const api = new ApiClient();
