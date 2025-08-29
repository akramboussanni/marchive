<script lang="ts">
	import BookShelf from '$lib/components/Books/BookShelf.svelte';
</script>

<svelte:head>
	<title>marchive - Public Digital Library</title>
	<meta name="description" content="Discover and download books from our public digital collection" />
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
	<BookShelf title="Public Library" />
</div>
