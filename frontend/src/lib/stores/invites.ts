import { writable } from 'svelte/store';
import { api } from './api';
import { showError, showSuccess } from './notifications';

export interface Invite {
	id: string;
	token: string;
	inviter_id: string;
	invitee_username?: string;
	invitee_id?: string;
	used_at?: string;
	revoked_at?: string;
	created_at: string;
}

export interface InviteResponse {
	token: string;
	invite_url: string;
	created_at: string;
}

export interface InviteListResponse {
	invites: Invite[];
	tokens: number;
}

export interface UseInviteRequest {
	token: string;
	username: string;
	password: string;
}

export const invites = writable<Invite[]>([]);
export const inviteTokens = writable(0);

export const inviteStore = {
	async createInvite(): Promise<InviteResponse | null> {
		try {
			const response = await api.post('/invites', {});
			if (response.ok) {
				const invite = await response.json();
				// Refresh the list and tokens
				await this.loadInvites();
				showSuccess('Invite Created', 'Your invite has been generated successfully!');
				return invite;
			} else {
				const error = await response.json();
				showError('Failed to Create Invite', error.error || 'An error occurred');
				return null;
			}
		} catch (error) {
			console.error('Create invite error:', error);
			showError('Failed to Create Invite', 'An error occurred while creating the invite');
			return null;
		}
	},

	async loadInvites(): Promise<void> {
		try {
			const response = await api.get('/invites');
			if (response.ok) {
				const data: InviteListResponse = await response.json();
				invites.set(data.invites);
				inviteTokens.set(data.tokens);
			} else {
				console.error('Failed to load invites:', response.status);
			}
		} catch (error) {
			console.error('Load invites error:', error);
		}
	},

	async revokeInvite(token: string): Promise<boolean> {
		try {
			const response = await api.post(`/invites/${token}/revoke`, {});
			if (response.ok) {
				// Refresh the list and tokens
				await this.loadInvites();
				showSuccess('Invite Revoked', 'Your invite has been revoked and token returned');
				return true;
			} else {
				const error = await response.json();
				showError('Failed to Revoke Invite', error.error || 'An error occurred');
				return false;
			}
		} catch (error) {
			console.error('Revoke invite error:', error);
			showError('Failed to Revoke Invite', 'An error occurred while revoking the invite');
			return false;
		}
	},

	async useInvite(token: string, username: string, password: string): Promise<boolean> {
		try {
			const request: UseInviteRequest = { token, username, password };
			const response = await api.post('/invites/use', request);
			if (response.ok) {
				showSuccess('Account Created', 'Your account has been created successfully!');
				return true;
			} else {
				const error = await response.json();
				showError('Failed to Create Account', error.error || 'An error occurred');
				return false;
			}
		} catch (error) {
			console.error('Use invite error:', error);
			showError('Failed to Create Account', 'An error occurred while creating your account');
			return false;
		}
	}
};
