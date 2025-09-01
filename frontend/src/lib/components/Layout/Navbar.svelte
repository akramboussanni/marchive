<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { Search, BookOpen, Settings, LogOut, User, Download, Lock, Menu, X, Gift, Mail } from 'lucide-svelte';
	import { auth, user, isAuthenticated, isAdmin } from '$lib/stores/auth';
	import { books } from '$lib/stores/books';
	import SearchModal from './SearchModal.svelte';
	import { onMount } from 'svelte';

	let searchQuery = '';
	let showSearchModal = false;
	let showUserMenu = false;
	let showMobileMenu = false;
	let requestCredits = 0;

	onMount(async () => {
		if ($isAuthenticated) {
			await loadRequestCredits();
		}
	});

	async function loadRequestCredits() {
		try {
			const response = await fetch('/api/auth/me/credits', {
				credentials: 'include'
			});
			if (response.ok) {
				const data = await response.json();
				requestCredits = data.request_credits;
			}
		} catch (error) {
			console.error('Failed to load request credits:', error);
		}
	}

	async function handleSearch() {
		if (searchQuery.trim()) {
			showSearchModal = true;
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			handleSearch();
		}
		if (event.key === '/' && event.ctrlKey) {
			event.preventDefault();
			const searchInput = document.getElementById('search-input') as HTMLInputElement;
			searchInput?.focus();
		}
	}

	async function handleLogout() {
		await auth.logout();
		showUserMenu = false;
	}

	async function handleDownloadRequested(event: CustomEvent) {
		// Refresh the explore page to update availability status
		// This will update the book status in the search results
		try {
			await books.explore(24, 0);
			// Also refresh credits in case they were used
			await loadRequestCredits();
		} catch (error) {
			console.error('Failed to refresh explore page:', error);
		}
	}
</script>

<nav class="bg-dark-900/95 backdrop-blur-sm border-b border-gray-800 sticky top-0 z-50">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex items-center justify-between h-16">
			<!-- Logo -->
			<div class="flex items-center">
				<a href="/" class="flex items-center space-x-2">
					<BookOpen class="h-8 w-8 text-primary-500" />
					<span class="text-xl font-bold text-gray-100 hidden sm:inline">marchive</span>
				</a>
			</div>

			<!-- Search Bar - Hidden on mobile, shown on larger screens -->
			<div class="navbar-desktop flex-1 max-w-3xl mx-8">
				<div class="relative w-full">
					<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
						{#if $isAuthenticated}
							<Search class="h-6 w-6 text-gray-400" />
						{:else}
							<Lock class="h-6 w-6 text-gray-500" />
						{/if}
					</div>
					<input
						id="search-input"
						type="text"
						bind:value={searchQuery}
						on:keydown={handleKeydown}
						placeholder={$isAuthenticated ? "Search for books, authors, or topics... (Ctrl+/ to focus)" : "Sign in to search books"}
						class="block w-full pl-12 pr-24 py-3 bg-dark-800 border border-gray-700 rounded-xl text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-lg"
						class:opacity-50={!$isAuthenticated}
						disabled={!$isAuthenticated}
					/>
					<button
						on:click={handleSearch}
						class="absolute inset-y-0 right-0 pr-4 flex items-center"
						disabled={!$isAuthenticated}
						class:opacity-50={!$isAuthenticated}
					>
						<kbd class="px-3 py-1.5 text-sm text-gray-400 bg-dark-700 rounded-lg border border-gray-600 font-mono">
							{#if $isAuthenticated}
								Enter
							{:else}
								Locked
							{/if}
						</kbd>
					</button>
				</div>
			</div>

			<!-- Mobile Search Button -->
			<div class="navbar-mobile">
				{#if $isAuthenticated}
					<button
						on:click={() => showSearchModal = true}
						class="btn-ghost p-2"
					>
						<Search class="h-5 w-5" />
					</button>
				{/if}
			</div>

			<!-- Desktop Navigation - Hidden on mobile -->
			<div class="navbar-desktop flex items-center space-x-6">
				{#if $isAuthenticated}
					<!-- Downloads -->
					<a
						href="/downloads"
						class="btn-ghost flex items-center space-x-3 px-4 py-2 rounded-lg hover:bg-dark-800 transition-colors"
						class:bg-dark-800={$page.url.pathname === '/downloads'}
					>
						<Download class="h-6 w-6" />
						<span class="text-gray-200 font-medium">Downloads</span>
					</a>

					<!-- Admin -->
					{#if $isAdmin}
						<a
							href="/admin"
							class="btn-ghost flex items-center space-x-3 px-4 py-2 rounded-lg hover:bg-dark-800 transition-colors"
							class:bg-dark-800={$page.url.pathname.startsWith('/admin')}
						>
							<Settings class="h-6 w-6" />
							<span class="text-gray-200 font-medium">Admin</span>
						</a>
					{/if}

					<!-- User Menu -->
					<div class="relative">
						<button
							on:click={() => showUserMenu = !showUserMenu}
							class="btn-ghost flex items-center space-x-3 px-4 py-2 rounded-lg hover:bg-dark-800 transition-colors"
						>
							<User class="h-6 w-6" />
							<span class="text-gray-200 font-medium">{$user?.username}</span>
						</button>

						{#if showUserMenu}
							<div class="absolute right-0 mt-2 w-48 bg-dark-800 rounded-lg shadow-lg border border-gray-700 py-1 z-50">
								<a
									href="/profile"
									class="block px-4 py-2 text-sm text-gray-300 hover:bg-dark-700"
									on:click={() => showUserMenu = false}
								>
									Profile Settings
								</a>
								<a
									href="/invites"
									class="block px-4 py-2 text-sm text-gray-300 hover:bg-dark-700"
									on:click={() => showUserMenu = false}
								>
									<Mail class="h-4 w-4 inline mr-2" />
									Invites
								</a>
								<a
									href="/redeem"
									class="block px-4 py-2 text-sm text-gray-300 hover:bg-dark-700"
									on:click={() => showUserMenu = false}
								>
									<Gift class="h-4 w-4 inline mr-2" />
									Redeem Code
								</a>
								<button
									on:click={handleLogout}
									class="w-full text-left block px-4 py-2 text-sm text-gray-300 hover:bg-dark-700"
								>
									<LogOut class="h-4 w-4 inline mr-2" />
									Sign Out
								</button>
							</div>
						{/if}
					</div>
				{:else}
					<a href="/login" class="btn-primary px-6 py-2.5 text-base font-medium">Sign In</a>
				{/if}
			</div>

			<!-- Mobile Menu Button -->
			<div class="navbar-mobile">
				<button
					on:click={() => showMobileMenu = !showMobileMenu}
					class="btn-ghost p-2"
				>
					{#if showMobileMenu}
						<X class="h-5 w-5" />
					{:else}
						<Menu class="h-5 w-5" />
					{/if}
				</button>
			</div>
		</div>

		<!-- Mobile Menu -->
		{#if showMobileMenu}
			<div class="navbar-mobile border-t border-gray-800 bg-dark-900/95 backdrop-blur-sm">
				<div class="px-2 pt-2 pb-3 space-y-1">
					{#if $isAuthenticated}
						<!-- Request Credits Display (Mobile) -->
						{#if requestCredits > 0}
							<div class="px-3 py-2">
								<div class="flex items-center justify-center space-x-2 px-3 py-2 bg-primary-500/20 border border-primary-500/30 rounded-lg">
									<Gift class="h-4 w-4 text-primary-400" />
									<span class="text-sm text-primary-300 font-medium">{requestCredits} request credits</span>
								</div>
								<div class="text-xs text-gray-500 text-center mt-1">
									Allow extra downloads beyond daily limit
								</div>
							</div>
						{/if}

						<!-- Mobile Search Bar -->
						<div class="px-3 py-2">
							<div class="relative">
								<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
									<Search class="h-5 w-5 text-gray-400" />
								</div>
								<input
									type="text"
									bind:value={searchQuery}
									on:keydown={handleKeydown}
									placeholder="Search books..."
									class="block w-full pl-10 pr-3 py-2 bg-dark-800 border border-gray-700 rounded-lg text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
								/>
							</div>
						</div>

						<!-- Mobile Navigation Links -->
						<a
							href="/downloads"
							class="block px-3 py-2 text-base font-medium text-gray-300 hover:bg-dark-800 hover:text-white rounded-md transition-colors"
							class:bg-dark-800={$page.url.pathname === '/downloads'}
							on:click={() => showMobileMenu = false}
						>
							<div class="flex items-center space-x-3">
								<Download class="h-5 w-5" />
								<span>Downloads</span>
							</div>
						</a>

						{#if $isAdmin}
							<a
								href="/admin"
								class="block px-3 py-2 text-base font-medium text-gray-300 hover:bg-dark-800 hover:text-white rounded-md transition-colors"
								class:bg-dark-800={$page.url.pathname.startsWith('/admin')}
								on:click={() => showMobileMenu = false}
							>
								<div class="flex items-center space-x-3">
									<Settings class="h-5 w-5" />
									<span>Admin</span>
								</div>
							</a>
						{/if}

						<a
							href="/profile"
							class="block px-3 py-2 text-base font-medium text-gray-300 hover:bg-dark-800 hover:text-white rounded-md transition-colors"
							on:click={() => showMobileMenu = false}
						>
							<div class="flex items-center space-x-3">
								<User class="h-5 w-5" />
								<span>Profile Settings</span>
							</div>
						</a>

						<a
							href="/invites"
							class="block px-3 py-2 text-base font-medium text-gray-300 hover:bg-dark-800 hover:text-white rounded-md transition-colors"
							on:click={() => showMobileMenu = false}
						>
							<div class="flex items-center space-x-3">
								<Mail class="h-5 w-5" />
								<span>Invites</span>
							</div>
						</a>

						<a
							href="/redeem"
							class="block px-3 py-2 text-base font-medium text-gray-300 hover:bg-dark-800 hover:text-white rounded-md transition-colors"
							on:click={() => showMobileMenu = false}
						>
							<div class="flex items-center space-x-3">
								<Gift class="h-5 w-5" />
								<span>Redeem Code</span>
							</div>
						</a>

						<button
							on:click={() => { handleLogout(); showMobileMenu = false; }}
							class="w-full text-left block px-3 py-2 text-base font-medium text-gray-300 hover:bg-dark-800 hover:text-white rounded-md transition-colors"
						>
							<div class="flex items-center space-x-3">
								<LogOut class="h-5 w-5" />
								<span>Sign Out</span>
							</div>
						</button>
					{:else}
						<!-- Mobile menu for unauthenticated users -->
						<div class="px-3 py-4">
							<div class="text-center">
								<p class="text-gray-400 text-sm mb-4">Sign in to access downloads and search</p>
								<a
									href="/login"
									class="btn-primary w-full flex items-center justify-center space-x-2"
									on:click={() => showMobileMenu = false}
								>
									<User class="h-5 w-5" />
									<span>Sign In</span>
								</a>
							</div>
						</div>
					{/if}
				</div>
			</div>
		{/if}
	</div>
</nav>

<!-- Search Modal -->
{#if showSearchModal}
	<SearchModal query={searchQuery} on:close={() => showSearchModal = false} on:downloadRequested={handleDownloadRequested} />
{/if}

<!-- Click outside to close user menu -->
{#if showUserMenu}
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div
		class="fixed inset-0 z-40"
		on:click={() => showUserMenu = false}
	></div>
{/if}
