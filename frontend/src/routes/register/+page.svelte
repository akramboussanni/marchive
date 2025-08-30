<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { inviteStore } from '$lib/stores/invites';
	import { isAuthenticated } from '$lib/stores/auth';

	let username = '';
	let password = '';
	let confirmPassword = '';
	let isLoading = false;
	let errorMessage = '';
	let inviteToken = '';

	onMount(() => {
		// Get token from URL query parameter
		const urlParams = new URLSearchParams($page.url.search);
		inviteToken = urlParams.get('token') || '';
		
		// If already authenticated, redirect to home
		if ($isAuthenticated) {
			goto('/');
		}
	});

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

		if (password.length < 8) {
			errorMessage = 'Password must be at least 8 characters long';
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

<div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
	<div class="sm:mx-auto sm:w-full sm:max-w-md">
		<div class="text-center">
			<h2 class="text-3xl font-extrabold text-gray-900">Create your account</h2>
			<p class="mt-2 text-sm text-gray-600">
				Use your invite to join Marchive
			</p>
		</div>
	</div>

	<div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
		<div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
			{#if !inviteToken}
				<div class="text-center">
					<div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 mb-4">
						<svg class="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
						</svg>
					</div>
					<h3 class="text-lg font-medium text-gray-900 mb-2">Invalid Invite</h3>
					<p class="text-gray-600 mb-4">
						You need a valid invite token to register. Please check your invite link.
					</p>
					<button
						on:click={goToLogin}
						class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
					>
						Go to Login
					</button>
				</div>
			{:else}
				<form class="space-y-6" on:submit|preventDefault={handleSubmit}>
					<!-- Invite Token Display -->
					<div class="bg-blue-50 border border-blue-200 rounded-md p-3">
						<div class="flex">
							<div class="flex-shrink-0">
								<svg class="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
									<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
								</svg>
							</div>
							<div class="ml-3">
								<p class="text-sm text-blue-700">
									<strong>Invite Token:</strong> {inviteToken.substring(0, 8)}...
								</p>
							</div>
						</div>
					</div>

					<!-- Username -->
					<div>
						<label for="username" class="block text-sm font-medium text-gray-700">
							Username
						</label>
						<div class="mt-1">
							<input
								id="username"
								name="username"
								type="text"
								required
								bind:value={username}
								class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Enter your username"
							/>
						</div>
					</div>

					<!-- Password -->
					<div>
						<label for="password" class="block text-sm font-medium text-gray-700">
							Password
						</label>
						<div class="mt-1">
							<input
								id="password"
								name="password"
								type="password"
								required
								bind:value={password}
								class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Enter your password"
							/>
						</div>
					</div>

					<!-- Confirm Password -->
					<div>
						<label for="confirmPassword" class="block text-sm font-medium text-gray-700">
							Confirm Password
						</label>
						<div class="mt-1">
							<input
								id="confirmPassword"
								name="confirmPassword"
								type="password"
								required
								bind:value={confirmPassword}
								class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
								placeholder="Confirm your password"
							/>
						</div>
					</div>

					<!-- Error Message -->
					{#if errorMessage}
						<div class="bg-red-50 border border-red-200 rounded-md p-3">
							<div class="flex">
								<div class="flex-shrink-0">
									<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
										<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
									</svg>
								</div>
								<div class="ml-3">
									<p class="text-sm text-red-700">{errorMessage}</p>
								</div>
							</div>
						</div>
					{/if}

					<!-- Submit Button -->
					<div>
						<button
							type="submit"
							disabled={isLoading}
							class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:bg-gray-400 disabled:cursor-not-allowed"
						>
							{isLoading ? 'Creating Account...' : 'Create Account'}
						</button>
					</div>

					<!-- Login Link -->
					<div class="text-center">
						<p class="text-sm text-gray-600">
							Already have an account?
							<a href="/login" class="font-medium text-blue-600 hover:text-blue-500">
								Sign in here
							</a>
						</p>
					</div>
				</form>
			{/if}
		</div>
	</div>
</div>
