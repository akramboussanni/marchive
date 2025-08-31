<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { inviteStore } from '$lib/stores/invites';
	import { isAuthenticated } from '$lib/stores/auth';
	import PasswordInput from '$lib/components/UI/PasswordInput.svelte';

	let username = '';
	let password = '';
	let confirmPassword = '';
	let isLoading = false;
	let errorMessage = '';
	let inviteToken = '';
	
	// Password validation state from components
	let passwordValidation: any = null;
	let passwordsMatch = false;

	onMount(() => {
		// Get token from URL query parameter
		const urlParams = new URLSearchParams($page.url.search);
		inviteToken = urlParams.get('token') || '';
	});

	// Reactive statement to check if passwords match
	$: passwordsMatch = Boolean(password && confirmPassword && password === confirmPassword);

	async function handleSubmit() {
		// Reset error message
		errorMessage = '';

		// Validation
		if (!username.trim()) {
			errorMessage = 'Username is required';
			return;
		}

		if (username.length < 3) {
			errorMessage = 'Username must be at least 3 characters long';
			return;
		}

		if (username.length > 20) {
			errorMessage = 'Username must be less than 20 characters long';
			return;
		}

		if (!password) {
			errorMessage = 'Password is required';
			return;
		}

		if (!passwordValidation?.isValid) {
			errorMessage = 'Password does not meet requirements';
			return;
		}

		if (password !== confirmPassword) {
			errorMessage = 'Passwords do not match';
			return;
		}

		if (!inviteToken) {
			errorMessage = 'Invite token is required';
			return;
		}

		isLoading = true;

		try {
			const success = await inviteStore.useInvite(inviteToken, username, password);
			if (success) {
				// Redirect to home page after successful registration
				goto('/');
			}
		} catch (error) {
			console.error('Registration error:', error);
			errorMessage = 'An error occurred during registration';
		} finally {
			isLoading = false;
		}
	}

	function goToLogin() {
		goto('/login');
	}
</script>

<svelte:head>
	<title>Register - Marchive</title>
</svelte:head>

<div class="min-h-screen bg-dark-950 flex items-center justify-center py-8 sm:py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-6 sm:space-y-8">
		<!-- Header -->
		<div class="text-center">
			<div class="flex justify-center">
				<div class="h-10 w-10 sm:h-12 sm:w-12 bg-primary-900/20 rounded-full flex items-center justify-center border border-primary-800">
					<svg class="h-6 w-6 sm:h-7 sm:w-7 text-primary-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
					</svg>
				</div>
			</div>
			<h2 class="mt-4 sm:mt-6 text-2xl sm:text-3xl font-bold text-gray-100">
				Create your account
			</h2>
			<p class="mt-2 text-sm text-gray-400">
				Use your invite to join Marchive
			</p>
		</div>

		<!-- Back to Home Button -->
		<div class="text-center">
			<a href="/" class="inline-flex items-center text-sm text-gray-400 hover:text-gray-300 transition-colors">
				<svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
				</svg>
				Back to Home
			</a>
		</div>

		<!-- Form -->
		<div class="bg-dark-800 rounded-lg border border-gray-700 p-6 sm:p-8">
			{#if !inviteToken}
				<div class="text-center">
					<div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-900/20 mb-4 border border-red-800">
						<svg class="h-6 w-6 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
						</svg>
					</div>
					<h3 class="text-lg font-medium text-gray-300 mb-2">Invalid Invite</h3>
					<p class="text-gray-400 mb-4">
						You need a valid invite token to register. Please check your invite link.
					</p>
					<button
						on:click={goToLogin}
						class="w-full px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
					>
						Go to Login
					</button>
				</div>
			{:else}
				<form on:submit|preventDefault={handleSubmit} class="space-y-6">
					{#if errorMessage}
						<div class="bg-red-900/20 border border-red-800 rounded-lg p-4">
							<p class="text-red-400 text-sm">{errorMessage}</p>
						</div>
					{/if}

					<div>
						<label for="username" class="block text-sm font-medium text-gray-300 mb-2">
							Username
						</label>
						<input
							id="username"
							name="username"
							type="text"
							autocomplete="username"
							required
							bind:value={username}
							class="w-full px-4 py-3 border border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 bg-dark-700 text-gray-100 placeholder-gray-400"
							placeholder="Enter your username"
							disabled={isLoading}
						/>
					</div>

					<PasswordInput
						id="password"
						label="Password"
						placeholder="Enter your password"
						bind:value={password}
						disabled={isLoading}
						required
						showRequirements={true}
						showStrength={true}
						bind:passwordValidation
					/>

					<PasswordInput
						id="confirmPassword"
						label="Confirm Password"
						placeholder="Confirm your password"
						bind:value={confirmPassword}
						disabled={isLoading}
						required
						showRequirements={false}
						showStrength={false}
						showMatchIndicator={true}
						confirmPassword={password}
					/>

					<button
						type="submit"
						disabled={isLoading || !passwordValidation?.isValid || !passwordsMatch}
						class="w-full px-4 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:bg-gray-600 disabled:cursor-not-allowed transition-colors font-medium"
					>
						{isLoading ? 'Creating Account...' : 'Create Account'}
					</button>
				</form>

				<div class="mt-6 text-center">
					<p class="text-sm text-gray-400">
						Already have an account? 
						<a href="/login" class="text-primary-400 hover:text-primary-300 transition-colors">
							Sign in here
						</a>
					</p>
				</div>
			{/if}
		</div>
	</div>
</div>
