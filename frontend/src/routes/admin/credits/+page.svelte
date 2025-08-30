<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Gift, Plus, Search, RefreshCw } from 'lucide-svelte';
	import { isAuthenticated, isAdmin } from '$lib/stores/auth';

	interface UserWithCredits {
		id: string;
		username: string;
		role: string;
		request_credits: number;
	}

	let users: UserWithCredits[] = [];
	let loading = true;
	let searchQuery = '';
	let showGrantModal = false;
	let selectedUser: UserWithCredits | null = null;
	let grantAmount = 1;
	let grantReason = '';

	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}

		if (!$isAdmin) {
			goto('/');
			return;
		}

		await loadUsers();
	});

	async function loadUsers() {
		loading = true;
		try {
			const response = await fetch('/api/admin/users/credits', {
				credentials: 'include'
			});
			if (response.ok) {
				const data = await response.json();
				users = data.users;
			}
		} catch (error) {
			console.error('Failed to load users:', error);
		} finally {
			loading = false;
		}
	}

	async function grantCredits() {
		if (!selectedUser || grantAmount <= 0) return;

		try {
			const response = await fetch('/api/admin/users/credits/grant', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include',
				body: JSON.stringify({
					user_id: parseInt(selectedUser.id),
					amount: grantAmount,
					reason: grantReason || 'Admin grant'
				})
			});

			if (response.ok) {
				// Refresh the user list
				await loadUsers();
				showGrantModal = false;
				selectedUser = null;
				grantAmount = 1;
				grantReason = '';
			}
		} catch (error) {
			console.error('Failed to grant credits:', error);
		}
	}

	function openGrantModal(user: UserWithCredits) {
		selectedUser = user;
		showGrantModal = true;
	}

	const filteredUsers = users.filter(user =>
		user.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
		user.role.toLowerCase().includes(searchQuery.toLowerCase())
	);
</script>

<svelte:head>
	<title>Request Credits Management - Admin - marchive</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
	<!-- Header -->
	<div class="mb-6 sm:mb-8">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl sm:text-2xl font-bold text-gray-100 flex items-center space-x-2">
					<Gift class="h-5 w-5 sm:h-6 sm:w-6" />
					<span>Request Credits Management</span>
				</h1>
				<p class="text-gray-400 mt-2 text-sm sm:text-base">
					Manage user request credits for downloads beyond daily limits
				</p>
			</div>
			<button
				on:click={loadUsers}
				class="btn-ghost flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-dark-800 transition-colors"
			>
				<RefreshCw class="h-4 w-4" />
				<span>Refresh</span>
			</button>
		</div>
	</div>

	<!-- Search Bar -->
	<div class="mb-6">
		<div class="relative max-w-md">
			<Search class="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search users..."
				class="w-full pl-10 pr-4 py-2 bg-dark-800 border border-gray-700 rounded-lg text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
			/>
		</div>
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-16">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
			<span class="ml-3 text-gray-400">Loading users...</span>
		</div>
	{:else}
		<!-- Users Table -->
		<div class="card overflow-hidden">
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead class="bg-dark-800">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">User</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Role</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Request Credits</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-800">
						{#each filteredUsers as user}
							<tr class="hover:bg-dark-800">
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="text-sm font-medium text-gray-200">{user.username}</div>
									<div class="text-sm text-gray-400">ID: {user.id}</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full {
										user.role === 'admin' ? 'bg-red-100 text-red-800' :
										user.role === 'moderator' ? 'bg-yellow-100 text-yellow-800' :
										'bg-green-100 text-green-800'
									}">
										{user.role}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center space-x-2">
										<Gift class="h-4 w-4 text-primary-400" />
										<span class="text-sm font-medium text-gray-200">{user.request_credits}</span>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<button
										on:click={() => openGrantModal(user)}
										class="btn-ghost flex items-center space-x-2 px-3 py-1.5 text-sm rounded-md hover:bg-primary-500/20 hover:text-primary-400 transition-colors"
									>
										<Plus class="h-4 w-4" />
										<span>Grant Credits</span>
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	{/if}
</div>

<!-- Grant Credits Modal -->
{#if showGrantModal && selectedUser}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-dark-800 rounded-lg p-6 max-w-md w-full mx-4">
			<h3 class="text-lg font-medium text-gray-100 mb-4">
				Grant Request Credits to {selectedUser.username}
			</h3>
			
			<div class="space-y-4">
				<div>
					<label class="block text-sm font-medium text-gray-300 mb-2">Amount</label>
					<input
						type="number"
						bind:value={grantAmount}
						min="1"
						max="100"
						class="w-full px-3 py-2 bg-dark-700 border border-gray-600 rounded-lg text-gray-100 focus:outline-none focus:ring-2 focus:ring-primary-500"
					/>
				</div>
				
				<div>
					<label class="block text-sm font-medium text-gray-300 mb-2">Reason (optional)</label>
					<input
						type="text"
						bind:value={grantReason}
						placeholder="e.g., Bonus for active user"
						class="w-full px-3 py-2 bg-dark-700 border border-gray-600 rounded-lg text-gray-100 focus:outline-none focus:ring-2 focus:ring-primary-500"
					/>
				</div>
			</div>
			
			<div class="flex space-x-3 mt-6">
				<button
					on:click={() => showGrantModal = false}
					class="flex-1 px-4 py-2 bg-gray-600 text-gray-200 rounded-lg hover:bg-gray-500 transition-colors"
				>
					Cancel
				</button>
				<button
					on:click={grantCredits}
					class="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-500 transition-colors"
				>
					Grant Credits
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.card {
		background-color: rgb(31 41 55);
		border: 1px solid rgb(75 85 99);
		border-radius: 0.5rem;
	}

	.btn-ghost {
		background: transparent;
		border: none;
		cursor: pointer;
		color: rgb(156 163 175);
	}
</style>
