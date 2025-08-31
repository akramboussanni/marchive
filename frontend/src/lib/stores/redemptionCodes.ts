import { writable } from 'svelte/store';
import { api } from './api';

export interface RedemptionCode {
	id: string;
	code: string;
	description: string;
	invite_tokens: number;
	request_credits: number;
	max_uses: number;
	current_uses: number;
	expires_at?: string;
	revoked_at?: string;
	created_by: string;
	created_at: string;
}

export interface CreateRedemptionCodeRequest {
	code: string;
	description: string;
	invite_tokens: number;
	request_credits: number;
	max_uses: number;
	expires_at?: number;
}

export interface RedemptionCodeListResponse {
	codes: RedemptionCode[];
	pagination: {
		limit: number;
		offset: number;
		total: number;
		has_next: boolean;
	};
}

export interface RedeemCodeRequest {
	code: string;
}

export interface RedeemCodeResponse {
	success: boolean;
	message: string;
	invite_tokens_granted?: number;
	request_credits_granted?: number;
	new_invite_tokens_total?: number;
	new_request_credits_total?: number;
}

class RedemptionCodeStore {
	async createCode(request: CreateRedemptionCodeRequest): Promise<RedemptionCode> {
		const response = await api.post('/admin/redemption-codes', request);
		return api.handleResponse<RedemptionCode>(response);
	}

	async listCodes(limit = 20, offset = 0): Promise<RedemptionCodeListResponse> {
		const response = await api.get(`/admin/redemption-codes?limit=${limit}&offset=${offset}`);
		return api.handleResponse<RedemptionCodeListResponse>(response);
	}

	async revokeCode(codeId: string): Promise<void> {
		const response = await api.post(`/admin/redemption-codes/${codeId}/revoke`, {});
		await api.handleResponse(response);
	}

	async deleteCode(codeId: string): Promise<void> {
		const response = await api.delete(`/admin/redemption-codes/${codeId}`);
		await api.handleResponse(response);
	}

	async redeemCode(request: RedeemCodeRequest): Promise<RedeemCodeResponse> {
		const response = await api.post('/redemption-codes/redeem', request);
		return api.handleResponse<RedeemCodeResponse>(response);
	}
}

export const redemptionCodeStore = new RedemptionCodeStore();

// Store for the current list of codes
export const redemptionCodes = writable<RedemptionCode[]>([]);
export const redemptionCodesLoading = writable(false);
export const redemptionCodesError = writable<string | null>(null);

// Store for pagination
export const redemptionCodesPagination = writable({
	limit: 20,
	offset: 0,
	total: 0,
	has_next: false
});

// Store for create code form
export const createCodeForm = writable<CreateRedemptionCodeRequest>({
	code: '',
	description: '',
	invite_tokens: 0,
	request_credits: 0,
	max_uses: 1
});

// Store for redeem code form
export const redeemCodeForm = writable<RedeemCodeRequest>({
	code: ''
});

// Store for redemption result
export const redemptionResult = writable<RedeemCodeResponse | null>(null);
