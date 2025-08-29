<script lang="ts">
	import type { Book } from '$lib/stores/books';

	export let book: Book;
	export let size: 'xs' | 'sm' | 'md' | 'lg' = 'md';

	const sizeClasses = {
		xs: 'w-16 h-20',
		sm: 'w-20 h-28',
		md: 'w-32 h-44',
		lg: 'w-40 h-56'
	};

	function generateFallbackCover(title: string): string {
		const colors = [
			'from-blue-900 to-blue-700',
			'from-purple-900 to-purple-700',
			'from-green-900 to-green-700',
			'from-red-900 to-red-700',
			'from-indigo-900 to-indigo-700',
			'from-pink-900 to-pink-700',
			'from-yellow-900 to-yellow-700',
			'from-gray-900 to-gray-700'
		];
		
		const hash = title.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
		return colors[hash % colors.length];
	}

	$: gradientClass = generateFallbackCover(book.title);
	$: displayTitle = book.title.length > 60 ? book.title.substring(0, 60) + '...' : book.title;
</script>

<div class={`relative ${sizeClasses[size]} rounded-lg overflow-hidden shadow-lg`}>
	{#if book.cover_url}
		<img
			src={book.cover_url}
			alt={book.title}
			class="w-full h-full object-cover"
			loading="lazy"
		/>
	{:else}
		<div class={`w-full h-full bg-gradient-to-br ${gradientClass} flex items-center justify-center p-2`}>
			<div class="text-center">
				<div class={`text-white font-medium ${size === 'xs' ? 'text-xs' : size === 'sm' ? 'text-xs' : 'text-sm'} leading-tight`}>
					{displayTitle}
				</div>
				{#if book.authors && size !== 'xs'}
					<div class={`text-gray-300 ${size === 'sm' ? 'text-xs' : 'text-xs'} mt-1 opacity-80`}>
						{book.authors.length > 30 ? book.authors.substring(0, 30) + '...' : book.authors}
					</div>
				{/if}
			</div>
		</div>
	{/if}

	<!-- Format badge -->
	{#if book.format && size !== 'xs'}
		<div class="absolute top-1 right-1 bg-black/70 text-white px-1 py-0.5 rounded text-xs font-medium">
			{book.format.toUpperCase()}
		</div>
	{/if}

	<!-- Download count badge (if available) -->
	{#if book.download_count !== undefined && book.download_count > 0 && size !== 'xs'}
		<div class="absolute bottom-1 left-1 bg-primary-600/90 text-white px-1 py-0.5 rounded text-xs font-medium">
			{book.download_count} downloads
		</div>
	{/if}
</div>
