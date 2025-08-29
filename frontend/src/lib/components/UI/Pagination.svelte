<script lang="ts">
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';

	export let currentPage: number;
	export let totalPages: number;
	export let onPageChange: (page: number) => void;
	export let showInfo = true;

	function handlePrevPage() {
		if (currentPage > 0) {
			onPageChange(currentPage - 1);
		}
	}

	function handleNextPage() {
		if (currentPage < totalPages - 1) {
			onPageChange(currentPage + 1);
		}
	}

	function handlePageClick(page: number) {
		onPageChange(page);
	}

	function getVisiblePages(): (number | string)[] {
		if (totalPages <= 7) {
			return Array.from({ length: totalPages }, (_, i) => i);
		}

		const pages: (number | string)[] = [];
		
		if (currentPage <= 3) {
			// Show first 5 pages + ... + last page
			for (let i = 0; i < 5; i++) {
				pages.push(i);
			}
			pages.push('...');
			pages.push(totalPages - 1);
		} else if (currentPage >= totalPages - 4) {
			// Show first page + ... + last 5 pages
			pages.push(0);
			pages.push('...');
			for (let i = totalPages - 5; i < totalPages; i++) {
				pages.push(i);
			}
		} else {
			// Show first page + ... + current page + ... + last page
			pages.push(0);
			pages.push('...');
			for (let i = currentPage - 1; i <= currentPage + 1; i++) {
				pages.push(i);
			}
			pages.push('...');
			pages.push(totalPages - 1);
		}
		
		return pages;
	}

	$: hasPrev = currentPage > 0;
	$: hasNext = currentPage < totalPages - 1;
	$: visiblePages = getVisiblePages();
</script>

<div class="flex items-center justify-between">
	{#if showInfo}
		<div class="text-sm text-gray-400">
			Page {currentPage + 1} of {totalPages}
		</div>
	{/if}
	
	<div class="flex items-center space-x-2">
		<!-- Previous Button -->
		<button
			on:click={handlePrevPage}
			disabled={!hasPrev}
			class="btn-secondary-sm flex items-center space-x-1"
			class:opacity-50={!hasPrev}
		>
			<ChevronLeft class="h-4 w-4" />
			<span>Previous</span>
		</button>

		<!-- Page Numbers -->
		<div class="flex items-center space-x-1">
			{#each visiblePages as page}
				{#if page === '...'}
					<span class="px-3 py-2 text-gray-400">...</span>
				{:else}
					<button
						on:click={() => handlePageClick(typeof page === 'number' ? page : 0)}
						class="px-3 py-2 rounded-md text-sm transition-colors"
						class:bg-primary-600={page === currentPage}
						class:text-white={page === currentPage}
						class:text-gray-300={page !== currentPage}
						class:hover:bg-primary-700={page === currentPage}
						class:hover:bg-gray-700={page !== currentPage}
					>
						{typeof page === 'number' ? page + 1 : '1'}
					</button>
				{/if}
			{/each}
		</div>

		<!-- Next Button -->
		<button
			on:click={handleNextPage}
			disabled={!hasNext}
			class="btn-secondary-sm flex items-center space-x-1"
			class:opacity-50={!hasNext}
		>
			<span>Next</span>
			<ChevronRight class="h-4 w-4" />
		</button>
	</div>
</div>

<style>
	.btn-secondary-sm {
		background-color: rgb(31 41 55);
		color: rgb(209 213 219);
		border: 1px solid rgb(75 85 99);
		font-weight: 500;
		padding: 0.375rem 0.75rem;
		border-radius: 0.375rem;
		transition: background-color 0.15s ease-in-out;
		cursor: pointer;
		font-size: 0.875rem;
	}
	.btn-secondary-sm:hover:not(:disabled) {
		background-color: rgb(55 65 81);
	}
	.btn-secondary-sm:disabled {
		cursor: not-allowed;
	}
</style>
