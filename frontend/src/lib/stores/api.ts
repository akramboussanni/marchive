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

		return fetch(url, config);
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
			throw new Error(error.error || 'An error occurred');
		}
		return response.json();
	}
}

export const api = new ApiClient();
