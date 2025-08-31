<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Gift, CheckCircle, XCircle, AlertCircle } from 'lucide-svelte';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import { 
		redemptionCodeStore, 
		redeemCodeForm, 
		redemptionResult 
	} from '$lib/stores/redemptionCodes';
	import { showSuccess, showError } from '$lib/stores/notifications';

	let loading = false;

	onMount(async () => {
		if (!$isAuthenticated) {
			// Try to refresh the token before redirecting to login
			const authSuccess = await auth.checkAuthWithRefresh();
			if (!authSuccess) {
				goto('/login');
				return;
			}
		}
	});

	async function handleRedeemCode() {
		if (!$redeemCodeForm.code.trim()) {
			showError('Error', 'Please enter a redemption code');
			return;
		}

		loading = true;
		redemptionResult.set(null);
		
		try {
			const result = await redemptionCodeStore.redeemCode($redeemCodeForm);
			redemptionResult.set(result);
			
			if (result.success) {
				showSuccess('Success', result.message);
				// Clear the form on success
				redeemCodeForm.set({ code: '' });
			} else {
				showError('Error', result.message);
			}
		} catch (error) {
			console.error('Failed to redeem code:', error);
			showError('Error', 'Failed to redeem code. Please try again.');
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

<div class="min-h-screen bg-dark-950 text-gray-100">
	<div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Header -->
		<div class="text-center mb-8">
			<div class="mx-auto w-16 h-16 bg-primary-900/20 rounded-full flex items-center justify-center mb-4 border border-primary-800">
				<Gift class="w-8 h-8 text-primary-400" />
			</div>
			<h1 class="text-3xl font-bold text-gray-100">Redeem Code</h1>
			<p class="text-gray-400 mt-2">Enter a redemption code to receive rewards</p>
		</div>

		<!-- Redemption Form -->
		<div class="bg-dark-800 rounded-lg border border-gray-700 p-6 mb-6">
			<form on:submit|preventDefault={handleRedeemCode} class="space-y-4">
				<div>
					<label for="code" class="block text-sm font-medium text-gray-300 mb-2">
						Redemption Code
					</label>
					<input
						id="code"
						type="text"
						bind:value={$redeemCodeForm.code}
						placeholder="Enter your redemption code here..."
						class="w-full px-4 py-3 border border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-lg font-mono bg-dark-700 text-gray-100"
						required
						disabled={loading}
					/>
				</div>

				<button
					type="submit"
					disabled={loading || !$redeemCodeForm.code.trim()}
					class="w-full px-4 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:bg-gray-600 disabled:cursor-not-allowed transition-colors font-medium text-lg"
				>
					{loading ? 'Redeeming...' : 'Redeem Code'}
				</button>
			</form>
		</div>

	<!-- Redemption Result -->
	{#if $redemptionResult}
		<div class="bg-dark-800 rounded-lg border border-gray-700 p-6 mb-6">
			{#if $redemptionResult.success}
				<div class="text-center">
					<div class="mx-auto w-16 h-16 bg-green-900/20 rounded-full flex items-center justify-center mb-4 border border-green-800">
						<CheckCircle class="w-8 h-8 text-green-400" />
					</div>
					<h3 class="text-xl font-semibold text-green-300 mb-2">Code Redeemed Successfully!</h3>
					<p class="text-green-200 mb-4">{$redemptionResult.message}</p>
					
					<div class="bg-green-900/20 rounded-lg p-4 mb-4 border border-green-800">
						<h4 class="font-medium text-green-300 mb-2">Rewards Received:</h4>
						<div class="space-y-2">
							{#if $redemptionResult.invite_tokens_granted}
								<div class="flex items-center justify-between">
									<span class="text-green-200">Invite Tokens:</span>
									<span class="font-semibold text-green-300">+{$redemptionResult.invite_tokens_granted}</span>
								</div>
							{/if}
							{#if $redemptionResult.request_credits_granted}
								<div class="flex items-center justify-between">
									<span class="text-green-200">Request Credits:</span>
									<span class="font-semibold text-green-300">+{$redemptionResult.request_credits_granted}</span>
								</div>
							{/if}
						</div>
					</div>

					{#if $redemptionResult.new_invite_tokens_total !== undefined || $redemptionResult.new_request_credits_total !== undefined}
						<div class="bg-primary-900/20 rounded-lg p-4 mb-4 border border-primary-800">
							<h4 class="font-medium text-primary-300 mb-2">New Totals:</h4>
							<div class="space-y-2">
								{#if $redemptionResult.new_invite_tokens_total !== undefined}
									<div class="flex items-center justify-between">
										<span class="text-primary-200">Total Invite Tokens:</span>
										<span class="font-semibold text-primary-300">{$redemptionResult.new_invite_tokens_total}</span>
									</div>
								{/if}
								{#if $redemptionResult.new_request_credits_total !== undefined}
									<div class="flex items-center justify-between">
										<span class="text-primary-200">Total Request Credits:</span>
										<span class="font-semibold text-primary-300">{$redemptionResult.new_request_credits_total}</span>
									</div>
								{/if}
							</div>
						</div>
					{/if}

					<button
						on:click={resetForm}
						class="px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
					>
						Redeem Another Code
					</button>
				</div>
			{:else}
				<div class="text-center">
					<div class="mx-auto w-16 h-16 bg-red-900/20 rounded-full flex items-center justify-center mb-4 border border-red-800">
						<XCircle class="w-8 h-8 text-red-400" />
					</div>
					<h3 class="text-xl font-semibold text-red-300 mb-2">Redemption Failed</h3>
					<p class="text-red-200 mb-4">{$redemptionResult.message}</p>
					
					<button
						on:click={resetForm}
						class="px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
					>
						Try Again
					</button>
				</div>
			{/if}
		</div>
	{/if}


	</div>
</div>
