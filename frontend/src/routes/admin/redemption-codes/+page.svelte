<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Plus, Trash2, X, Calendar, Users, Gift, AlertCircle } from 'lucide-svelte';
	import { isAuthenticated, isAdmin } from '$lib/stores/auth';
	import { 
		redemptionCodeStore, 
		redemptionCodes, 
		redemptionCodesLoading, 
		redemptionCodesError,
		redemptionCodesPagination,
		createCodeForm,
		type RedemptionCode,
		type CreateRedemptionCodeRequest
	} from '$lib/stores/redemptionCodes';
	import { notifications } from '$lib/stores/notifications';
	import { Pagination } from '$lib/components/UI/Pagination.svelte';

	let showCreateModal = false;
	let loading = false;
	let currentPage = 0;

	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}

		if (!$isAdmin) {
			goto('/');
			return;
		}

		await loadCodes();
	});

	async function loadCodes() {
		redemptionCodesLoading.set(true);
		redemptionCodesError.set(null);
		
		try {
			const response = await redemptionCodeStore.listCodes(20, currentPage * 20);
			redemptionCodes.set(response.codes);
			redemptionCodesPagination.set(response.pagination);
		} catch (error) {
			console.error('Failed to load redemption codes:', error);
			redemptionCodesError.set('Failed to load redemption codes');
		} finally {
			redemptionCodesLoading.set(false);
		}
	}

	async function handleCreateCode() {
		loading = true;
		try {
			await redemptionCodeStore.createCode($createCodeForm);
			notifications.add('success', 'Redemption code created successfully!');
			showCreateModal = false;
			createCodeForm.set({
				code: '',
				description: '',
				invite_tokens: 0,
				request_credits: 0,
				max_uses: 1
			});
			await loadCodes();
		} catch (error) {
			console.error('Failed to create redemption code:', error);
			notifications.add('error', 'Failed to create redemption code');
		} finally {
			loading = false;
		}
	}

	async function handleRevokeCode(code: RedemptionCode) {
		if (!confirm(`Are you sure you want to revoke the code "${code.code}"?`)) {
			return;
		}

		try {
			await redemptionCodeStore.revokeCode(code.id);
			notifications.add('success', 'Code revoked successfully!');
			await loadCodes();
		} catch (error) {
			console.error('Failed to revoke code:', error);
			notifications.add('error', 'Failed to revoke code');
		}
	}

	async function handleDeleteCode(code: RedemptionCode) {
		if (!confirm(`Are you sure you want to permanently delete the code "${code.code}"? This action cannot be undone.`)) {
			return;
		}

		try {
			await redemptionCodeStore.deleteCode(code.id);
			notifications.add('success', 'Code deleted successfully!');
			await loadCodes();
		} catch (error) {
			console.error('Failed to delete code:', error);
			notifications.add('error', 'Failed to delete code');
		}
	}

	async function handleGiveEveryoneInvite() {
		if (!confirm('Are you sure you want to give 1 invite token to ALL users? This action cannot be undone.')) {
			return;
		}

		loading = true;
		try {
			await redemptionCodeStore.giveEveryoneInvite();
			notifications.add('success', '1 invite token given to all users successfully!');
		} catch (error) {
			console.error('Failed to give everyone invite:', error);
			notifications.add('error', 'Failed to give everyone invite');
		} finally {
			loading = false;
		}
	}

	function formatDate(timestamp: string | number): string {
		const ts = typeof timestamp === 'string' ? parseInt(timestamp) : timestamp;
		return new Date(ts * 1000).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getStatusBadge(code: RedemptionCode) {
		if (code.revoked_at) {
			return { text: 'Revoked', class: 'bg-red-100 text-red-800' };
		}
		if (code.expires_at && parseInt(code.expires_at) < Date.now() / 1000) {
			return { text: 'Expired', class: 'bg-yellow-100 text-yellow-800' };
		}
		if (code.current_uses >= code.max_uses) {
			return { text: 'Maxed Out', class: 'bg-gray-100 text-gray-800' };
		}
		return { text: 'Active', class: 'bg-green-100 text-green-800' };
	}

	async function handlePageChange(page: number) {
		currentPage = page;
		await loadCodes();
	}
</script>

<svelte:head>
	<title>Redemption Codes - Admin Dashboard</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<!-- Header -->
	<div class="mb-8">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-3xl font-bold text-gray-900">Redemption Codes</h1>
				<p class="text-gray-600 mt-2">Manage redemption codes that users can redeem for rewards</p>
			</div>
			<div class="flex space-x-3">
				<button
					on:click={handleGiveEveryoneInvite}
					disabled={loading}
					class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors flex items-center space-x-2"
				>
					<Gift class="w-4 h-4" />
					<span>Give Everyone 1 Invite</span>
				</button>
				<button
					on:click={() => showCreateModal = true}
					class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-2"
				>
					<Plus class="w-4 h-4" />
					<span>Create Code</span>
				</button>
			</div>
		</div>
	</div>

	<!-- Codes List -->
	<div class="bg-white rounded-lg shadow-sm border border-gray-200">
		<div class="px-6 py-4 border-b border-gray-200">
			<h3 class="text-lg font-semibold text-gray-900">All Redemption Codes</h3>
		</div>

		{#if $redemptionCodesLoading}
			<div class="px-6 py-8 text-center">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
				<p class="text-gray-500 mt-2">Loading codes...</p>
			</div>
		{:else if $redemptionCodesError}
			<div class="px-6 py-8 text-center">
				<AlertCircle class="w-12 h-12 text-red-400 mx-auto mb-3" />
				<p class="text-red-600">{$redemptionCodesError}</p>
				<button
					on:click={loadCodes}
					class="mt-3 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
				>
					Try Again
				</button>
			</div>
		{:else if $redemptionCodes.length === 0}
			<div class="px-6 py-8 text-center text-gray-500">
				<Gift class="w-16 h-16 text-gray-300 mx-auto mb-3" />
				<p>No redemption codes created yet.</p>
				<p class="text-sm mt-1">Click "Create Code" above to get started!</p>
			</div>
		{:else}
			<div class="divide-y divide-gray-200">
				{#each $redemptionCodes as code}
					<div class="px-6 py-4">
						<div class="flex items-center justify-between">
							<div class="flex-1">
								<div class="flex items-center space-x-3 mb-2">
									<span class="font-mono text-lg bg-gray-100 px-3 py-1 rounded text-gray-700">
										{code.code}
									</span>
									<span class={`px-2 py-1 rounded-full text-xs font-medium ${getStatusBadge(code).class}`}>
										{getStatusBadge(code).text}
									</span>
								</div>
								
								<p class="text-gray-900 font-medium mb-1">{code.description}</p>
								
								<div class="flex items-center space-x-6 text-sm text-gray-600">
									<div class="flex items-center space-x-1">
										<Users class="w-4 h-4" />
										<span>{code.current_uses} / {code.max_uses} uses</span>
									</div>
									
									{#if code.invite_tokens > 0}
										<div class="flex items-center space-x-1">
											<Gift class="w-4 h-4" />
											<span>{code.invite_tokens} invite tokens</span>
										</div>
									{/if}
									
									{#if code.request_credits > 0}
										<div class="flex items-center space-x-1">
											<Gift class="w-4 h-4" />
											<span>{code.request_credits} request credits</span>
										</div>
									{/if}
									
									{#if code.expires_at}
										<div class="flex items-center space-x-1">
											<Calendar class="w-4 h-4" />
											<span>Expires: {formatDate(code.expires_at)}</span>
										</div>
									{/if}
									
									<div class="text-gray-500">
										Created: {formatDate(code.created_at)}
									</div>
								</div>
							</div>
							
							<div class="flex space-x-2">
								{#if getStatusBadge(code).text === 'Active'}
									<button
										on:click={() => handleRevokeCode(code)}
										class="px-3 py-1 text-sm bg-yellow-100 text-yellow-700 rounded hover:bg-yellow-200 transition-colors"
									>
										Revoke
									</button>
								{/if}
								<button
									on:click={() => handleDeleteCode(code)}
									class="px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200 transition-colors"
								>
									<Trash2 class="w-4 h-4" />
								</button>
							</div>
						</div>
					</div>
				{/each}
			</div>

			<!-- Pagination -->
			{#if $redemptionCodesPagination.total > $redemptionCodesPagination.limit}
				<div class="px-6 py-4 border-t border-gray-200">
					<Pagination
						currentPage={currentPage}
						totalPages={Math.ceil($redemptionCodesPagination.total / $redemptionCodesPagination.limit)}
						onPageChange={handlePageChange}
					/>
				</div>
			{/if}
		{/if}
	</div>
</div>

<!-- Create Code Modal -->
{#if showCreateModal}
	<div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
		<div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
			<div class="mt-3">
				<div class="flex items-center justify-between mb-4">
					<h3 class="text-lg font-medium text-gray-900">Create Redemption Code</h3>
					<button
						on:click={() => showCreateModal = false}
						class="text-gray-400 hover:text-gray-600"
					>
						<X class="w-5 h-5" />
					</button>
				</div>

				<form on:submit|preventDefault={handleCreateCode} class="space-y-4">
					<div>
						<label for="code" class="block text-sm font-medium text-gray-700">Code</label>
						<input
							id="code"
							type="text"
							bind:value={$createCodeForm.code}
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
							placeholder="e.g., WELCOME2024"
							required
						/>
					</div>

					<div>
						<label for="description" class="block text-sm font-medium text-gray-700">Description</label>
						<input
							id="description"
							type="text"
							bind:value={$createCodeForm.description}
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
							placeholder="e.g., Welcome bonus for new users"
							required
						/>
					</div>

					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="invite_tokens" class="block text-sm font-medium text-gray-700">Invite Tokens</label>
							<input
								id="invite_tokens"
								type="number"
								bind:value={$createCodeForm.invite_tokens}
								min="0"
								max="100"
								class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
							/>
						</div>

						<div>
							<label for="request_credits" class="block text-sm font-medium text-gray-700">Request Credits</label>
							<input
								id="request_credits"
								type="number"
								bind:value={$createCodeForm.request_credits}
								min="0"
								max="1000"
								class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
							/>
						</div>
					</div>

					<div>
						<label for="max_uses" class="block text-sm font-medium text-gray-700">Max Uses</label>
						<input
							id="max_uses"
							type="number"
							bind:value={$createCodeForm.max_uses}
							min="1"
							max="1000000"
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
							required
						/>
					</div>

					<div>
						<label for="expires_at" class="block text-sm font-medium text-gray-700">Expires At (Optional)</label>
						<input
							id="expires_at"
							type="datetime-local"
							on:change={(e) => {
								if (e.target.value) {
									$createCodeForm.expires_at = Math.floor(new Date(e.target.value).getTime() / 1000);
								} else {
									$createCodeForm.expires_at = undefined;
								}
							}}
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
						/>
					</div>

					<div class="flex justify-end space-x-3 pt-4">
						<button
							type="button"
							on:click={() => showCreateModal = false}
							class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
						>
							Cancel
						</button>
						<button
							type="submit"
							disabled={loading || ($createCodeForm.invite_tokens === 0 && $createCodeForm.request_credits === 0)}
							class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
						>
							{loading ? 'Creating...' : 'Create Code'}
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}
