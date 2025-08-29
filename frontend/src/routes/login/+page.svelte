<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { BookOpen, Eye, EyeOff, Loader2 } from 'lucide-svelte';
	import { auth, isAuthenticated } from '$lib/stores/auth';

	let username = '';
	let password = '';
	let loading = false;
	let error = '';
	let showPassword = false;

	onMount(() => {
		// Redirect if already authenticated
		if ($isAuthenticated) {
			goto('/');
		}
	});

	async function handleSubmit() {
		if (!username || !password) {
			error = 'Please fill in all fields';
			return;
		}

		loading = true;
		error = '';

		try {
			const success = await auth.login(username, password);
			if (success) {
				goto('/');
			} else {
				error = 'Invalid username or password';
			}
		} catch (err) {
			error = 'Login failed. Please try again.';
		} finally {
			loading = false;
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			handleSubmit();
		}
	}
</script>

<svelte:head>
			<title>Sign In - marchive</title>
</svelte:head>

<div class="min-h-screen bg-dark-950 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<!-- Header -->
		<div class="text-center">
			<div class="flex justify-center">
				<BookOpen class="h-12 w-12 text-primary-500" />
			</div>
			<h2 class="mt-6 text-3xl font-bold text-gray-100">
				Welcome back
			</h2>
			<p class="mt-2 text-sm text-gray-400">
				Sign in to access your digital library
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
		<form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
			<div class="space-y-4">
				<!-- Username -->
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
						on:keydown={handleKeydown}
						class="input"
						placeholder="Enter your username"
					/>
				</div>

				<!-- Password -->
				<div>
					<label for="password" class="block text-sm font-medium text-gray-300 mb-2">
						Password
					</label>
					<div class="relative">
						{#if showPassword}
							<input
								id="password"
								name="password"
								type="text"
								autocomplete="current-password"
								required
								bind:value={password}
								on:keydown={handleKeydown}
								class="input pr-10"
								placeholder="Enter your password"
							/>
						{:else}
							<input
								id="password"
								name="password"
								type="password"
								autocomplete="current-password"
								required
								bind:value={password}
								on:keydown={handleKeydown}
								class="input pr-10"
								placeholder="Enter your password"
							/>
						{/if}
						<button
							type="button"
							class="absolute inset-y-0 right-0 pr-3 flex items-center"
							on:click={() => showPassword = !showPassword}
						>
							{#if showPassword}
								<EyeOff class="h-5 w-5 text-gray-400" />
							{:else}
								<Eye class="h-5 w-5 text-gray-400" />
							{/if}
						</button>
					</div>
				</div>
			</div>

			<!-- Error Message -->
			{#if error}
				<div class="bg-red-900/50 border border-red-800 text-red-300 px-4 py-3 rounded-lg">
					{error}
				</div>
			{/if}

			<!-- Submit Button -->
			<button
				type="submit"
				disabled={loading}
				class="w-full btn-primary flex items-center justify-center py-3"
			>
				{#if loading}
					<Loader2 class="h-5 w-5 animate-spin mr-2" />
					<span>Signing in...</span>
				{:else}
					<span>Sign in</span>
				{/if}
			</button>
		</form>

	</div>
</div>
