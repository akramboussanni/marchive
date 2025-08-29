import { writable } from 'svelte/store';

export interface Notification {
	id: string;
	type: 'success' | 'error' | 'warning' | 'info';
	title: string;
	message?: string;
	duration?: number; // in milliseconds, 0 means no auto-dismiss
}

export const notifications = writable<Notification[]>([]);

let notificationId = 0;

export function addNotification(notification: Omit<Notification, 'id'>): string {
	const id = `notification-${notificationId++}`;
	const newNotification: Notification = {
		id,
		duration: 5000, // Default 5 seconds
		...notification
	};

	notifications.update(n => [...n, newNotification]);

	// Auto-dismiss if duration is set
	if (newNotification.duration && newNotification.duration > 0) {
		setTimeout(() => {
			removeNotification(id);
		}, newNotification.duration);
	}

	return id;
}

export function removeNotification(id: string) {
	notifications.update(n => n.filter(notification => notification.id !== id));
}

export function clearAllNotifications() {
	notifications.set([]);
}

// Convenience functions
export const showSuccess = (title: string, message?: string) => 
	addNotification({ type: 'success', title, message });

export const showError = (title: string, message?: string) => 
	addNotification({ type: 'error', title, message, duration: 0 }); // Errors don't auto-dismiss

export const showWarning = (title: string, message?: string) => 
	addNotification({ type: 'warning', title, message });

export const showInfo = (title: string, message?: string) => 
	addNotification({ type: 'info', title, message });


