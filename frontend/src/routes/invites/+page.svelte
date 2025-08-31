<script lang="ts">
	import { onMount } from 'svelte';
	import { invites, inviteTokens, inviteStore } from '$lib/stores/invites';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	onMount(async () => {
		if (!$isAuthenticated) {
			// Try to refresh the token before redirecting to login
			const authSuccess = await auth.checkAuthWithRefresh();
			if (!authSuccess) {
				goto('/login');
				return;
			}
		}
		await inviteStore.loadInvites();
	});

	async function handleCreateInvite() {
		const invite = await inviteStore.createInvite();
		if (invite) {
			// Copy invite URL to clipboard
			navigator.clipboard.writeText(invite.invite_url);
			// Refresh the page data
			await inviteStore.loadInvites();
		}
	}

	async function handleRevokeInvite(token: string) {
		if (confirm('Are you sure you want to revoke this invite? You will get your token back.')) {
			await inviteStore.revokeInvite(token);
			// Refresh the page data
			await inviteStore.loadInvites();
		}
	}

	function formatDate(timestamp: string): string {
		return new Date(parseInt(timestamp) * 1000).toLocaleDateString();
	}

	function getInviteStatus(invite: any): string {
		if (invite.revoked_at) return 'Revoked';
		if (invite.used_at) return 'Used';
		return 'Active';
	}

	function getStatusColor(status: string): string {
		switch (status) {
			case 'Active': return 'text-green-400 bg-green-900/20 border border-green-800';
			case 'Used': return 'text-primary-400 bg-primary-900/20 border border-primary-800';
			case 'Revoked': return 'text-red-400 bg-red-900/20 border border-red-800';
			default: return 'text-gray-400 bg-gray-900/20 border border-gray-700';
		}
	}
</script>

<svelte:head>
	<title>Invites - Marchive</title>
</svelte:head>

<div class="min-h-screen bg-dark-950 text-gray-100">
	<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="text-3xl font-bold text-gray-100 mb-2">Invites</h1>
			<p class="text-gray-400">
				Invite others to join Marchive. Each user gets 1 invite token.
			</p>
		</div>

		<!-- Invite Token Status -->
		<div class="bg-dark-800 rounded-lg border border-gray-700 p-6 mb-6">
			<div class="flex items-center justify-between">
				<div>
					<h2 class="text-lg font-semibold text-gray-100">Your Invite Tokens</h2>
					<p class="text-gray-400">You have {$inviteTokens} invite token{$inviteTokens !== 1 ? 's' : ''} available</p>
				</div>
				<button
					on:click={handleCreateInvite}
					disabled={$inviteTokens <= 0}
					class="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:bg-gray-600 disabled:cursor-not-allowed transition-colors"
				>
					{$inviteTokens > 0 ? 'Generate Invite' : 'No Tokens Available'}
				</button>
			</div>
		</div>

		<!-- Invites List -->
		<div class="bg-dark-800 rounded-lg border border-gray-700">
			<div class="px-6 py-4 border-b border-gray-700">
				<h3 class="text-lg font-semibold text-gray-100">Your Invites</h3>
			</div>
			
			{#if $invites.length === 0}
				<div class="px-6 py-8 text-center text-gray-400">
					<p>You haven't created any invites yet.</p>
					{#if $inviteTokens > 0}
						<p class="mt-2">Click "Generate Invite" above to get started!</p>
					{/if}
				</div>
			{:else}
				<div class="divide-y divide-gray-700">
					{#each $invites as invite}
						<div class="px-6 py-4">
							<div class="flex items-center justify-between">
								<div class="flex-1">
									<div class="flex items-center space-x-3">
										<span class="font-mono text-sm bg-dark-700 px-2 py-1 rounded text-gray-300">
											{invite.token.substring(0, 8)}...
										</span>
										<span class={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(getInviteStatus(invite))}`}>
											{getInviteStatus(invite)}
										</span>
									</div>
									
									<div class="mt-2 text-sm text-gray-400">
										<p>Created: {formatDate(invite.created_at)}</p>
										{#if invite.used_at}
											<p>Used by: {invite.invitee_username} on {formatDate(invite.used_at)}</p>
										{:else if invite.revoked_at}
											<p>Revoked on: {formatDate(invite.revoked_at)}</p>
										{:else}
											<p class="text-green-400">Active - waiting to be used</p>
										{/if}
									</div>
								</div>
								
								{#if getInviteStatus(invite) === 'Active'}
									<div class="flex space-x-2">
										<button
											on:click={() => navigator.clipboard.writeText(`${window.location.origin}/register?token=${invite.token}`)}
											class="px-3 py-1 text-sm bg-dark-700 text-gray-300 rounded hover:bg-dark-600 transition-colors"
										>
											Copy Link
										</button>
										<button
											on:click={() => handleRevokeInvite(invite.token)}
											class="px-3 py-1 text-sm bg-red-900/20 text-red-400 rounded hover:bg-red-900/30 transition-colors border border-red-800"
										>
											Revoke
										</button>
									</div>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>


	</div>
</div>
