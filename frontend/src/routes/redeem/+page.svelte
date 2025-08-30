<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Gift, CheckCircle, XCircle, AlertCircle } from 'lucide-svelte';
	import { isAuthenticated } from '$lib/stores/auth';
	import { 
		redemptionCodeStore, 
		redeemCodeForm, 
		redemptionResult 
	} from '$lib/stores/redemptionCodes';
	import { notifications } from '$lib/stores/notifications';

	let loading = false;

	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}
	});

	async function handleRedeemCode() {
		if (!$redeemCodeForm.code.trim()) {
			notifications.add('error', 'Please enter a redemption code');
			return;
		}

		loading = true;
		redemptionResult.set(null);
		
		try {
			const result = await redemptionCodeStore.redeemCode($redeemCodeForm);
			redemptionResult.set(result);
			
			if (result.success) {
				notifications.add('success', result.message);
				// Clear the form on success
				redeemCodeForm.set({ code: '' });
			} else {
				notifications.add('error', result.message);
			}
		} catch (error) {
			console.error('Failed to redeem code:', error);
			notifications.add('error', 'Failed to redeem code. Please try again.');
		} finally {
			loading = false;
		}
	}

	function resetForm() {
		redeemCodeForm.set({ code: '' });
		redemptionResult.set(null);
	}
</script>

<svelte:head>
	<title>Redeem Code - Marchive</title>
</svelte:head>

<div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
	<!-- Header -->
	<div class="text-center mb-8">
		<div class="mx-auto w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mb-4">
			<Gift class="w-8 h-8 text-blue-600" />
		</div>
		<h1 class="text-3xl font-bold text-gray-900">Redeem Code</h1>
		<p class="text-gray-600 mt-2">Enter a redemption code to receive rewards</p>
	</div>

	<!-- Redemption Form -->
	<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
		<form on:submit|preventDefault={handleRedeemCode} class="space-y-4">
			<div>
				<label for="code" class="block text-sm font-medium text-gray-700 mb-2">
					Redemption Code
				</label>
				<input
					id="code"
					type="text"
					bind:value={$redeemCodeForm.code}
					placeholder="Enter your redemption code here..."
					class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-lg font-mono"
					required
					disabled={loading}
				/>
			</div>

			<button
				type="submit"
				disabled={loading || !$redeemCodeForm.code.trim()}
				class="w-full px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium text-lg"
			>
				{loading ? 'Redeeming...' : 'Redeem Code'}
			</button>
		</form>
	</div>

	<!-- Redemption Result -->
	{#if $redemptionResult}
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
			{#if $redemptionResult.success}
				<div class="text-center">
					<div class="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4">
						<CheckCircle class="w-8 h-8 text-green-600" />
					</div>
					<h3 class="text-xl font-semibold text-green-900 mb-2">Code Redeemed Successfully!</h3>
					<p class="text-green-700 mb-4">{$redemptionResult.message}</p>
					
					<div class="bg-green-50 rounded-lg p-4 mb-4">
						<h4 class="font-medium text-green-900 mb-2">Rewards Received:</h4>
						<div class="space-y-2">
							{#if $redemptionResult.invite_tokens_granted}
								<div class="flex items-center justify-between">
									<span class="text-green-700">Invite Tokens:</span>
									<span class="font-semibold text-green-900">+{$redemptionResult.invite_tokens_granted}</span>
								</div>
							{/if}
							{#if $redemptionResult.request_credits_granted}
								<div class="flex items-center justify-between">
									<span class="text-green-700">Request Credits:</span>
									<span class="font-semibold text-green-900">+{$redemptionResult.request_credits_granted}</span>
								</div>
							{/if}
						</div>
					</div>

					{#if $redemptionResult.new_invite_tokens_total !== undefined || $redemptionResult.new_request_credits_total !== undefined}
						<div class="bg-blue-50 rounded-lg p-4 mb-4">
							<h4 class="font-medium text-blue-900 mb-2">New Totals:</h4>
							<div class="space-y-2">
								{#if $redemptionResult.new_invite_tokens_total !== undefined}
									<div class="flex items-center justify-between">
										<span class="text-blue-700">Total Invite Tokens:</span>
										<span class="font-semibold text-blue-900">{$redemptionResult.new_invite_tokens_total}</span>
									</div>
								{/if}
								{#if $redemptionResult.new_request_credits_total !== undefined}
									<div class="flex items-center justify-between">
										<span class="text-blue-700">Total Request Credits:</span>
										<span class="font-semibold text-blue-900">{$redemptionResult.new_request_credits_total}</span>
									</div>
								{/if}
							</div>
						</div>
					{/if}

					<button
						on:click={resetForm}
						class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
					>
						Redeem Another Code
					</button>
				</div>
			{:else}
				<div class="text-center">
					<div class="mx-auto w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mb-4">
						<XCircle class="w-8 h-8 text-red-600" />
					</div>
					<h3 class="text-xl font-semibold text-red-900 mb-2">Redemption Failed</h3>
					<p class="text-red-700 mb-4">{$redemptionResult.message}</p>
					
					<button
						on:click={resetForm}
						class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
					>
						Try Again
					</button>
				</div>
			{/if}
		</div>
	{/if}

	<!-- Information Section -->
	<div class="bg-blue-50 rounded-lg p-6">
		<h3 class="text-lg font-semibold text-blue-900 mb-3 flex items-center">
			<AlertCircle class="w-5 h-5 mr-2" />
			How Redemption Codes Work
		</h3>
		<div class="text-blue-800 space-y-2 text-sm">
			<p>• Redemption codes can give you invite tokens and/or request credits</p>
			<p>• Each code can only be redeemed once per user</p>
			<p>• Codes may have expiration dates or usage limits</p>
			<p>• Invite tokens let you invite new users to join Marchive</p>
			<p>• Request credits let you download books beyond your daily limit</p>
		</div>
	</div>
</div>
