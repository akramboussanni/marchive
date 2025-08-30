<script lang="ts">
	import { onMount } from 'svelte';
	import { Loader2, RefreshCw, Heart, Star, Gift } from 'lucide-svelte';
	import type { Book } from '$lib/stores/books';
	import { books, exploredBooks, userFavorites, isFavoritesLoading } from '$lib/stores/books';
	import { isAuthenticated } from '$lib/stores/auth';
	import { localFavorites } from '$lib/utils/localFavorites';
	import BookCard from './BookCard.svelte';
	// import Pagination from '$lib/components/UI/Pagination.svelte';

	export let title = 'Public Library';
	export let showLoadMore = true;

	let loading = false;
	let loadingMore = false;
	let currentPage = 0;
	let favoritesPage = 0;
	const pageSize = 24;

	// Download status
	let downloadStatus: any = null;
	let statusLoading = false;

	onMount(() => {
		loadBooks();
		if ($isAuthenticated) {
			loadFavorites();
			loadDownloadStatus();
		}
	});

	async function loadDownloadStatus() {
		if (!$isAuthenticated) return;
		
		statusLoading = true;
		try {
			const response = await fetch('/api/books/download-status', {
				credentials: 'include'
			});
			if (response.ok) {
				downloadStatus = await response.json();
			}
		} catch (error) {
			console.error('Failed to load download status:', error);
		} finally {
			statusLoading = false;
		}
	}

	async function loadBooks() {
		loading = true;
		try {
			await books.explore(pageSize, 0);
			currentPage = 0;
		} catch (error) {
			console.error('Failed to load books:', error);
		} finally {
			loading = false;
		}
	}

	async function loadFavorites() {
		try {
			await books.getUserFavorites(pageSize, 0);
			favoritesPage = 0;
		} catch (error) {
			console.error('Failed to load favorites:', error);
		}
	}

	async function refreshAll() {
		// Refresh both favorites and public library
		await Promise.all([
			loadBooks(),
			$isAuthenticated ? loadFavorites() : Promise.resolve(),
			$isAuthenticated ? loadDownloadStatus() : Promise.resolve()
		]);
	}

	async function loadMoreBooks() {
		if (!$exploredBooks || loadingMore) return;

		loadingMore = true;
		try {
			const nextPage = currentPage + 1;
			const result = await books.explore(pageSize, nextPage * pageSize);
			
			// Append new books to existing ones
			exploredBooks.update(current => {
				if (!current) return result;
				return {
					...result,
					books: [...current.books, ...result.books]
				};
			});
			
			currentPage = nextPage;
		} catch (error) {
			console.error('Failed to load more books:', error);
		} finally {
			loadingMore = false;
		}
	}

	async function loadMoreFavorites() {
		if (!$userFavorites || isFavoritesLoading) return;

		try {
			const nextPage = favoritesPage + 1;
			const result = await books.getUserFavorites(pageSize, nextPage * pageSize);
			
			// Append new favorites to existing ones
			userFavorites.update(current => {
				if (!current) return result;
				return {
					...result,
					books: [...current.books, ...result.books]
				};
			});
			
			favoritesPage = nextPage;
		} catch (error) {
			console.error('Failed to load more favorites:', error);
		}
	}

	function handlePageChange(page: number) {
		currentPage = page;
		loadBooks();
	}

	function handleFavoritesPageChange(page: number) {
		favoritesPage = page;
		loadFavorites();
	}

	$: hasMoreBooks = $exploredBooks?.pagination.has_next ?? false;
	$: displayBooks = $exploredBooks?.books ?? [];
	$: hasMoreFavorites = $userFavorites?.pagination.has_next ?? false;
	$: displayFavorites = $userFavorites?.books ?? [];
	$: localFavoritesCount = $localFavorites?.length ?? 0;
	$: totalPages = Math.ceil(($exploredBooks?.pagination.total ?? 0) / pageSize);
	$: favoritesTotalPages = Math.ceil(($userFavorites?.pagination.total ?? 0) / pageSize);

	// Watch for auth changes to load favorites
	$: if ($isAuthenticated && !$userFavorites) {
		loadFavorites();
	}
</script>

<div class="space-y-6">

	<!-- Download Status Section -->
	{#if $isAuthenticated && downloadStatus && !statusLoading}
		<div class="bg-dark-800/50 rounded-lg p-4 sm:p-6 border border-gray-700">
			<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-4 space-y-3 sm:space-y-0">
				<div class="flex items-center space-x-2">
					<Gift class="h-5 w-5 text-primary-400" />
					<h2 class="text-lg font-semibold text-gray-100">Your Download Status</h2>
				</div>
				<div class="text-xs text-gray-500">
					Daily limit: {downloadStatus.daily_limit} • Credits allow extra downloads
				</div>
			</div>
			
			<div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
				<!-- Daily Downloads -->
				<div class="text-center p-3 bg-dark-700 rounded-lg">
					<div class="text-xl font-bold text-gray-100">{downloadStatus.downloads_used}</div>
					<div class="text-sm text-gray-400">Used Today</div>
					<div class="text-xs text-gray-500">of {downloadStatus.daily_limit}</div>
				</div>
				
				<!-- Remaining Downloads -->
				<div class="text-center p-3 bg-dark-700 rounded-lg">
					<div class="text-xl font-bold text-green-400">{downloadStatus.downloads_remaining}</div>
					<div class="text-sm text-gray-400">Remaining</div>
				</div>
				
				<!-- Request Credits -->
				<div class="text-center p-3 bg-dark-700 rounded-lg">
					<div class="text-xl font-bold text-primary-400 flex items-center justify-center space-x-2">
						<Gift class="h-5 w-5" />
						<span>{downloadStatus.request_credits}</span>
					</div>
					<div class="text-sm text-gray-400">Request Credits</div>
					<div class="text-xs text-gray-500">for extra downloads</div>
				</div>
				
				<!-- Next Reset -->
				<div class="text-center p-3 bg-dark-700 rounded-lg">
					<div class="text-lg font-bold text-yellow-400">
						{new Date(downloadStatus.next_reset).toLocaleTimeString('en-US', { 
							hour: '2-digit', 
							minute: '2-digit' 
						})}
					</div>
					<div class="text-sm text-gray-400">Next Reset</div>
					<div class="text-xs text-gray-500">Daily limit resets</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Favorites Section -->
	{#if $isAuthenticated || localFavoritesCount > 0}
		<div class="bg-dark-800/50 rounded-lg p-4 sm:p-6 border border-gray-700">
			<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-4 space-y-3 sm:space-y-0">
				<div class="flex items-center space-x-2">
					<Heart class="h-6 w-6 text-red-400" />
					<h2 class="text-xl font-semibold text-gray-100">
						{#if $isAuthenticated}
							Your Favorites
						{:else}
							Local Favorites
						{/if}
					</h2>
					{#if localFavoritesCount > 0}
						<span class="bg-red-500/20 text-red-400 text-xs px-2 py-1 rounded-full">
							{localFavoritesCount}
						</span>
					{/if}
				</div>
				

			</div>

			{#if $isAuthenticated && $isFavoritesLoading && !loading}
				<div class="flex items-center justify-center py-8">
					<Loader2 class="h-6 w-6 animate-spin text-primary-500" />
					<span class="ml-3 text-gray-400">Loading favorites...</span>
				</div>
			{:else if ($isAuthenticated && displayFavorites.length > 0) || (!$isAuthenticated && localFavoritesCount > 0)}
				<div class="book-shelf mb-4">
					{#if $isAuthenticated}
						{#each displayFavorites as book (book.hash)}
							<BookCard {book} />
						{/each}
					{:else}
						{#each $localFavorites as favorite (favorite.bookHash)}
							<!-- For local favorites, we need to find the book data from the main library -->
							{#if $exploredBooks?.books}
								{@const book = $exploredBooks.books.find(b => b.hash === favorite.bookHash)}
								{#if book}
									<BookCard {book} />
								{/if}
							{/if}
						{/each}
					{/if}
				</div>

				{#if $isAuthenticated && hasMoreFavorites}
					<div class="flex justify-center">
						<button
							on:click={loadMoreFavorites}
							class="btn-secondary flex items-center justify-center space-x-2 w-full sm:w-auto max-w-xs"
						>
							<Loader2 class="h-4 w-4" />
							<span>Load More Favorites</span>
						</button>
					</div>
				{/if}
			{:else}
				<div class="text-center py-8">
					<div class="text-4xl mb-3">💖</div>
					<h3 class="text-lg font-medium text-gray-300 mb-2">
						{#if $isAuthenticated}
							No favorites yet
						{:else}
							No local favorites
						{/if}
					</h3>
					<p class="text-gray-500">
						{#if $isAuthenticated}
							Start adding books to your favorites to see them here
						{:else}
							Click the heart icon on any book to add it to your local favorites
						{/if}
					</p>
				</div>
			{/if}
		</div>
	{/if}

	<!-- Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
		<div>
			<h1 class="text-2xl font-bold text-gray-100">{title}</h1>
			{#if $exploredBooks}
				<p class="text-sm text-gray-400 mt-1">
					{$exploredBooks.pagination.total} books available
				</p>
			{/if}
		</div>
		
		<button
			on:click={refreshAll}
			disabled={loading || ($isAuthenticated && $isFavoritesLoading)}
			class="btn-secondary flex items-center justify-center sm:justify-start space-x-2 w-full sm:w-auto"
		>
			<RefreshCw class="h-4 w-4 {(loading || ($isAuthenticated && $isFavoritesLoading)) ? 'animate-spin' : ''}" />
			<span>Refresh All</span>
		</button>
	</div>

	<!-- Loading State -->
	{#if loading || ($isAuthenticated && $isFavoritesLoading)}
		<div class="flex items-center justify-center py-16">
			<Loader2 class="h-8 w-8 animate-spin text-primary-500" />
			<span class="ml-3 text-gray-400">
				{#if loading && $isAuthenticated && $isFavoritesLoading}
					Refreshing everything...
				{:else if loading}
					Loading library...
				{:else if $isAuthenticated && $isFavoritesLoading}
					Loading favorites...
				{/if}
			</span>
		</div>
	
	<!-- Books Grid -->
	{:else if displayBooks.length > 0}
		<div class="book-shelf">
			{#each displayBooks as book (book.hash)}
				<BookCard {book} />
			{/each}
		</div>

		<!-- Load More Button -->
		{#if showLoadMore && hasMoreBooks}
			<div class="flex justify-center pt-8">
				<button
					on:click={loadMoreBooks}
					disabled={loadingMore}
					class="btn-secondary flex items-center justify-center space-x-2 w-full sm:w-auto max-w-xs"
				>
					{#if loadingMore}
						<Loader2 class="h-4 w-4 animate-spin" />
						<span>Loading more...</span>
					{:else}
						<span>Load More Books</span>
					{/if}
				</button>
			</div>
		{/if}

	<!-- Empty State -->
	{:else}
		<div class="text-center py-16">
			<div class="text-6xl mb-4">📚</div>
			<h3 class="text-lg font-medium text-gray-300 mb-2">No books available yet</h3>
			<p class="text-gray-500 mb-6">Books will appear here as users download and share them with the community</p>
			<button
				on:click={() => document.getElementById('search-input')?.focus()}
				class="btn-primary"
			>
				Search for Books
			</button>
		</div>
	{/if}
</div>

<style>
	.book-shelf {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 1.5rem;
	}

	.btn-primary {
		background-color: rgb(79 70 229);
		color: white;
		font-weight: 500;
		padding: 0.5rem 1rem;
		border-radius: 0.5rem;
		transition: background-color 0.15s ease-in-out;
		border: none;
		cursor: pointer;
	}
	.btn-primary:hover {
		background-color: rgb(67 56 202);
	}

	.btn-secondary {
		background-color: rgb(31 41 55);
		color: rgb(209 213 219);
		border: 1px solid rgb(75 85 99);
		font-weight: 500;
		padding: 0.5rem 1rem;
		border-radius: 0.5rem;
		transition: background-color 0.15s ease-in-out;
		cursor: pointer;
	}
	.btn-secondary:hover {
		background-color: rgb(55 65 81);
	}
	.btn-secondary:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
