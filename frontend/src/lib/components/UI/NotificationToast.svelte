<script lang="ts">
	import { fly } from 'svelte/transition';
	import { CheckCircle, XCircle, AlertTriangle, Info, X } from 'lucide-svelte';
	import { removeNotification, type Notification } from '$lib/stores/notifications';

	export let notification: Notification;

	$: Icon = getIcon(notification.type);

	function getIcon(type: string) {
		switch (type) {
			case 'success':
				return CheckCircle;
			case 'error':
				return XCircle;
			case 'warning':
				return AlertTriangle;
			case 'info':
				return Info;
			default:
				return Info;
		}
	}

	function getColors(type: string) {
		switch (type) {
			case 'success':
				return 'border-green-500 bg-green-900/20 text-green-400';
			case 'error':
				return 'border-red-500 bg-red-900/20 text-red-400';
			case 'warning':
				return 'border-yellow-500 bg-yellow-900/20 text-yellow-400';
			case 'info':
				return 'border-blue-500 bg-blue-900/20 text-blue-400';
			default:
				return 'border-gray-500 bg-gray-900/20 text-gray-400';
		}
	}

	function dismiss() {
		removeNotification(notification.id);
	}
</script>

<div
	class="notification-toast border-l-4 rounded-lg p-4 shadow-lg backdrop-blur-sm {getColors(notification.type)}"
	transition:fly={{ x: 300, duration: 300 }}
>
	<div class="flex items-start space-x-3">
		<Icon class="h-5 w-5 flex-shrink-0 mt-0.5" />
		
		<div class="flex-1 min-w-0">
			<h4 class="font-medium text-sm">{notification.title}</h4>
			{#if notification.message}
				<p class="text-xs opacity-90 mt-1">{notification.message}</p>
			{/if}
		</div>
		
		<button
			on:click={dismiss}
			class="flex-shrink-0 opacity-70 hover:opacity-100 transition-opacity"
			aria-label="Dismiss notification"
		>
			<X class="h-4 w-4" />
		</button>
	</div>
</div>

<style>
	.notification-toast {
		background-color: rgb(17 24 39 / 0.9);
		border: 1px solid rgb(55 65 81);
		backdrop-filter: blur(8px);
	}
</style>
