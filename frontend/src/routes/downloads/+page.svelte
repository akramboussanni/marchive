<script lang="ts">
	import { onMount } from 'svelte';
	import { Download, Clock, CheckCircle, XCircle, AlertCircle, Loader2, RefreshCw, Gift } from 'lucide-svelte';
	import { books, type DownloadJob } from '$lib/stores/books';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import Pagination from '$lib/components/UI/Pagination.svelte';
	import { showError, showSuccess } from '$lib/stores/notifications';

	let downloads: DownloadJob[] = [];
	let isLoading = false;
	let downloadingHashes = new Set<string>(); // Track which downloads are in progress
	let currentPage = 0;
	let totalPages = 1;
	let totalDownloads = 0;
	const limit = 20;

	// Download status
	let downloadStatus: any = null;
	let statusLoading = false;

	onMount(async () => {
		if (!$isAuthenticated) {
			// Try to refresh the token before redirecting to login
			const authSuccess = await auth.checkAuthWithRefresh();
			if (!authSuccess) {
				goto('/login');
				return;
			}
		}
		loadDownloads();
		loadDownloadStatus();
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

	async function loadDownloads() {
		if (!$isAuthenticated) return;
		
		isLoading = true;
		try {
			const offset = currentPage * limit;
			const response = await books.getUserDownloads(limit, offset);
			
			downloads = response.jobs;
			totalDownloads = response.pagination.total;
			totalPages = Math.max(1, Math.ceil(totalDownloads / limit));
		} catch (error) {
			console.error('Failed to load downloads:', error);
			showError('Failed to Load Downloads', 'Unable to load your download history. Please try again.');
		} finally {
			isLoading = false;
		}
	}

	async function handlePageChange(page: number) {
		currentPage = page;
		await loadDownloads();
	}

	async function refreshDownloads() {
		await loadDownloads();
		showSuccess('Downloads Refreshed', 'Your download list has been updated.');
	}

	function getStatusIcon(status: string) {
		switch (status.toLowerCase()) {
			case 'completed':
			case 'ready':
				return CheckCircle;
			case 'failed':
			case 'error':
				return XCircle;
			case 'processing':
			case 'downloading':
				return Loader2;
			case 'pending':
			case 'queued':
				return Clock;
			default:
				return AlertCircle;
		}
	}

	function getStatusColor(status: string) {
		switch (status.toLowerCase()) {
			case 'completed':
			case 'ready':
				return 'text-green-500';
			case 'failed':
			case 'error':
				return 'text-red-500';
			case 'processing':
			case 'downloading':
				return 'text-blue-500';
			case 'pending':
			case 'queued':
				return 'text-yellow-500';
			default:
				return 'text-gray-500';
		}
	}

	function getStatusText(status: string) {
		switch (status.toLowerCase()) {
			case 'completed':
			case 'ready':
				return 'Ready for Download';
			case 'failed':
			case 'error':
				return 'Download Failed';
			case 'processing':
			case 'downloading':
				return 'Processing';
			case 'pending':
			case 'queued':
				return 'Queued';
			default:
				return status;
		}
	}

	function formatDate(timestamp: string | number | undefined) {
		if (!timestamp) return 'Unknown';
		
		const date = new Date(typeof timestamp === 'string' ? parseInt(timestamp) * 1000 : timestamp * 1000);
		return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
	}

	async function handleDownload(hash: string, title: string) {
		if (downloadingHashes.has(hash)) return; // Prevent multiple downloads
		
		downloadingHashes.add(hash);
		
		try {
			// Use the books store download method for proper authentication
			const blob = await books.downloadBook(hash, title);
			
			// Create a download link
			const url = window.URL.createObjectURL(blob);
			const link = document.createElement('a');
			link.href = url;
			link.download = `${title || 'book'}.pdf`; // Default to PDF, adjust as needed
			
			// Trigger download
			document.body.appendChild(link);
			link.click();
			document.body.removeChild(link);
			
			// Clean up
			window.URL.revokeObjectURL(url);
			
			showSuccess('Download Started', 'Your download has begun. Check your downloads folder.');
		} catch (error) {
			console.error('Download error:', error);
			showError('Download Failed', 'Unable to start download. Please try again.');
		} finally {
			downloadingHashes.delete(hash);
		}
	}

	function getDownloadUrl(hash: string) {
		return `/api/books/${hash}/download`;
	}
</script>

<svelte:head>
	<title>Downloads - marchive</title>
</svelte:head>

<div class="min-h-screen bg-dark-950 text-gray-100">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Header -->
		<div class="flex items-center justify-between mb-8">
			<div>
				<h1 class="text-3xl font-bold text-gray-100 flex items-center space-x-3">
					<Download class="h-8 w-8 text-primary-500" />
					<span>My Downloads</span>
					{#if totalDownloads > 0}
						<span class="text-lg text-gray-400 font-normal">({totalDownloads})</span>
					{/if}
				</h1>
				<p class="text-gray-400 mt-2">Track your book download progress and access completed downloads</p>
			</div>
			
			<button
				on:click={refreshDownloads}
				disabled={isLoading}
				class="btn-secondary flex items-center space-x-2"
				class:opacity-50={isLoading}
			>
				{#if isLoading}
					<Loader2 class="h-4 w-4 animate-spin" />
				{:else}
					<RefreshCw class="h-4 w-4" />
				{/if}
				<span>Refresh</span>
			</button>
		</div>

		<!-- Download Status -->
		{#if downloadStatus && !statusLoading}
			<div class="mb-8 p-6 bg-dark-800 border border-gray-700 rounded-lg">
				<h2 class="text-lg font-semibold text-gray-100 mb-4">Download Status</h2>
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
					<!-- Daily Downloads -->
					<div class="text-center p-4 bg-dark-700 rounded-lg">
						<div class="text-2xl font-bold text-gray-100">{downloadStatus.downloads_used}</div>
						<div class="text-sm text-gray-400">Used Today</div>
						<div class="text-xs text-gray-500">of {downloadStatus.daily_limit}</div>
					</div>
					
					<!-- Remaining Downloads -->
					<div class="text-center p-4 bg-dark-700 rounded-lg">
						<div class="text-2xl font-bold text-green-400">{downloadStatus.downloads_remaining}</div>
						<div class="text-sm text-gray-400">Remaining Today</div>
					</div>
					
					<!-- Request Credits -->
					<div class="text-center p-4 bg-dark-700 rounded-lg">
						<div class="text-2xl font-bold text-primary-400 flex items-center justify-center space-x-2">
							<Gift class="h-5 w-5" />
							<span>{downloadStatus.request_credits}</span>
						</div>
						<div class="text-sm text-gray-400">Request Credits</div>
						<div class="text-xs text-gray-500">for extra downloads</div>
					</div>
					
					<!-- Next Reset -->
					<div class="text-center p-4 bg-dark-700 rounded-lg">
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

		<!-- Downloads List -->
		{#if isLoading}
			<div class="flex items-center justify-center py-12">
				<Loader2 class="h-8 w-8 animate-spin text-primary-500" />
				<span class="ml-3 text-gray-400">Loading downloads...</span>
			</div>
		{:else if downloads.length === 0}
			<div class="text-center py-12">
				<Download class="h-16 w-16 text-gray-600 mx-auto mb-4" />
				<h3 class="text-xl font-medium text-gray-400 mb-2">No downloads yet</h3>
				<p class="text-gray-500 mb-6">Start downloading books to see them here</p>
				<a href="/" class="btn-primary">Explore Books</a>
			</div>
		{:else}
			<div class="space-y-4">
				{#each downloads as download (download.id)}
					<div class="card p-6">
						<div class="flex items-start space-x-4">
							<!-- Status Icon -->
							<div class="flex-shrink-0 mt-1">
								{#if download.cover_data}
									<img 
										src={`data:image/jpeg;base64,${download.cover_data}`} 
										alt="Book cover"
										class="w-16 h-20 object-cover rounded-lg"
									/>
								{:else}
									<div class="w-16 h-20 bg-gradient-to-br from-primary-900 to-dark-800 rounded-lg flex items-center justify-center">
										<Download class="h-6 w-6 text-gray-400" />
									</div>
								{/if}
							</div>

							<!-- Download Info -->
							<div class="flex-1 min-w-0">
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<h3 class="text-lg font-semibold text-gray-100 mb-1">
											{download.title || 'Unknown Title'}
										</h3>
										
										{#if download.authors}
											<p class="text-gray-400 mb-2">by {download.authors}</p>
										{/if}
										
										<div class="flex items-center space-x-4 text-sm text-gray-500 mb-3">
											{#if download.format}
												<span>{download.format.toUpperCase()}</span>
											{/if}
											{#if download.size}
												<span>{download.size}</span>
											{/if}
											{#if download.language}
												<span>{download.language}</span>
											{/if}
										</div>
									</div>

									<!-- Status Badge -->
									<div class="flex items-center space-x-2">
										<svelte:component 
											this={getStatusIcon(download.status)} 
											class="h-5 w-5 {getStatusColor(download.status)}" 
										/>
										<span class="text-sm font-medium {getStatusColor(download.status)}">
											{getStatusText(download.status)}
										</span>
									</div>
								</div>

								<!-- Progress Bar -->
								{#if download.status === 'processing' || download.status === 'downloading'}
									<div class="mt-3">
										<div class="flex items-center justify-between text-sm text-gray-400 mb-1">
											<span>Progress</span>
											<span>{download.progress || 0}%</span>
										</div>
										<div class="w-full bg-dark-700 rounded-full h-2">
											<div 
												class="bg-primary-500 h-2 rounded-full transition-all duration-300"
												style="width: {download.progress || 0}%"
											></div>
										</div>
									</div>
								{/if}

								<!-- Error Message -->
								{#if download.error_msg}
									<div class="mt-3 p-3 bg-red-900/20 border border-red-800 rounded-lg">
										<p class="text-red-400 text-sm">
											<AlertCircle class="h-4 w-4 inline mr-2" />
											{download.error_msg}
										</p>
									</div>
								{/if}

								<!-- Timestamps -->
								<div class="flex items-center justify-between mt-4 text-xs text-gray-500">
									<div class="flex items-center space-x-4">
										{#if download.created_at}
											<span>Requested: {formatDate(download.created_at)}</span>
										{/if}
										{#if download.updated_at && download.updated_at !== download.created_at}
											<span>Updated: {formatDate(download.updated_at)}</span>
										{/if}
									</div>

									<!-- Download Button -->
									{#if download.status === 'ready' || download.status === 'completed'}
										<button
											on:click={() => handleDownload(download.book_hash, download.title)}
											disabled={downloadingHashes.has(download.book_hash)}
											class="btn-primary text-xs flex items-center space-x-1"
											class:opacity-50={downloadingHashes.has(download.book_hash)}
										>
											{#if downloadingHashes.has(download.book_hash)}
												<Loader2 class="h-3 w-3 animate-spin" />
												<span>Downloading...</span>
											{:else}
												<Download class="h-3 w-3" />
												<span>Download</span>
											{/if}
										</button>
									{/if}
								</div>
							</div>
						</div>
					</div>
				{/each}
			</div>

			<!-- Pagination -->
			{#if totalPages > 1}
				<div class="mt-8">
					<Pagination
						currentPage={currentPage}
						totalPages={totalPages}
						onPageChange={handlePageChange}
						showInfo={true}
					/>
				</div>
			{/if}
		{/if}
	</div>
</div>
