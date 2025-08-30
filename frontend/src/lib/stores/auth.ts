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
				await this.checkAuth();
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

	async checkAuth(retryCount = 0): Promise<void> {
		try {
			console.log('Checking authentication...');
			const response = await api.get('/auth/me');
			if (response.ok) {
				const userData = await response.json();
				console.log('Auth check successful:', userData.username);
				user.set(userData);
				isAuthenticated.set(true);
				isAdmin.set(userData.role === 'admin');
				return;
			} else if (response.status === 401) {
				console.log('Auth check failed (401), attempting token refresh...');
				// Unauthorized - try to refresh the token first
				const refreshSuccess = await this.refreshToken();
				if (refreshSuccess) {
					// Token refreshed successfully, we're now authenticated
					return;
				}
				
				// Refresh failed - clear auth state
				console.log('Token refresh failed, clearing auth state...');
				this.clearAuthState();
			} else {
				// Other errors - retry up to 2 times
				if (retryCount < 2) {
					console.log(`Auth check failed, retrying... (${retryCount + 1}/2)`);
					setTimeout(() => this.checkAuth(retryCount + 1), 1000 * (retryCount + 1));
					return;
				}
				
				// Max retries reached
				console.log('Max auth check retries reached');
				this.clearAuthState();
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
			console.log('Max auth check retries reached');
			this.clearAuthState();
			showError('Connection Error', 'Unable to connect to the server. Please check your internet connection.');
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

	async makeAuthenticatedRequest<T>(
		requestFn: () => Promise<T>,
		onAuthFailure?: () => void
	): Promise<T | null> {
		try {
			return await requestFn();
		} catch (error: any) {
			if (error?.status === 401) {
				console.log('Request failed with 401, attempting token refresh...');
				// Try to refresh the token
				const refreshSuccess = await this.refreshToken();
				if (refreshSuccess) {
					console.log('Token refreshed, retrying request...');
					// Retry the request
					try {
						return await requestFn();
					} catch (retryError) {
						console.error('Request failed after token refresh:', retryError);
						return null;
					}
				} else {
					console.log('Token refresh failed, calling failure handler...');
					// Refresh failed, call the failure handler
					onAuthFailure?.();
					return null;
				}
			}
			throw error;
		}
	},

	async handleUnauthorized(): Promise<boolean> {
		console.log('Handling unauthorized request...');
		// This method can be called when any API request returns 401
		// It will attempt to refresh the token and return true if successful
		return await this.refreshToken();
	}
};
