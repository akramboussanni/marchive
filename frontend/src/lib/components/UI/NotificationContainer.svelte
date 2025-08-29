<script lang="ts">
	import { notifications } from '$lib/stores/notifications';
	import NotificationToast from './NotificationToast.svelte';
</script>

{#if $notifications.length > 0}
	<div class="fixed top-4 right-4 z-50 space-y-2 max-w-sm">
		{#each $notifications as notification (notification.id)}
			<NotificationToast {notification} />
		{/each}
	</div>
{/if}

