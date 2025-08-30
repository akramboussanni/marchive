<script lang="ts">
	import { onMount } from 'svelte';
	import { User, Key, Shield, Calendar, Download, Loader2, LogOut } from 'lucide-svelte';
	import { auth, user, isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	let currentPassword = '';
	let newPassword = '';
	let confirmPassword = '';
	let changingPassword = false;
	let passwordError = '';
	let passwordSuccess = false;

	onMount(async () => {
		if (!$isAuthenticated) {
			// Try to refresh the token before redirecting to login
			const authSuccess = await auth.checkAuthWithRefresh();
			if (!authSuccess) {
				goto('/login');
				return;
			}
		}
	});

	async function handlePasswordChange() {
		if (!currentPassword || !newPassword || !confirmPassword) {
			passwordError = 'Please fill in all fields';
			return;
		}

		if (newPassword !== confirmPassword) {
			passwordError = 'New passwords do not match';
			return;
		}

		if (newPassword.length < 8) {
			passwordError = 'Password must be at least 8 characters long';
			return;
		}

		changingPassword = true;
		passwordError = '';
		passwordSuccess = false;

		try {
			const success = await auth.changePassword(currentPassword, newPassword);
			if (success) {
				passwordSuccess = true;
				currentPassword = '';
				newPassword = '';
				confirmPassword = '';
			} else {
				passwordError = 'Current password is incorrect';
			}
		} catch (error) {
			passwordError = 'Failed to change password. Please try again.';
		} finally {
			changingPassword = false;
		}
	}

	function formatDate(timestamp: string | number) {
		const ts = typeof timestamp === 'string' ? parseInt(timestamp) : timestamp;
		return new Date(ts * 1000).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}

	function getRoleBadgeColor(role: string) {
		switch (role) {
			case 'admin':
				return 'bg-red-900/50 text-red-300 border-red-800';
			case 'user':
				return 'bg-blue-900/50 text-blue-300 border-blue-800';
			default:
				return 'bg-gray-900/50 text-gray-300 border-gray-800';
		}
	}
</script>

<svelte:head>
			<title>Profile - marchive</title>
</svelte:head>

{#if $user}
	<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Header -->
		<div class="mb-8">
			<h1 class="text-2xl font-bold text-gray-100 flex items-center space-x-2">
				<User class="h-6 w-6" />
				<span>Profile Settings</span>
			</h1>
			<p class="text-gray-400 mt-2">
				Manage your account settings and preferences
			</p>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
			<!-- Profile Information -->
			<div class="lg:col-span-2 space-y-6">
				<!-- Account Details -->
				<div class="card p-6">
					<h2 class="text-lg font-semibold text-gray-100 mb-4 flex items-center space-x-2">
						<User class="h-5 w-5" />
						<span>Account Details</span>
					</h2>

					<div class="space-y-4">
						<div>
							<label class="block text-sm font-medium text-gray-300 mb-1">
								Username
							</label>
							<div class="text-gray-100 bg-dark-800 px-3 py-2 rounded-lg border border-gray-700">
								{$user.username}
							</div>
						</div>



						<div>
							<label class="block text-sm font-medium text-gray-300 mb-1">
								Role
							</label>
							<div class={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${getRoleBadgeColor($user.role)}`}>
								<Shield class="h-3 w-3 mr-1" />
								{$user.role.charAt(0).toUpperCase() + $user.role.slice(1)}
							</div>
						</div>

						<div>
							<label class="block text-sm font-medium text-gray-300 mb-1">
								Member Since
							</label>
							<div class="text-gray-400 flex items-center space-x-1">
								<Calendar class="h-4 w-4" />
								<span>{formatDate($user.created_at)}</span>
							</div>
						</div>
					</div>
				</div>

				<!-- Change Password -->
				<div class="card p-6">
					<h2 class="text-lg font-semibold text-gray-100 mb-4 flex items-center space-x-2">
						<Key class="h-5 w-5" />
						<span>Change Password</span>
					</h2>

					<form on:submit|preventDefault={handlePasswordChange} class="space-y-4">
						<div>
							<label for="current-password" class="block text-sm font-medium text-gray-300 mb-1">
								Current Password
							</label>
							<input
								id="current-password"
								type="password"
								bind:value={currentPassword}
								class="input"
								placeholder="Enter your current password"
							/>
						</div>

						<div>
							<label for="new-password" class="block text-sm font-medium text-gray-300 mb-1">
								New Password
							</label>
							<input
								id="new-password"
								type="password"
								bind:value={newPassword}
								class="input"
								placeholder="Enter your new password"
							/>
						</div>

						<div>
							<label for="confirm-password" class="block text-sm font-medium text-gray-300 mb-1">
								Confirm New Password
							</label>
							<input
								id="confirm-password"
								type="password"
								bind:value={confirmPassword}
								class="input"
								placeholder="Confirm your new password"
							/>
						</div>

						{#if passwordError}
							<div class="bg-red-900/50 border border-red-800 text-red-300 px-4 py-3 rounded-lg">
								{passwordError}
							</div>
						{/if}

						{#if passwordSuccess}
							<div class="bg-green-900/50 border border-green-800 text-green-300 px-4 py-3 rounded-lg">
								Password changed successfully!
							</div>
						{/if}

						<button
							type="submit"
							disabled={changingPassword}
							class="btn-primary flex items-center space-x-2"
						>
							{#if changingPassword}
								<Loader2 class="h-4 w-4 animate-spin" />
								<span>Changing Password...</span>
							{:else}
								<Key class="h-4 w-4" />
								<span>Change Password</span>
							{/if}
						</button>
					</form>
				</div>
			</div>

			<!-- Quick Stats -->
			<div class="space-y-6">
				<div class="card p-6">
					<h2 class="text-lg font-semibold text-gray-100 mb-4">
						Quick Actions
					</h2>

					<div class="space-y-3">
						<a href="/downloads" class="w-full btn-secondary flex items-center space-x-2">
							<Download class="h-4 w-4" />
							<span>View Downloads</span>
						</a>

						{#if $user.role === 'admin'}
							<a href="/admin" class="w-full btn-secondary flex items-center space-x-2">
								<Shield class="h-4 w-4" />
								<span>Admin Panel</span>
							</a>
						{/if}

						<button
							on:click={() => auth.logoutEverywhere()}
							class="w-full btn-danger flex items-center space-x-2"
						>
							<LogOut class="h-4 w-4" />
							<span>Logout Everywhere</span>
						</button>
					</div>
				</div>

				<!-- Security Notice -->
				<div class="card p-6">
					<h3 class="text-md font-medium text-gray-200 mb-2">
						Security Tips
					</h3>
					<ul class="text-sm text-gray-400 space-y-2">
						<li>• Use a strong, unique password</li>
						<li>• Don't share your account credentials</li>
						<li>• Log out from shared devices</li>
						<li>• Contact admin if you notice suspicious activity</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
{/if}

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

	.btn-secondary {
		background-color: rgb(31 41 55);
		color: rgb(209 213 219);
		border: 1px solid rgb(75 85 99);
		font-weight: 500;
		padding: 0.5rem 1rem;
		border-radius: 0.5rem;
		transition: background-color 0.15s ease-in-out;
		cursor: pointer;
	}
	
	.btn-secondary:hover {
		background-color: rgb(55 65 81);
	}

	.btn-danger {
		background-color: rgb(220 38 38);
		color: white;
		font-weight: 500;
		padding: 0.5rem 1rem;
		border-radius: 0.5rem;
		transition: background-color 0.15s ease-in-out;
		border: none;
		cursor: pointer;
	}
	
	.btn-danger:hover {
		background-color: rgb(185 28 28);
	}

	.input {
		background-color: rgb(31 41 55);
		color: rgb(209 213 219);
		border: 1px solid rgb(75 85 99);
		padding: 0.5rem 0.75rem;
		border-radius: 0.375rem;
		width: 100%;
		transition: border-color 0.15s ease-in-out;
	}
	
	.input:focus {
		outline: none;
		border-color: rgb(79 70 229);
	}
</style>
