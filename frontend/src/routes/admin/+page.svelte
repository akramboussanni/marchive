<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { BarChart3, Users, BookOpen, Download, Activity, TrendingUp, Gift } from 'lucide-svelte';
	import { isAuthenticated, isAdmin } from '$lib/stores/auth';
	import { admin, systemStats, type SystemStats } from '$lib/stores/admin';

	let stats: SystemStats | null = null;
	let loading = true;

	onMount(async () => {
		if (!$isAuthenticated) {
			// Try to refresh the token before redirecting to login
			const authSuccess = await auth.checkAuthWithRefresh();
			if (!authSuccess) {
				goto('/login');
				return;
			}
		}

		if (!$isAdmin) {
			goto('/');
			return;
		}

		await loadStats();
	});

	async function loadStats() {
		loading = true;
		try {
			console.log('Loading admin stats...');
			stats = await admin.getSystemStats();
			console.log('Stats loaded:', stats);
		} catch (error) {
			console.error('Failed to load stats:', error);
		} finally {
			loading = false;
		}
	}

	function formatNumber(num: number): string {
		if (num >= 1000000) {
			return (num / 1000000).toFixed(1) + 'M';
		} else if (num >= 1000) {
			return (num / 1000).toFixed(1) + 'K';
		}
		return num.toString();
	}

	function formatDate(timestamp: string | number): string {
		const ts = typeof timestamp === 'string' ? parseInt(timestamp) : timestamp;
		return new Date(ts * 1000).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<svelte:head>
			<title>Admin Dashboard - marchive</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
	<!-- Header -->
	<div class="mb-6 sm:mb-8">
		<h1 class="text-xl sm:text-2xl font-bold text-gray-100 flex items-center space-x-2">
			<BarChart3 class="h-5 w-5 sm:h-6 sm:w-6" />
			<span>Admin Dashboard</span>
		</h1>
		<p class="text-gray-400 mt-2 text-sm sm:text-base">
			System overview and management tools
		</p>
	</div>

	<!-- Quick Actions -->
	<div class="mb-8">
		<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
			<a href="/admin/users" class="card p-6 hover:bg-dark-800 transition-colors inline-block">
				<div class="flex items-center space-x-3">
					<Users class="h-8 w-8 text-primary-500" />
					<div>
						<h3 class="font-medium text-gray-100">Manage Users</h3>
						<p class="text-sm text-gray-400">View and edit user accounts</p>
					</div>
				</div>
			</a>
			
			<a href="/admin/credits" class="card p-6 hover:bg-dark-800 transition-colors inline-block">
				<div class="flex items-center space-x-3">
					<Gift class="h-8 w-8 text-primary-500" />
					<div>
						<h3 class="font-medium text-gray-100">Request Credits</h3>
						<p class="text-sm text-gray-400">Manage user request credits</p>
					</div>
				</div>
			</a>
		</div>
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-16">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
			<span class="ml-3 text-gray-400">Loading dashboard...</span>
		</div>
	{:else if !stats}
		<div class="flex items-center justify-center py-16">
			<div class="text-center">
				<div class="text-6xl mb-4">⚠️</div>
				<h3 class="text-lg font-medium text-gray-300 mb-2">Failed to load dashboard</h3>
				<p class="text-gray-500 mb-6">Unable to fetch system statistics</p>
				<button
					on:click={loadStats}
					class="btn-primary"
				>
					Retry
				</button>
			</div>
		</div>
	{:else if stats}
		<!-- Stats Cards -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6 mb-6 sm:mb-8">
			<!-- Total Users -->
			<div class="card p-6">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<Users class="h-8 w-8 text-blue-400" />
					</div>
					<div class="ml-4">
						<p class="text-2xl font-bold text-gray-100">
							{formatNumber(stats.total_users)}
						</p>
						<p class="text-sm text-gray-400">Total Users</p>
					</div>
				</div>
			</div>

			<!-- Total Books -->
			<div class="card p-6">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<BookOpen class="h-8 w-8 text-green-400" />
					</div>
					<div class="ml-4">
						<p class="text-2xl font-bold text-gray-100">
							{formatNumber(stats.total_books)}
						</p>
						<p class="text-sm text-gray-400">Total Books</p>
					</div>
				</div>
			</div>

			<!-- Total Downloads -->
			<div class="card p-6">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<Download class="h-8 w-8 text-purple-400" />
					</div>
					<div class="ml-4">
						<p class="text-2xl font-bold text-gray-100">
							{formatNumber(stats.total_downloads)}
						</p>
						<p class="text-sm text-gray-400">Total Downloads</p>
					</div>
				</div>
			</div>

			<!-- Active Users -->
			<div class="card p-6">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<Activity class="h-8 w-8 text-yellow-400" />
					</div>
					<div class="ml-4">
						<p class="text-6xl font-bold text-gray-100">
							{formatNumber(stats.active_users_24h)}
						</p>
						<p class="text-sm text-gray-400">Active (24h)</p>
					</div>
				</div>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6 sm:gap-8">
			<!-- Recent Downloads -->
			<div class="card p-6">
				<h2 class="text-lg font-semibold text-gray-100 mb-4 flex items-center space-x-2">
					<Download class="h-5 w-5" />
					<span>Recent Downloads</span>
				</h2>

				{#if stats.recent_downloads.length > 0}
					<div class="space-y-3">
						{#each stats.recent_downloads as download}
							<div class="flex items-center justify-between py-2 border-b border-gray-800 last:border-b-0">
								<div class="flex-1 min-w-0">
									<p class="text-sm font-medium text-gray-200 truncate">
										{download.title}
									</p>
									<p class="text-xs text-gray-400">
										User ID: {download.user_id}
									</p>
								</div>
								<div class="text-xs text-gray-500">
									{formatDate(download.created_at)}
								</div>
							</div>
						{/each}
					</div>
				{:else}
					<p class="text-gray-400 text-center py-4">No recent downloads</p>
				{/if}
			</div>

			<!-- Top Books -->
			<div class="card p-6">
				<h2 class="text-lg font-semibold text-gray-100 mb-4 flex items-center space-x-2">
					<TrendingUp class="h-5 w-5" />
					<span>Popular Books</span>
				</h2>

				{#if stats.top_books.length > 0}
					<div class="space-y-3">
						{#each stats.top_books as book}
							<div class="flex items-center justify-between py-2 border-b border-gray-800 last:border-b-0">
								<div class="flex-1 min-w-0">
									<p class="text-sm font-medium text-gray-200 truncate">
										{book.title}
									</p>
									<p class="text-xs text-gray-400">
										{book.authors}
									</p>
								</div>
								<div class="text-xs text-primary-400 font-medium">
									{book.download_count} downloads
								</div>
							</div>
						{/each}
					</div>
				{:else}
					<p class="text-gray-400 text-center py-4">No download data available</p>
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	.card {
		background-color: rgb(31 41 55);
		border: 1px solid rgb(75 85 99);
		border-radius: 0.5rem;
		transition: all 0.2s ease-in-out;
	}
	
	.card:hover {
		border-color: rgb(107 114 128);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
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
</style>
