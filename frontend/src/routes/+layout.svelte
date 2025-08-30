<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import '../app.css';
	import Navbar from '$lib/components/Layout/Navbar.svelte';
	import NotificationContainer from '$lib/components/UI/NotificationContainer.svelte';
	import { auth, isAuthenticated } from '$lib/stores/auth';

	onMount(async () => {
		// Check authentication status on app load
		await auth.checkAuth();
	});

	// Don't show navbar on login page
	$: showNavbar = !$page.url.pathname.startsWith('/login');

	// Define public routes that don't require authentication
	const publicRoutes = ['/', '/login'];
	
	// Handle authentication state changes
	$: if (!$isAuthenticated && !publicRoutes.includes($page.url.pathname)) {
		// User is on a protected page but not authenticated
		// Try to refresh the token automatically
		auth.checkAuth().then(() => {
			// If still not authenticated after refresh attempt, redirect to login
			if (!$isAuthenticated) {
				goto('/login');
			}
		});
	}
</script>

{#if showNavbar}
	<Navbar />
{/if}

<main class="min-h-screen">
	<slot />
</main>

<!-- Global notification container -->
<NotificationContainer />
