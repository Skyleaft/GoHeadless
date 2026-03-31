<script lang="ts">
	import { onMount } from 'svelte';
	import { adminApi, type Role } from '$lib/api/admin';
	import { toast } from '$lib/stores/toast';
	import Button from '$lib/shared/Button.svelte';
	import Input from '$lib/shared/Input.svelte';
	import Modal from '$lib/shared/Modal.svelte';

	let roles = $state<Role[]>([]);
	let loading = $state(true);
	
	let showCreateModal = $state(false);
	let newRoleName = $state('');
	let newRoleDesc = $state('');
	let creating = $state(false);

	onMount(async () => {
		try {
			roles = await adminApi.getRoles();
		} catch (err: any) {
			toast.error(err.message || 'Failed to load roles');
		} finally {
			loading = false;
		}
	});

	async function handleCreateRole() {
		if (!newRoleName) {
			toast.error('Role name is required');
			return;
		}
		creating = true;
		try {
			const newRole = await adminApi.createRole({
				name: newRoleName,
				description: newRoleDesc,
				permissions: []
			});
			roles = [...roles, newRole];
			toast.success('Role created successfully');
			showCreateModal = false;
			newRoleName = '';
			newRoleDesc = '';
		} catch (err: any) {
			toast.error(err.message || 'Failed to create role');
		} finally {
			creating = false;
		}
	}

	async function handleDeleteRole(id: string | undefined) {
		if (!id) return;
		if (!confirm('Are you sure you want to delete this role? Users with this role will lose access!')) return;
		try {
			await adminApi.deleteRole(id);
			roles = roles.filter(r => r.id !== id);
			toast.success('Role deleted');
		} catch (err: any) {
			toast.error(err.message || 'Failed to delete role');
		}
	}
</script>

<div class="flex flex-col gap-6 animate-fade-in">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold">Role Management</h1>
			<p class="text-sm text-[--text-muted]">Define access scopes and permissions</p>
		</div>
		<Button onclick={() => (showCreateModal = true)}>+ Create Role</Button>
	</div>

	{#if loading}
		<div class="flex justify-center py-12">
			<div class="h-8 w-8 rounded-full border-4 border-[--brand] border-t-transparent animate-spin"></div>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each roles as role}
				<div class="card p-5 group hover:shadow-lg transition-all h-full flex flex-col">
					<div class="flex items-start justify-between mb-2">
						<div class="h-10 w-10 rounded-lg bg-[--brand-light] text-[--brand] flex items-center justify-center">
							🛡️
						</div>
						<div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
							<a 
								href="/admin/roles/{role.id}" 
								class="h-8 w-8 rounded bg-[--surface-alt] flex items-center justify-center hover:bg-[--brand-light] hover:text-[--brand]"
								title="Edit Role"
							>✎</a>
							<button 
								onclick={() => handleDeleteRole(role.id)}
								class="h-8 w-8 rounded bg-[--surface-alt] flex items-center justify-center hover:bg-red-50 hover:text-red-500"
								title="Delete Role"
							>🗑</button>
						</div>
					</div>
					<h3 class="font-bold text-lg">{role.name}</h3>
					<p class="text-sm text-[--text-muted] mt-1 flex-1 line-clamp-2">{role.description || 'No description provided'}</p>
					
					<div class="mt-4 pt-4 border-t border-[--border] flex items-center justify-between">
						<span class="text-xs font-medium text-[--text-secondary]">
							{role.permissions?.length || 0} Permissions set
						</span>
						<a 
							href="/admin/roles/{role.id}" 
							class="text-xs font-bold text-[--brand] hover:underline"
						>Manage Permissions →</a>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<Modal bind:open={showCreateModal} title="Create New Role">
	<form onsubmit={(e) => { e.preventDefault(); handleCreateRole(); }} class="flex flex-col gap-4">
		<Input label="Role Name" bind:value={newRoleName} placeholder="e.g. Content Writer" required />
		<Input label="Description" bind:value={newRoleDesc} placeholder="Can create and edit blog posts..." />
		
		<div class="flex justify-end gap-3 mt-4">
			<Button variant="ghost" onclick={() => (showCreateModal = false)}>Cancel</Button>
			<Button type="submit" loading={creating}>Create Role</Button>
		</div>
	</form>
</Modal>
