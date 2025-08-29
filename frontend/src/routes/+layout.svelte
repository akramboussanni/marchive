<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import '../app.css';
	import Navbar from '$lib/components/Layout/Navbar.svelte';
	import NotificationContainer from '$lib/components/UI/NotificationContainer.svelte';
	import { auth } from '$lib/stores/auth';

	onMount(() => {
		// Check authentication status on app load
		auth.checkAuth();
	});

	// Don't show navbar on login page
	$: showNavbar = !$page.url.pathname.startsWith('/login');
</script>

{#if showNavbar}
	<Navbar />
{/if}

<main class="min-h-screen">
	<slot />
</main>

<!-- Global notification container -->
<NotificationContainer />
