<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Search, Plus, Users, Edit, Trash2, Key, Shield, UserX, Loader2 } from 'lucide-svelte';
	import { isAuthenticated, isAdmin } from '$lib/stores/auth';
	import { admin, usersList, type UserWithStats } from '$lib/stores/admin';

	let searchQuery = '';
	let selectedRole = '';
	let loading = true;
	let searchTimeout: NodeJS.Timeout;
	let showCreateModal = false;
	let showEditModal = false;
	let showDeleteModal = false;
	let selectedUser: UserWithStats | null = null;

	// Create user form
	let newUser = {
		username: '',
		password: '',
		role: 'user'
	};

	// Edit user form
	let editUser = {
		username: '',
		role: ''
	};

	let currentPage = 0;
	const pageSize = 20;

	onMount(async () => {
		if (!$isAuthenticated) {
			goto('/login');
			return;
		}

		if (!$isAdmin) {
			goto('/');
			return;
		}

		await loadUsers();
	});

	async function loadUsers() {
		loading = true;
		try {
			if (searchQuery.trim() || selectedRole) {
				await admin.searchUsers(searchQuery, selectedRole, pageSize, 0);
			} else {
				await admin.getUsers(pageSize, 0);
			}
			currentPage = 0;
		} catch (error) {
			console.error('Failed to load users:', error);
		} finally {
			loading = false;
		}
	}

	function debounceSearch() {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(loadUsers, 300);
	}

	async function loadMoreUsers() {
		if (!$usersList || loading) return;

		const nextPage = currentPage + 1;
		try {
			if (searchQuery.trim() || selectedRole) {
				await admin.searchUsers(searchQuery, selectedRole, pageSize, nextPage * pageSize);
			} else {
				await admin.getUsers(pageSize, nextPage * pageSize);
			}
			currentPage = nextPage;
		} catch (error) {
			console.error('Failed to load more users:', error);
		}
	}

	async function handleCreateUser() {
		try {
			await admin.createUser(newUser);
			showCreateModal = false;
			newUser = { username: '', password: '', role: 'user' };
			await loadUsers();
		} catch (error) {
			console.error('Failed to create user:', error);
			alert('Failed to create user');
		}
	}

	async function handleUpdateUser() {
		if (!selectedUser) return;

		try {
			await admin.updateUser(selectedUser.id, editUser);
			showEditModal = false;
			selectedUser = null;
			await loadUsers();
		} catch (error) {
			console.error('Failed to update user:', error);
			alert('Failed to update user');
		}
	}

	async function handleDeleteUser() {
		if (!selectedUser) return;

		try {
			await admin.deleteUser(selectedUser.id);
			showDeleteModal = false;
			selectedUser = null;
			await loadUsers();
		} catch (error) {
			console.error('Failed to delete user:', error);
			alert('Failed to delete user');
		}
	}

	async function invalidateUserSessions(user: UserWithStats) {
		try {
			await admin.invalidateUserSessions(user.id);
			alert('User sessions invalidated successfully');
		} catch (error) {
			console.error('Failed to invalidate sessions:', error);
			alert('Failed to invalidate sessions');
		}
	}

	function openEditModal(user: UserWithStats) {
		selectedUser = user;
		editUser = {
			username: user.username,
			role: user.role,
		};
		showEditModal = true;
	}

	function openDeleteModal(user: UserWithStats) {
		selectedUser = user;
		showDeleteModal = true;
	}

	function formatDate(timestamp: string | number) {
		const ts = typeof timestamp === 'string' ? parseInt(timestamp) : timestamp;
		return new Date(ts * 1000).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
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

	$: if (searchQuery !== undefined) {
		debounceSearch();
	}

	$: if (selectedRole !== undefined) {
		loadUsers();
	}

	$: hasMoreUsers = $usersList?.pagination.has_next ?? false;
	$: displayUsers = $usersList?.users ?? [];
</script>

<svelte:head>
			<title>User Management - marchive Admin</title>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
	<!-- Header -->
	<div class="flex items-center justify-between mb-8">
		<div>
			<h1 class="text-2xl font-bold text-gray-100 flex items-center space-x-2">
				<Users class="h-6 w-6" />
				<span>User Management</span>
			</h1>
			<p class="text-gray-400 mt-2">
				Manage user accounts and permissions
			</p>
		</div>

		<button
			on:click={() => showCreateModal = true}
			class="btn-primary flex items-center space-x-2"
		>
			<Plus class="h-4 w-4" />
			<span>Add User</span>
		</button>
	</div>

	<!-- Filters -->
	<div class="card p-6 mb-6">
		<div class="flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4">
			<!-- Search -->
			<div class="flex-1 relative">
				<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
					<Search class="h-5 w-5 text-gray-400" />
				</div>
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search by username..."
					class="block w-full pl-10 input"
				/>
			</div>

			<!-- Role Filter -->
			<select bind:value={selectedRole} class="input w-auto">
				<option value="">All Roles</option>
				<option value="user">Users</option>
				<option value="admin">Admins</option>
			</select>
		</div>
	</div>

	<!-- Users List -->
	{#if loading}
		<div class="flex items-center justify-center py-16">
			<Loader2 class="h-8 w-8 animate-spin text-primary-500" />
			<span class="ml-3 text-gray-400">Loading users...</span>
		</div>
	{:else if displayUsers.length > 0}
		<div class="card overflow-hidden">
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-800">
					<thead class="bg-dark-800">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
								User
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
								Role
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
								Downloads
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
								Joined
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-dark-900 divide-y divide-gray-800">
						{#each displayUsers as user (user.id)}
							<tr class="hover:bg-dark-800">
								<td class="px-6 py-4 whitespace-nowrap">
									<div>
										<div class="text-sm font-medium text-gray-100">
											{user.username}
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${getRoleBadgeColor(user.role)}`}>
										<Shield class="h-3 w-3 mr-1" />
										{user.role.charAt(0).toUpperCase() + user.role.slice(1)}
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
									{user.download_count}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-400">
									{formatDate(user.created_at)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<div class="flex items-center space-x-2">
										<button
											on:click={() => openEditModal(user)}
											class="text-primary-400 hover:text-primary-300"
											title="Edit User"
										>
											<Edit class="h-4 w-4" />
										</button>
										<button
											on:click={() => invalidateUserSessions(user)}
											class="text-yellow-400 hover:text-yellow-300"
											title="Invalidate Sessions"
										>
											<UserX class="h-4 w-4" />
										</button>
										<button
											on:click={() => openDeleteModal(user)}
											class="text-red-400 hover:text-red-300"
											title="Delete User"
										>
											<Trash2 class="h-4 w-4" />
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>

		<!-- Load More -->
		{#if hasMoreUsers}
			<div class="flex justify-center mt-6">
				<button
					on:click={loadMoreUsers}
					class="btn-secondary"
				>
					Load More Users
				</button>
			</div>
		{/if}
	{:else}
		<div class="text-center py-16">
			<Users class="h-16 w-16 text-gray-600 mx-auto mb-4" />
			<h3 class="text-lg font-medium text-gray-300 mb-2">No users found</h3>
			<p class="text-gray-500">Try adjusting your search criteria</p>
		</div>
	{/if}
</div>

<!-- Create User Modal -->
{#if showCreateModal}
	<div class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
		<div class="bg-dark-900 rounded-xl border border-gray-800 w-full max-w-md">
			<div class="p-6">
				<h2 class="text-lg font-semibold text-gray-100 mb-4">Create New User</h2>
				
				<form on:submit|preventDefault={handleCreateUser} class="space-y-4">
					<div>
						<label class="block text-sm font-medium text-gray-300 mb-1">Username</label>
						<input
							type="text"
							bind:value={newUser.username}
							required
							class="input"
							placeholder="Enter username"
						/>
					</div>



					<div>
						<label class="block text-sm font-medium text-gray-300 mb-1">Password</label>
						<input
							type="password"
							bind:value={newUser.password}
							required
							class="input"
							placeholder="Enter password"
						/>
					</div>

					<div>
						<label class="block text-sm font-medium text-gray-300 mb-1">Role</label>
						<select bind:value={newUser.role} class="input">
							<option value="user">User</option>
							<option value="admin">Admin</option>
						</select>
					</div>

					<div class="flex space-x-3 pt-4">
						<button type="submit" class="flex-1 btn-primary">
							Create User
						</button>
						<button
							type="button"
							on:click={() => showCreateModal = false}
							class="flex-1 btn-secondary"
						>
							Cancel
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}

<!-- Edit User Modal -->
{#if showEditModal && selectedUser}
	<div class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
		<div class="bg-dark-900 rounded-xl border border-gray-800 w-full max-w-md">
			<div class="p-6">
				<h2 class="text-lg font-semibold text-gray-100 mb-4">Edit User</h2>
				
				<form on:submit|preventDefault={handleUpdateUser} class="space-y-4">
					<div>
						<label class="block text-sm font-medium text-gray-300 mb-1">Username</label>
						<input
							type="text"
							bind:value={editUser.username}
							required
							class="input"
						/>
					</div>



					<div>
						<label class="block text-sm font-medium text-gray-300 mb-1">Role</label>
						<select bind:value={editUser.role} class="input">
							<option value="user">User</option>
							<option value="admin">Admin</option>
						</select>
					</div>



					<div class="flex space-x-3 pt-4">
						<button type="submit" class="flex-1 btn-primary">
							Update User
						</button>
						<button
							type="button"
							on:click={() => showEditModal = false}
							class="flex-1 btn-secondary"
						>
							Cancel
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}

<!-- Delete User Modal -->
{#if showDeleteModal && selectedUser}
	<div class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
		<div class="bg-dark-900 rounded-xl border border-gray-800 w-full max-w-md">
			<div class="p-6">
				<h2 class="text-lg font-semibold text-gray-100 mb-4">Delete User</h2>
				<p class="text-gray-400 mb-6">
					Are you sure you want to delete user <strong>{selectedUser.username}</strong>? 
					This action cannot be undone.
				</p>
				
				<div class="flex space-x-3">
					<button
						on:click={handleDeleteUser}
						class="flex-1 bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
					>
						Delete User
					</button>
					<button
						on:click={() => showDeleteModal = false}
						class="flex-1 btn-secondary"
					>
						Cancel
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
