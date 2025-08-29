import { writable } from 'svelte/store';
import { goto } from '$app/navigation';
import { api } from './api';
import { showError, showWarning } from './notifications';

export interface User {
	id: string;
	username: string;
	role: string;
	created_at: string;
}

export const user = writable<User | null>(null);
export const isAuthenticated = writable(false);
export const isAdmin = writable(false);

export const auth = {
	async login(username: string, password: string): Promise<boolean> {
		try {
			const response = await api.post('/auth/login', { username, password });
			if (response.ok) {
				await this.checkAuth();
				return true;
			}
			return false;
		} catch (error) {
			console.error('Login error:', error);
			return false;
		}
	},

	async logout(): Promise<void> {
		try {
			await api.post('/auth/logout', {});
		} catch (error) {
			console.error('Logout error:', error);
		} finally {
			user.set(null);
			isAuthenticated.set(false);
			isAdmin.set(false);
			goto('/login');
		}
	},

	async checkAuth(retryCount = 0): Promise<void> {
		try {
			const response = await api.get('/auth/me');
			if (response.ok) {
				const userData = await response.json();
				user.set(userData);
				isAuthenticated.set(true);
				isAdmin.set(userData.role === 'admin');
			} else if (response.status === 401) {
				// Unauthorized - clear auth state but don't show error on first check
				user.set(null);
				isAuthenticated.set(false);
				isAdmin.set(false);
			} else {
				// Other errors - retry up to 2 times
				if (retryCount < 2) {
					console.log(`Auth check failed, retrying... (${retryCount + 1}/2)`);
					setTimeout(() => this.checkAuth(retryCount + 1), 1000 * (retryCount + 1));
					return;
				}
				
				// Max retries reached
				user.set(null);
				isAuthenticated.set(false);
				isAdmin.set(false);
				showWarning('Connection Issue', 'Unable to verify authentication status. Please refresh the page.');
			}
		} catch (error) {
			console.error('Auth check error:', error);
			
			// Network error - retry up to 2 times
			if (retryCount < 2) {
				console.log(`Auth check failed, retrying... (${retryCount + 1}/2)`);
				setTimeout(() => this.checkAuth(retryCount + 1), 1000 * (retryCount + 1));
				return;
			}
			
			// Max retries reached
			user.set(null);
			isAuthenticated.set(false);
			isAdmin.set(false);
			showError('Connection Error', 'Unable to connect to the server. Please check your internet connection.');
		}
	},

	async changePassword(currentPassword: string, newPassword: string): Promise<boolean> {
		try {
			const response = await api.post('/auth/change-password', {
				current_password: currentPassword,
				new_password: newPassword
			});
			return response.ok;
		} catch (error) {
			console.error('Change password error:', error);
			return false;
		}
	},

	async logoutEverywhere(): Promise<void> {
		try {
			await api.post('/auth/logout-all', {});
		} catch (error) {
			console.error('Logout everywhere error:', error);
		} finally {
			user.set(null);
			isAuthenticated.set(false);
			isAdmin.set(false);
			goto('/login');
		}
	}
};
