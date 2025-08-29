import { writable } from 'svelte/store';
import { api } from './api';

export interface UserWithStats {
	id: string;
	username: string;
	role: string;
	created_at: string;
	download_count: number;
	last_active?: string;
}

export interface SystemStats {
	total_users: number;
	total_books: number;
	total_downloads: number;
	active_users_24h: number;
	recent_downloads: Array<{
		id: string;
		user_id: string;
		md5: string;
		title: string;
		created_at: string;
	}>;
	top_books: Array<{
		hash: string;
		title: string;
		authors: string;
		download_count: number;
	}>;
}

export interface UserListResponse {
	users: UserWithStats[];
	pagination: {
		limit: number;
		offset: number;
		total: number;
		has_next: boolean;
	};
}

export const systemStats = writable<SystemStats | null>(null);
export const usersList = writable<UserListResponse | null>(null);

export const admin = {
	async getSystemStats(): Promise<SystemStats> {
		const response = await api.get('/admin/stats');
		const data = await api.handleResponse<SystemStats>(response);
		systemStats.set(data);
		return data;
	},

	async getUsers(limit = 20, offset = 0): Promise<UserListResponse> {
		const response = await api.get(`/admin/users?limit=${limit}&offset=${offset}`);
		const data = await api.handleResponse<UserListResponse>(response);
		usersList.set(data);
		return data;
	},

	async searchUsers(query?: string, role?: string, limit = 20, offset = 0): Promise<UserListResponse> {
		const response = await api.post('/admin/users/search', { query, role, limit, offset });
		const data = await api.handleResponse<UserListResponse>(response);
		usersList.set(data);
		return data;
	},

	async createUser(userData: {
		username: string;
		password: string;
		role?: string;
	}): Promise<UserWithStats> {
		const response = await api.post('/admin/users', userData);
		return api.handleResponse<UserWithStats>(response);
	},

	async getUser(userId: number): Promise<UserWithStats> {
		const response = await api.get(`/admin/users/${userId}`);
		return api.handleResponse<UserWithStats>(response);
	},

	async updateUser(userId: string, updates: {
		username?: string;
		role?: string;
	}): Promise<UserWithStats> {
		const response = await api.put(`/admin/users/${userId}`, updates);
		return api.handleResponse<UserWithStats>(response);
	},

	async deleteUser(userId: string): Promise<void> {
		const response = await api.delete(`/admin/users/${userId}`);
		await api.handleResponse(response);
	},

	async changeUserPassword(userId: number, newPassword: string): Promise<void> {
		const response = await api.post(`/admin/users/${userId}/password`, { new_password: newPassword });
		await api.handleResponse(response);
	},

	async invalidateUserSessions(userId: string): Promise<void> {
		const response = await api.post(`/admin/users/${userId}/invalidate-sessions`, {});
		await api.handleResponse(response);
	}
};
