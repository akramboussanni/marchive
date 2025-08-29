<script lang="ts">
	import { Heart, Star } from 'lucide-svelte';
	import { isAuthenticated } from '$lib/stores/auth';
	import { books, userFavorites } from '$lib/stores/books';
	import { localFavorites } from '$lib/utils/localFavorites';
	import { showSuccess, showError } from '$lib/stores/notifications';

	export let bookHash: string;
	export let size: 'sm' | 'md' = 'md';

	let isLoading = false;

	// Reactive favorite state
	$: isFavorited = $isAuthenticated ? 
		($userFavorites?.books?.some((book: any) => book.hash === bookHash) ?? false) : 
		localFavorites.isFavorited(bookHash);

	async function toggleFavorite() {
		if (isLoading) return;
		
		isLoading = true;
		
		try {
			if ($isAuthenticated) {
				// Use backend API
				const result = await books.toggleFavorite(bookHash);
				isFavorited = result.isFavorited;
				showSuccess('Success', result.message);
				
				// Refresh favorites list for authenticated users
				await books.getUserFavorites(24, 0);
			} else {
				// Use local storage
				if (isFavorited) {
					localFavorites.remove(bookHash);
					isFavorited = false;
					showSuccess('Success', 'Book removed from favorites');
				} else {
					localFavorites.add(bookHash);
					isFavorited = true;
					showSuccess('Success', 'Book added to favorites');
				}
			}
		} catch (error) {
			console.error('Failed to toggle favorite:', error);
			showError('Error', 'Failed to update favorite status');
		} finally {
			isLoading = false;
		}
	}


</script>

<button
	on:click={toggleFavorite}
	disabled={isLoading}
	class="favorite-btn"
	class:size-sm={size === 'sm'}
	class:size-md={size === 'md'}
	class:is-favorited={isFavorited}
	class:is-loading={isLoading}
	class:always-visible={isFavorited}
	aria-label="{isFavorited ? 'Remove from favorites' : 'Add to favorites'}"
>
	<Heart class="heart-icon" />
</button>

<style>
	.favorite-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		border: none;
		cursor: pointer;
		transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
		background: rgba(17, 24, 39, 0.9);
		backdrop-filter: blur(8px);
		border-radius: 50%;
		position: relative;
		overflow: hidden;
		opacity: 0;
		transform: scale(0.8);
	}

	/* Show on hover or when favorited */
	.favorite-btn.always-visible {
		opacity: 1;
		transform: scale(1);
	}

	/* Show on container hover */
	:global(.book-card:hover) .favorite-btn {
		opacity: 1;
		transform: scale(1);
	}

	.favorite-btn.size-sm {
		width: 2rem;
		height: 2rem;
	}

	.favorite-btn.size-md {
		width: 2.5rem;
		height: 2.5rem;
	}

	.heart-icon {
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		z-index: 1;
		position: relative;
	}

	.favorite-btn.size-sm .heart-icon {
		width: 0.625rem;
		height: 0.625rem;
	}

	.favorite-btn.size-md .heart-icon {
		width: 0.75rem;
		height: 0.75rem;
	}

	/* Unfavorited state */
	.favorite-btn:not(.is-favorited) {
		border: 1px solid rgba(156, 163, 175, 0.4);
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
	}

	.favorite-btn:not(.is-favorited) .heart-icon {
		color: rgb(156, 163, 175);
		fill: none;
	}

	.favorite-btn:not(.is-favorited):hover {
		border-color: rgba(239, 68, 68, 0.6);
		box-shadow: 0 4px 16px rgba(239, 68, 68, 0.3);
		transform: scale(1.1);
	}

	.favorite-btn:not(.is-favorited):hover .heart-icon {
		color: rgb(239, 68, 68);
		transform: scale(1.1);
	}

	/* Favorited state */
	.favorite-btn.is-favorited {
		background: rgba(239, 68, 68, 0.9);
		border: 1px solid rgba(239, 68, 68, 0.6);
		box-shadow: 0 4px 16px rgba(239, 68, 68, 0.4);
	}

	.favorite-btn.is-favorited .heart-icon {
		color: white;
		fill: white;
		animation: heartBeat 0.6s ease-in-out;
	}

	.favorite-btn.is-favorited:hover {
		transform: scale(1.1);
		box-shadow: 0 6px 20px rgba(239, 68, 68, 0.5);
	}

	.favorite-btn.is-favorited:hover .heart-icon {
		transform: scale(1.2);
	}

	/* Loading state */
	.favorite-btn.is-loading {
		pointer-events: none;
	}

	.favorite-btn.is-loading .heart-icon {
		animation: pulse 1.5s ease-in-out infinite;
	}

	/* Active state */
	.favorite-btn:active {
		transform: scale(0.9);
	}

	@keyframes heartBeat {
		0% {
			transform: scale(1);
		}
		14% {
			transform: scale(1.3);
		}
		28% {
			transform: scale(1);
		}
		42% {
			transform: scale(1.3);
		}
		70% {
			transform: scale(1);
		}
	}

	@keyframes pulse {
		0%, 100% {
			opacity: 1;
		}
		50% {
			opacity: 0.6;
		}
	}

	/* Focus state for accessibility */
	.favorite-btn:focus {
		outline: 2px solid rgba(59, 130, 246, 0.5);
		outline-offset: 2px;
	}
</style>
