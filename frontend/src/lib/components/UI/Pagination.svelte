<script lang="ts">
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';
	import { createEventDispatcher } from 'svelte';

	export let currentPage: number;
	export let totalPages: number;
	export let maxVisiblePages: number = 5;

	const dispatch = createEventDispatcher<{ pageChange: number }>();

	function handlePageChange(page: number) {
		if (page >= 0 && page < totalPages && page !== currentPage) {
			dispatch('pageChange', page);
		}
	}

	function getVisiblePages() {
		const pages: (number | string)[] = [];
		
		if (totalPages <= maxVisiblePages) {
			// Show all pages if total is small
			for (let i = 0; i < totalPages; i++) {
				pages.push(i);
			}
		} else {
			// Always show first page
			pages.push(0);
			
			let start = Math.max(1, currentPage - Math.floor(maxVisiblePages / 2));
			let end = Math.min(totalPages - 1, start + maxVisiblePages - 3);
			
			// Adjust start if we're near the end
			if (end === totalPages - 1) {
				start = Math.max(1, end - maxVisiblePages + 3);
			}
			
			// Add ellipsis if there's a gap
			if (start > 1) {
				pages.push('...');
			}
			
			// Add middle pages
			for (let i = start; i <= end; i++) {
				pages.push(i);
			}
			
			// Add ellipsis if there's a gap
			if (end < totalPages - 1) {
				pages.push('...');
			}
			
			// Always show last page
			pages.push(totalPages - 1);
		}
		
		return pages;
	}

	$: visiblePages = getVisiblePages();
</script>

<div class="flex items-center justify-center space-x-1">
	<!-- Previous button -->
	<button
		on:click={() => handlePageChange(currentPage - 1)}
		disabled={currentPage === 0}
		class="px-3 py-2 text-sm font-medium text-gray-400 bg-dark-800 border border-gray-700 rounded-md hover:bg-dark-700 hover:text-gray-300 disabled:bg-dark-900 disabled:text-gray-600 disabled:cursor-not-allowed transition-colors"
		aria-label="Previous page"
	>
		<ChevronLeft class="w-4 h-4" />
	</button>

	<!-- Page numbers -->
	{#each visiblePages as page}
		{#if typeof page === 'number'}
			<button
				on:click={() => handlePageChange(page)}
				class="px-3 py-2 text-sm font-medium rounded-md transition-colors {page === currentPage 
					? 'bg-primary-600 text-white border border-primary-600' 
					: 'text-gray-400 bg-dark-800 border border-gray-700 hover:bg-dark-700 hover:text-gray-300'}"
				aria-label="Page {page + 1}"
				aria-current={page === currentPage ? 'page' : undefined}
			>
				{page + 1}
			</button>
		{:else}
			<span class="px-3 py-2 text-sm text-gray-500">...</span>
		{/if}
	{/each}

	<!-- Next button -->
	<button
		on:click={() => handlePageChange(currentPage + 1)}
		disabled={currentPage === totalPages - 1}
		class="px-3 py-2 text-sm font-medium text-gray-400 bg-dark-800 border border-gray-700 rounded-md hover:bg-dark-700 hover:text-gray-300 disabled:bg-dark-900 disabled:text-gray-600 disabled:cursor-not-allowed transition-colors"
		aria-label="Next page"
	>
		<ChevronRight class="w-4 h-4" />
	</button>
</div>
