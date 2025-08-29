<script lang="ts">
	import { Download, Loader2, Clock, CheckCircle, ExternalLink } from 'lucide-svelte';
	import type { Book } from '$lib/stores/books';
	import { books } from '$lib/stores/books';
	import { isAuthenticated } from '$lib/stores/auth';
	import { showError, showSuccess } from '$lib/stores/notifications';
	import BookCover from './BookCover.svelte';
	import FavoriteButton from '$lib/components/UI/FavoriteButton.svelte';

	export let book: Book;
	export let showDownloadButton = true;
	export let size: 'sm' | 'md' | 'lg' = 'md';

	let isDownloading = false;

	async function handleDownload() {
		if (!$isAuthenticated) {
			showError('Authentication Required', 'Please sign in to download books');
			return;
		}

		isDownloading = true;
		try {
			await books.requestDownload(book.hash, book.title);
			showSuccess('Download Started', 'Check your downloads page for progress.');
		} catch (error) {
			console.error('Download failed:', error);
			showError('Download Failed', 'You may have reached your daily limit.');
		} finally {
			isDownloading = false;
		}
	}

	function handleDownloadFile() {
		if (book.status === 'ready') {
			window.open(`/api/books/${book.hash}/download`, '_blank');
		}
	}

	$: isReady = book.status === 'ready';
	$: isProcessing = book.status === 'processing';
	$: hasError = book.status === 'error';
</script>

<div class="book-card group">
	<!-- Book Cover with Favorite Button -->
	<div class="mb-3 relative">
		<BookCover {book} {size} />
		<div class="absolute top-2 right-2">
			<FavoriteButton bookHash={book.hash} size="sm" />
		</div>
	</div>

	<!-- Book Info -->
	<div class="space-y-2">
		<h3 class="font-medium text-gray-100 text-sm leading-tight line-clamp-2 group-hover:text-white transition-colors">
			{book.title}
		</h3>
		
		{#if book.authors}
			<p class="text-xs text-gray-400 line-clamp-1">
				{book.authors}
			</p>
		{/if}

		<!-- Metadata -->
		<div class="flex items-center space-x-1 text-xs text-gray-500">
			{#if book.format}
				<span class="bg-dark-800 px-1.5 py-0.5 rounded">
					{book.format.toUpperCase()}
				</span>
			{/if}
			{#if book.size}
				<span>{book.size}</span>
			{/if}
		</div>

		<!-- Status & Download -->
		{#if showDownloadButton}
			<div class="pt-2">
				{#if isReady}
					<button
						on:click={handleDownloadFile}
						class="w-full btn-primary text-xs py-1.5 flex items-center justify-center space-x-1"
					>
						<ExternalLink class="h-3 w-3" />
						<span>Download File</span>
					</button>
				{:else if isProcessing}
					<div class="w-full btn-secondary text-xs py-1.5 flex items-center justify-center space-x-1 cursor-not-allowed">
						<Clock class="h-3 w-3" />
						<span>Processing...</span>
					</div>
				{:else if hasError}
					<div class="w-full bg-red-900/50 text-red-300 text-xs py-1.5 flex items-center justify-center space-x-1 rounded-lg border border-red-800">
						<span>❌ Error</span>
					</div>
				{:else}
					<button
						on:click={handleDownload}
						disabled={isDownloading}
						class="w-full btn-primary text-xs py-1.5 flex items-center justify-center space-x-1"
					>
						{#if isDownloading}
							<Loader2 class="h-3 w-3 animate-spin" />
							<span>Requesting...</span>
						{:else}
							<Download class="h-3 w-3" />
							<span>Download</span>
						{/if}
					</button>
				{/if}
			</div>
		{/if}

		<!-- Download Count -->
		{#if book.download_count !== undefined && book.download_count > 0}
			<p class="text-xs text-gray-500 flex items-center space-x-1">
				<CheckCircle class="h-3 w-3" />
				<span>{book.download_count} downloads</span>
			</p>
		{/if}
	</div>
</div>

<style>
	.line-clamp-1 {
		display: -webkit-box;
		-webkit-line-clamp: 1;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
	
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
