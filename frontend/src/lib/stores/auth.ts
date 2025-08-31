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
			console.log('Attempting login...');
			const response = await api.post('/auth/login', { username, password });
			if (response.ok) {
				console.log('Login successful, checking auth...');
				// Reset failed refresh count on successful login
				api.resetFailedRefreshCount();
				await this.checkAuthWithRefresh();
				return true;
			}
			console.log('Login failed:', response.status);
			return false;
		} catch (error) {
			console.error('Login error:', error);
			return false;
		}
	},

	async logout(): Promise<void> {
		try {
			console.log('Logging out...');
			await api.post('/auth/logout', {});
		} catch (error) {
			console.error('Logout error:', error);
		} finally {
			console.log('Clearing auth state...');
			user.set(null);
			isAuthenticated.set(false);
			isAdmin.set(false);
			goto('/login');
		}
	},

	/**
	 * Enhanced authentication check that always attempts refresh when /me fails
	 * This is the main method that should be used for checking auth state
	 */
	async checkAuthWithRefresh(): Promise<boolean> {
		try {
			console.log('Checking authentication with automatic refresh...');
			const response = await api.get('/auth/me');
			if (response.ok) {
				const userData = await response.json();
				console.log('Auth check successful:', userData.username);
				user.set(userData);
				isAuthenticated.set(true);
				isAdmin.set(userData.role === 'admin');
				return true;
			}
			
			// Any failure of /me (including 401, 500, network errors, etc.) triggers refresh
			console.log('Auth check failed, attempting token refresh...');
			const refreshSuccess = await this.refreshToken();
			if (refreshSuccess) {
				console.log('Token refresh successful after /me failure');
				return true;
			}
			
			// Refresh failed - clear auth state
			console.log('Token refresh failed after /me failure, clearing auth state...');
			this.clearAuthState();
			return false;
		} catch (error) {
			console.error('Auth check with refresh error:', error);
			
			// Even network errors trigger refresh attempt
			console.log('Network error during auth check, attempting token refresh...');
			const refreshSuccess = await this.refreshToken();
			if (refreshSuccess) {
				console.log('Token refresh successful after network error');
				return true;
			}
			
			// Refresh failed - clear auth state
			console.log('Token refresh failed after network error, clearing auth state...');
			this.clearAuthState();
			return false;
		}
	},

	async refreshToken(): Promise<boolean> {
		try {
			console.log('Refreshing token...');
			const response = await api.post('/auth/refresh', {});
			if (response.ok) {
				console.log('Token refresh successful, verifying...');
				// Token refreshed, verify by calling /me
				const meResponse = await api.get('/auth/me');
				if (meResponse.ok) {
					const userData = await meResponse.json();
					console.log('Token verification successful:', userData.username);
					user.set(userData);
					isAuthenticated.set(true);
					isAdmin.set(userData.role === 'admin');
					// Reset failed refresh count on successful refresh
					api.resetFailedRefreshCount();
					return true;
				}
			}
			
			// Refresh failed
			console.log('Token refresh failed');
			this.clearAuthState();
			return false;
		} catch (error) {
			console.error('Token refresh error:', error);
			this.clearAuthState();
			return false;
		}
	},

	clearAuthState(): void {
		console.log('Clearing authentication state...');
		user.set(null);
		isAuthenticated.set(false);
		isAdmin.set(false);
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
			console.log('Logging out everywhere...');
			await api.post('/auth/logout-all', {});
		} catch (error) {
			console.error('Logout everywhere error:', error);
		} finally {
			console.log('Clearing auth state...');
			user.set(null);
			isAuthenticated.set(false);
			isAdmin.set(false);
			goto('/login');
		}
	},

	async handleUnauthorized(): Promise<boolean> {
		console.log('Handling unauthorized request...');
		// This method can be called when any API request returns 401
		// It will attempt to refresh the token and return true if successful
		return await this.refreshToken();
	}
};
