<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { Search, BookOpen, Settings, LogOut, User, Download, Lock } from 'lucide-svelte';
	import { auth, user, isAuthenticated, isAdmin } from '$lib/stores/auth';
	import { books } from '$lib/stores/books';
	import SearchModal from './SearchModal.svelte';

	let searchQuery = '';
	let showSearchModal = false;
	let showUserMenu = false;

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
</script>

<nav class="bg-dark-900/95 backdrop-blur-sm border-b border-gray-800 sticky top-0 z-50">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex items-center justify-between h-16">
			<!-- Logo -->
			<div class="flex items-center">
				<a href="/" class="flex items-center space-x-2">
					<BookOpen class="h-8 w-8 text-primary-500" />
					<span class="text-xl font-bold text-gray-100">marchive</span>
				</a>
			</div>

			<!-- Search Bar -->
			<div class="flex-1 max-w-2xl mx-8">
				<div class="relative">
					<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
						{#if $isAuthenticated}
							<Search class="h-5 w-5 text-gray-400" />
						{:else}
							<Lock class="h-5 w-5 text-gray-500" />
						{/if}
					</div>
					<input
						id="search-input"
						type="text"
						bind:value={searchQuery}
						on:keydown={handleKeydown}
						placeholder={$isAuthenticated ? "Search books... (Ctrl+/ to focus)" : "Sign in to search books"}
						class="block w-full pl-10 pr-12 py-2 bg-dark-800 border border-gray-700 rounded-lg text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
						class:opacity-50={!$isAuthenticated}
						disabled={!$isAuthenticated}
					/>
					<button
						on:click={handleSearch}
						class="absolute inset-y-0 right-0 pr-3 flex items-center"
						disabled={!$isAuthenticated}
						class:opacity-50={!$isAuthenticated}
					>
						<kbd class="px-2 py-1 text-xs text-gray-400 bg-dark-700 rounded border border-gray-600">
							{#if $isAuthenticated}
								Enter
							{:else}
								Locked
							{/if}
						</kbd>
					</button>
				</div>
			</div>

			<!-- Navigation -->
			<div class="flex items-center space-x-4">
				{#if $isAuthenticated}
					<!-- Downloads -->
					<a
						href="/downloads"
						class="btn-ghost flex items-center space-x-2"
						class:bg-dark-800={$page.url.pathname === '/downloads'}
					>
						<Download class="h-5 w-5" />
						<span class="hidden sm:inline">Downloads</span>
					</a>

					<!-- Admin -->
					{#if $isAdmin}
						<a
							href="/admin"
							class="btn-ghost flex items-center space-x-2"
							class:bg-dark-800={$page.url.pathname.startsWith('/admin')}
						>
							<Settings class="h-5 w-5" />
							<span class="hidden sm:inline">Admin</span>
						</a>
					{/if}

					<!-- User Menu -->
					<div class="relative">
						<button
							on:click={() => showUserMenu = !showUserMenu}
							class="btn-ghost flex items-center space-x-2"
						>
							<User class="h-5 w-5" />
							<span class="hidden sm:inline">{$user?.username}</span>
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
					<a href="/login" class="btn-primary">Sign In</a>
				{/if}
			</div>
		</div>
	</div>
</nav>

<!-- Search Modal -->
{#if showSearchModal}
	<SearchModal
		bind:query={searchQuery}
		on:close={() => showSearchModal = false}
	/>
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
