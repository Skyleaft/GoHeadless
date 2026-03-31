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
		if (
			!confirm('Are you sure you want to delete this role? Users with this role will lose access!')
		)
			return;
		try {
			await adminApi.deleteRole(id);
			roles = roles.filter((r) => r.id !== id);
			toast.success('Role deleted');
		} catch (err: any) {
			toast.error(err.message || 'Failed to delete role');
		}
	}
</script>

<div class="max-w-6xl animate-fade-in gap-10 py-6 px-2 mx-auto flex w-full flex-col">
	<div class="relative flex items-center justify-between">
		<div class="pl-6 relative">
			<div
				class="left-0 w-1 h-12 shadow-lg absolute top-1/2 -translate-y-1/2 rounded-full bg-[var(--brand)] shadow-[var(--brand)]/40"
			></div>
			<h1 class="text-4xl font-black tracking-tighter" style="color: var(--text-primary)">
				Security Protocol
			</h1>
			<p class="text-sm font-bold mt-1 tracking-tight opacity-50" style="color: var(--text-muted)">
				Define granular access scopes and RBAC matrices.
			</p>
		</div>
		<Button
			variant="primary"
			size="lg"
			onclick={() => (showCreateModal = true)}
			class="shadow-[var(--brand)]/20"
		>
			Initialize Role
		</Button>
	</div>

	{#if loading}
		<div class="py-24 gap-4 flex flex-col items-center justify-center opacity-40">
			<div
				class="h-12 w-12 animate-spin rounded-full border-4 border-[var(--brand)] border-t-transparent"
			></div>
			<p class="text-xs font-black tracking-widest uppercase">Decrypting Keys...</p>
		</div>
	{:else}
		<div class="md:grid-cols-2 lg:grid-cols-3 gap-6 grid grid-cols-1">
			{#each roles as role}
				<div
					class="role-card p-8 group hover:-translate-y-1 relative flex h-full flex-col overflow-hidden transition-all"
				>
					<div class="mb-6 relative z-10 flex items-start justify-between">
						<div
							class="h-14 w-14 rounded-2xl text-2xl shadow-inner flex items-center justify-center bg-[var(--brand)]/10 text-[var(--brand)]"
						>
							🛡️
						</div>
						<div class="gap-2 flex">
							<a
								href="/admin/roles/{role.id}"
								class="h-10 w-10 rounded-xl hover:text-white hover:shadow-lg flex items-center justify-center bg-[var(--surface-alt)] transition-all hover:bg-[var(--brand)] hover:shadow-[var(--brand)]/30"
								title="Edit Role">✎</a
							>
							<button
								onclick={() => handleDeleteRole(role.id)}
								class="h-10 w-10 rounded-xl hover:bg-red-500 hover:text-white hover:shadow-lg hover:shadow-red-500/30 flex items-center justify-center bg-[var(--surface-alt)] transition-all"
								title="Delete Role">🗑</button
							>
						</div>
					</div>

					<div class="relative z-10 mt-auto">
						<h3 class="text-xl font-black tracking-tight mb-2" style="color: var(--text-primary)">
							{role.name}
						</h3>
						<p
							class="text-xs font-bold leading-relaxed mb-6 line-clamp-2 opacity-50"
							style="color: var(--text-muted)"
						>
							{role.description || 'No specialized description provided for this scope.'}
						</p>

						<div class="pt-6 flex items-center justify-between border-t border-[var(--border)]">
							<div class="gap-0.5 flex flex-col">
								<span
									class="font-black text-[10px] tracking-[0.2em] text-[var(--text-muted)] uppercase opacity-40"
									>Matrix Depth</span
								>
								<span class="text-xs font-black" style="color: var(--text-primary)">
									{role.permissions?.length || 0} Permissions
								</span>
							</div>
							<a
								href="/admin/roles/{role.id}"
								class="h-9 px-4 rounded-xl font-black tracking-widest hover:text-white flex items-center justify-center bg-[var(--surface-alt)] text-[10px] uppercase transition-all hover:bg-[var(--brand)]"
								>Manage</a
							>
						</div>
					</div>

					<!-- Background Accent -->
					<div
						class="-right-12 -bottom-12 w-48 h-48 blur-3xl pointer-events-none absolute rounded-full bg-[var(--brand)]/5 transition-colors group-hover:bg-[var(--brand)]/10"
					></div>
				</div>
			{:else}
				<div
					class="col-span-full py-24 border-2 border-dashed border-[var(--border)] rounded-[2.5rem] flex flex-col items-center justify-center opacity-30 gap-4"
				>
					<span class="text-5xl">🔐</span>
					<div class="text-center">
						<p class="text-lg font-black tracking-tight">Zero Protocols Defined</p>
						<p class="text-xs font-bold uppercase tracking-widest">
							Initialize your first security role above
						</p>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<Modal bind:open={showCreateModal} title="Establish New Security Protocol">
	<form
		onsubmit={(e) => {
			e.preventDefault();
			handleCreateRole();
		}}
		class="gap-6 p-2 flex flex-col"
	>
		<div class="space-y-6">
			<Input
				label="Protocol Identifier"
				bind:value={newRoleName}
				placeholder="e.g. System Architect"
				required
				class="font-bold"
			/>
			<Input
				label="Scope Description"
				bind:value={newRoleDesc}
				placeholder="Describe the operational boundaries of this role..."
				type="textarea"
				class="font-bold"
			/>
		</div>

		<div class="gap-6 pt-4 flex justify-end border-t border-[var(--border)]">
			<button
				type="button"
				onclick={() => (showCreateModal = false)}
				class="font-black tracking-widest text-[10px] uppercase opacity-40 transition-opacity hover:opacity-100"
				>Discard Access</button
			>
			<Button type="submit" loading={creating} class="w-48">Initialize Role</Button>
		</div>
	</form>
</Modal>

<style>
	.role-card {
		background: rgba(var(--surface-rgb), 0.7);
		backdrop-filter: blur(40px);
		border-radius: 2rem;
		box-shadow: var(--shadow-xl);
		border: 1px solid var(--border);
	}

	:global(.dark) .role-card {
		background: rgba(var(--surface-rgb), 0.4);
	}
</style>
