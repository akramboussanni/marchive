<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { Search, X, Download, Loader2, ChevronLeft, ChevronRight, Lock } from 'lucide-svelte';
	import { books, isSearching, type Book, type SearchResponse } from '$lib/stores/books';
	import { isAuthenticated } from '$lib/stores/auth';
	import { showError, showSuccess } from '$lib/stores/notifications';
	import BookCover from '../Books/BookCover.svelte';

	export let query = '';

	const dispatch = createEventDispatcher();

	let currentSearchResponse: SearchResponse | null = null;
	let searchResults: Book[] = [];
	let searchTimeout: NodeJS.Timeout;
	let downloadingBooks = new Set<string>();
	let currentPage = 0;
	let searchLimit = 20;

	onMount(() => {
		if (query) {
			performSearch();
		}
	});

	function handleClose() {
		dispatch('close');
	}

	function debounceSearch() {
		if (!$isAuthenticated) {
			showError('Authentication Required', 'Please sign in to search for books. Search is only available to authenticated users since you cannot download books without an account.');
			return;
		}
		
		clearTimeout(searchTimeout);
		currentPage = 0; // Reset to first page on new search
		searchTimeout = setTimeout(performSearch, 1500);
	}

	async function performSearch(page = 0) {
		if (!$isAuthenticated) {
			showError('Authentication Required', 'Please sign in to search for books. Search is only available to authenticated users since you cannot download books without an account.');
			return;
		}

		if (!query.trim()) {
			searchResults = [];
			currentSearchResponse = null;
			return;
		}

		try {
			const offset = page * searchLimit;
			const result = await books.search(query, searchLimit, offset);
			currentSearchResponse = result;
			searchResults = result.books;
			currentPage = page;
		} catch (error) {
			console.error('Search failed:', error);
			searchResults = [];
			currentSearchResponse = null;
		}
	}

	async function handlePrevPage() {
		if (currentPage > 0) {
			await performSearch(currentPage - 1);
		}
	}

	async function handleNextPage() {
		if (currentSearchResponse?.pagination.has_next) {
			await performSearch(currentPage + 1);
		}
	}

	async function handleDownload(book: Book, index: number) {
		if (!$isAuthenticated) {
			showError('Authentication Required', 'Please sign in to download books');
			return;
		}

		downloadingBooks.add(book.hash);
		downloadingBooks = downloadingBooks;

		try {
			// Use cached download if search ID is available, otherwise fallback to regular download
			if (currentSearchResponse?.search_id) {
				await books.requestCachedDownload(currentSearchResponse.search_id, index, book);
			} else {
				await books.requestDownload(book.hash, book.title);
			}
			showSuccess('Download Started', `"${book.title}" has been added to your downloads. Check your downloads page for progress.`);
		} catch (error) {
			console.error('Download failed:', error);
			const errorMessage = (error as Error).message || 'Unknown error';
			
			if (errorMessage.includes('already requested')) {
				showError('Already Requested', `You have already requested "${book.title}". Check your downloads page.`);
			} else if (errorMessage.includes('daily download limit')) {
				showError('Download Limit Reached', 'You have reached your daily download limit. Try again tomorrow.');
			} else if (errorMessage.includes('available for download')) {
				showSuccess('Book Ready', `"${book.title}" is already available in your downloads!`);
			} else {
				showError('Download Failed', 'The book is temporarily unavailable. Please try again later.');
			}
		} finally {
			downloadingBooks.delete(book.hash);
			downloadingBooks = downloadingBooks;
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			handleClose();
		}
	}

	$: if (query) {
		debounceSearch();
	}
</script>

<!-- Modal Backdrop -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
	class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-start justify-center pt-4 sm:pt-16 px-2 sm:px-0"
	on:click={handleClose}
>
	<!-- Modal Content -->
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div
		class="bg-dark-900 rounded-xl border border-gray-800 w-full max-w-4xl mx-auto max-h-[90vh] sm:max-h-[80vh] overflow-hidden"
		on:click|stopPropagation
	>
		<!-- Header -->
		<div class="border-b border-gray-800 p-4">
			<div class="flex items-center space-x-3">
				<Search class="h-5 w-5 text-gray-400" />
				<input
					bind:value={query}
					on:keydown={handleKeydown}
					placeholder="Search for books..."
					class="flex-1 bg-transparent text-gray-100 placeholder-gray-400 focus:outline-none text-lg"
					autofocus
				/>
				<button
					on:click={handleClose}
					class="btn-ghost p-2"
				>
					<X class="h-5 w-5" />
				</button>
			</div>
		</div>

		<!-- Results -->
		<div class="p-3 sm:p-4 max-h-[60vh] overflow-y-auto">
			{#if $isSearching}
				<div class="flex items-center justify-center py-8 sm:py-12">
					<Loader2 class="h-6 w-6 sm:h-8 sm:w-8 animate-spin text-primary-500" />
					<span class="ml-3 text-gray-400 text-sm sm:text-base">Searching...</span>
				</div>
			{:else if searchResults.length > 0}
				<div class="grid grid-cols-1 gap-3 sm:gap-4">
					{#each searchResults as book, index}
						{@const globalIndex = currentPage * searchLimit + index}
						<div class="card p-4 hover:bg-dark-800 transition-all duration-200">
							<div class="flex space-x-4">
								<!-- Book Cover -->
								<div class="flex-shrink-0">
									<BookCover {book} size="sm" />
								</div>

								<!-- Book Details -->
								<div class="flex-1 min-w-0">
									<h3 class="text-sm font-medium text-gray-100 line-clamp-2">
										{book.title}
									</h3>
									<p class="text-xs text-gray-400 mt-1">
										{book.authors}
									</p>
									<div class="flex items-center space-x-2 mt-2 text-xs text-gray-500">
										<span>{book.format}</span>
										<span>•</span>
										<span>{book.size}</span>
										{#if book.language}
											<span>•</span>
											<span>{book.language}</span>
										{/if}
									</div>


									<!-- Download Button -->
									<div class="mt-3">
										<button
											on:click={() => handleDownload(book, globalIndex)}
											disabled={downloadingBooks.has(book.hash)}
											class="btn-primary text-xs px-3 py-1 flex items-center space-x-1"
										>
											{#if downloadingBooks.has(book.hash)}
												<Loader2 class="h-3 w-3 animate-spin" />
												<span>Downloading...</span>
											{:else}
												<Download class="h-3 w-3" />
												<span>Download</span>
											{/if}
										</button>
									</div>
								</div>
							</div>
						</div>
					{/each}
				</div>

				<!-- Pagination Controls -->
				{#if currentSearchResponse && currentSearchResponse.total > searchLimit}
					<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mt-6 pt-4 border-t border-dark-700 space-y-3 sm:space-y-0">
						<div class="text-sm text-gray-400 text-center sm:text-left">
							Showing {currentPage * searchLimit + 1}-{Math.min((currentPage + 1) * searchLimit, currentSearchResponse.total)} of {currentSearchResponse.total} results
						</div>
						<div class="flex items-center justify-center space-x-2">
							<button
								on:click={handlePrevPage}
								disabled={currentPage === 0 || $isSearching}
								class="btn-secondary-sm flex items-center space-x-1 disabled:opacity-50 disabled:cursor-not-allowed"
							>
								<ChevronLeft class="h-4 w-4" />
								<span class="hidden sm:inline">Previous</span>
							</button>
							<span class="px-3 py-1 text-sm text-gray-400">
								Page {currentPage + 1} of {Math.ceil(currentSearchResponse.total / searchLimit)}
							</span>
							<button
								on:click={handleNextPage}
								disabled={!currentSearchResponse.pagination.has_next || $isSearching}
								class="btn-secondary-sm flex items-center space-x-1 disabled:opacity-50 disabled:cursor-not-allowed"
							>
								<span class="hidden sm:inline">Next</span>
								<ChevronRight class="h-4 w-4" />
							</button>
						</div>
					</div>
				{/if}
			{:else if query.trim()}
				<div class="text-center py-12">
					<Search class="h-12 w-12 text-gray-600 mx-auto mb-4" />
					<p class="text-gray-400">No books found for "{query}"</p>
					<p class="text-sm text-gray-500 mt-2">Try different keywords or check your spelling</p>
				</div>
			{:else}
				<div class="text-center py-12">
					<Search class="h-12 w-12 text-gray-600 mx-auto mb-4" />
					<p class="text-gray-400">Start typing to search for books</p>
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}

	.btn-secondary-sm {
		background-color: rgb(31 41 55);
		color: rgb(209 213 219);
		border: 1px solid rgb(75 85 99);
		font-weight: 500;
		padding: 0.25rem 0.75rem;
		border-radius: 0.5rem;
		transition: background-color 0.15s ease-in-out;
		font-size: 0.875rem;
	}
	.btn-secondary-sm:hover {
		background-color: rgb(55 65 81);
	}
	.btn-secondary-sm:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
