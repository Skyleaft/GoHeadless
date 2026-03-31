<script lang="ts">
	import { onMount } from 'svelte';
	import { adminApi } from '$lib/api/admin';
	import { toast } from '$lib/stores/toast';
	import Button from '$lib/shared/Button.svelte';
	import Input from '$lib/shared/Input.svelte';
	import Modal from '$lib/shared/Modal.svelte';
	import type { User } from '$lib/stores/auth';
	import type { Role } from '$lib/api/admin';

	let users = $state<User[]>([]);
	let roles = $state<Role[]>([]);
	let loading = $state(true);
	
	let showCreateModal = $state(false);
	let newUsername = $state('');
	let newPassword = $state('');
	let newRoleId = $state('');
	let creating = $state(false);

	onMount(async () => {
		try {
			const [u, r] = await Promise.all([
				adminApi.getUsers(),
				adminApi.getRoles()
			]);
			users = u;
			roles = r;
		} catch (err: any) {
			toast.error(err.message || 'Failed to load users');
		} finally {
			loading = false;
		}
	});

	async function handleCreateUser() {
		if (!newUsername || !newPassword || !newRoleId) {
			toast.error('All fields are required');
			return;
		}
		creating = true;
		try {
			const newUser = await adminApi.createUser({
				username: newUsername,
				password: newPassword,
				role_id: newRoleId
			});
			users = [...users, newUser];
			toast.success('User created successfully');
			showCreateModal = false;
			newUsername = '';
			newPassword = '';
		} catch (err: any) {
			toast.error(err.message || 'Failed to create user');
		} finally {
			creating = false;
		}
	}

	async function handleDeleteUser(id: string) {
		if (!confirm('Are you sure you want to delete this user?')) return;
		try {
			await adminApi.deleteUser(id);
			users = users.filter(u => u.id !== id);
			toast.success('User deleted');
		} catch (err: any) {
			toast.error(err.message || 'Failed to delete user');
		}
	}
</script>

<div class="flex flex-col gap-6 animate-fade-in">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold">User Management</h1>
			<p class="text-sm text-[--text-muted]">Manage system administrators and editors</p>
		</div>
		<Button onclick={() => (showCreateModal = true)}>+ Create User</Button>
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<div class="h-8 w-8 rounded-full border-4 border-[--brand] border-t-transparent animate-spin"></div>
		</div>
	{:else}
		<div class="card overflow-hidden">
			<table class="w-full text-left text-sm">
				<thead class="bg-[--surface-alt] border-b border-[--border]">
					<tr>
						<th class="px-6 py-3 font-semibold text-[--text-muted] uppercase tracking-wider text-xs">Username</th>
						<th class="px-6 py-3 font-semibold text-[--text-muted] uppercase tracking-wider text-xs">Role</th>
						<th class="px-6 py-3 font-semibold text-[--text-muted] uppercase tracking-wider text-xs">Status</th>
						<th class="px-6 py-3 font-semibold text-[--text-muted] uppercase tracking-wider text-xs text-right">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-[--border]">
					{#each users as user}
						<tr class="hover:bg-[--surface-alt]/50 transition-colors">
							<td class="px-6 py-4">
								<div class="flex items-center gap-3">
									<div class="h-8 w-8 rounded-full bg-[--brand-light] text-[--brand] flex items-center justify-center font-bold">
										{user.username[0].toUpperCase()}
									</div>
									<span class="font-medium">{user.username}</span>
									{#if user.is_initial_admin}
										<span class="px-2 py-0.5 rounded text-[10px] bg-amber-100 text-amber-700 font-bold uppercase tracking-tight">Superadmin</span>
									{/if}
								</div>
							</td>
							<td class="px-6 py-4">
								<span class="text-[--text-secondary]">
									{roles.find(r => r.id === user.role_id)?.name || user.role_id}
								</span>
							</td>
							<td class="px-6 py-4">
								{#if user.active_status}
									<span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium bg-emerald-100 text-emerald-800">
										<span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span> Active
									</span>
								{:else}
									<span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-800">
										<span class="h-1.5 w-1.5 rounded-full bg-slate-400"></span> Disabled
									</span>
								{/if}
							</td>
							<td class="px-6 py-4 text-right">
								<button 
									onclick={() => handleDeleteUser(user.id)}
									disabled={user.is_initial_admin}
									class="text-red-400 hover:text-red-600 disabled:opacity-30 disabled:cursor-not-allowed"
									title={user.is_initial_admin ? "Cannot delete initial admin" : "Delete user"}
								>
									🗑
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<Modal bind:open={showCreateModal} title="Create New User">
	<form onsubmit={(e) => { e.preventDefault(); handleCreateUser(); }} class="flex flex-col gap-4">
		<Input id="newUsername" label="Username" bind:value={newUsername} placeholder="e.g. editor_joe" required />
		<Input id="newPassword" label="Password" type="password" bind:value={newPassword} placeholder="••••••••" required />
		
		<div class="flex flex-col gap-1.5">
			<label for="newRoleId" class="text-sm font-medium">Assign Role</label>
			<select 
				id="newRoleId"
				bind:value={newRoleId}
				class="h-10 rounded-lg border border-[--border] bg-[--surface] px-3 text-sm focus:border-[--brand] focus:ring-1 focus:ring-[--brand] outline-none"
			>
				<option value="">Select a role...</option>
				{#each roles as role}
					<option value={role.id}>{role.name}</option>
				{/each}
			</select>
		</div>

		<div class="flex justify-end gap-3 mt-4">
			<Button variant="ghost" onclick={() => (showCreateModal = false)}>Cancel</Button>
			<Button type="submit" loading={creating}>Create User</Button>
		</div>
	</form>
</Modal>
