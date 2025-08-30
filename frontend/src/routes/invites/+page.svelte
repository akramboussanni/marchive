<script lang="ts">
	import { onMount } from 'svelte';
	import { invites, inviteTokens, inviteStore } from '$lib/stores/invites';
	import { isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}
		await inviteStore.loadInvites();
	});

	async function handleCreateInvite() {
		const invite = await inviteStore.createInvite();
		if (invite) {
			// Copy invite URL to clipboard
			navigator.clipboard.writeText(invite.invite_url);
		}
	}

	async function handleRevokeInvite(token: string) {
		if (confirm('Are you sure you want to revoke this invite? You will get your token back.')) {
			await inviteStore.revokeInvite(token);
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
			case 'Active': return 'text-green-600 bg-green-100';
			case 'Used': return 'text-blue-600 bg-blue-100';
			case 'Revoked': return 'text-red-600 bg-red-100';
			default: return 'text-gray-600 bg-gray-100';
		}
	}
</script>

<svelte:head>
	<title>Invites - Marchive</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<div class="max-w-4xl mx-auto">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="text-3xl font-bold text-gray-900 mb-2">Invites</h1>
			<p class="text-gray-600">
				Invite others to join Marchive. Each user gets 1 invite token.
			</p>
		</div>

		<!-- Invite Token Status -->
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
			<div class="flex items-center justify-between">
				<div>
					<h2 class="text-lg font-semibold text-gray-900">Your Invite Tokens</h2>
					<p class="text-gray-600">You have {$inviteTokens} invite token{$inviteTokens !== 1 ? 's' : ''} available</p>
				</div>
				<button
					on:click={handleCreateInvite}
					disabled={$inviteTokens <= 0}
					class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
				>
					{$inviteTokens > 0 ? 'Generate Invite' : 'No Tokens Available'}
				</button>
			</div>
		</div>

		<!-- Invites List -->
		<div class="bg-white rounded-lg shadow-sm border border-gray-200">
			<div class="px-6 py-4 border-b border-gray-200">
				<h3 class="text-lg font-semibold text-gray-900">Your Invites</h3>
			</div>
			
			{#if $invites.length === 0}
				<div class="px-6 py-8 text-center text-gray-500">
					<p>You haven't created any invites yet.</p>
					{#if $inviteTokens > 0}
						<p class="mt-2">Click "Generate Invite" above to get started!</p>
					{/if}
				</div>
			{:else}
				<div class="divide-y divide-gray-200">
					{#each $invites as invite}
						<div class="px-6 py-4">
							<div class="flex items-center justify-between">
								<div class="flex-1">
									<div class="flex items-center space-x-3">
										<span class="font-mono text-sm bg-gray-100 px-2 py-1 rounded text-gray-700">
											{invite.token.substring(0, 8)}...
										</span>
										<span class={`px-2 py-1 rounded-full text-xs font-medium ${getStatusColor(getInviteStatus(invite))}`}>
											{getInviteStatus(invite)}
										</span>
									</div>
									
									<div class="mt-2 text-sm text-gray-600">
										<p>Created: {formatDate(invite.created_at)}</p>
										{#if invite.used_at}
											<p>Used by: {invite.invitee_username} on {formatDate(invite.used_at)}</p>
										{:else if invite.revoked_at}
											<p>Revoked on: {formatDate(invite.revoked_at)}</p>
										{:else}
											<p class="text-green-600">Active - waiting to be used</p>
										{/if}
									</div>
								</div>
								
								{#if getInviteStatus(invite) === 'Active'}
									<div class="flex space-x-2">
										<button
											on:click={() => navigator.clipboard.writeText(`${window.location.origin}/register?token=${invite.token}`)}
											class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
										>
											Copy Link
										</button>
										<button
											on:click={() => handleRevokeInvite(invite.token)}
											class="px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200 transition-colors"
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

		<!-- Instructions -->
		<div class="mt-8 bg-blue-50 rounded-lg p-6">
			<h3 class="text-lg font-semibold text-blue-900 mb-3">How Invites Work</h3>
			<div class="text-blue-800 space-y-2">
				<p>• Each user starts with 1 invite token</p>
				<p>• Generate an invite to share with someone</p>
				• When they use the invite, your token is consumed</p>
				<p>• You can revoke unused invites to get your token back</p>
				<p>• New users get 1 invite token when they join</p>
			</div>
		</div>
	</div>
</div>
