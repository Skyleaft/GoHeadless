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

<div class="max-w-6xl mx-auto w-full animate-fade-in flex flex-col gap-10 py-6 px-2">
	<div class="flex items-center justify-between relative">
		<div class="relative pl-6">
			<div class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-12 rounded-full bg-[var(--brand)] shadow-lg shadow-[var(--brand)]/40"></div>
			<h1 class="text-4xl font-black tracking-tighter" style="color: var(--text-primary)">
				User Registry
			</h1>
			<p class="text-sm font-bold mt-1 opacity-50 tracking-tight" style="color: var(--text-muted)">
				Manage system access for administrators and architects.
			</p>
		</div>
		<Button variant="primary" size="lg" onclick={() => (showCreateModal = true)} class="shadow-[var(--brand)]/20">
			Provision User
		</Button>
	</div>

	{#if loading}
		<div class="flex flex-col items-center justify-center py-24 gap-4 opacity-40">
			<div class="h-12 w-12 rounded-full border-4 border-[var(--brand)] border-t-transparent animate-spin"></div>
			<p class="text-xs font-black uppercase tracking-widest">Syncing Identity Hub...</p>
		</div>
	{:else}
		<div class="registry-card overflow-hidden">
			<table class="w-full text-left text-sm border-collapse">
				<thead>
					<tr class="bg-[var(--surface-alt)]/30 border-b border-[var(--border)]">
						<th class="px-8 py-5 text-[10px] font-black uppercase tracking-[0.2em] opacity-40" style="color: var(--text-muted)">Identity Node</th>
						<th class="px-8 py-5 text-[10px] font-black uppercase tracking-[0.2em] opacity-40" style="color: var(--text-muted)">Assigned Rights</th>
						<th class="px-8 py-5 text-[10px] font-black uppercase tracking-[0.2em] opacity-40" style="color: var(--text-muted)">System Integrity</th>
						<th class="px-8 py-5 text-[10px] font-black uppercase tracking-[0.2em] opacity-40 text-right" style="color: var(--text-muted)">Operations</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-[var(--border)]">
					{#each users as user}
						<tr class="hover:bg-[var(--surface-alt)]/20 transition-colors group">
							<td class="px-8 py-6">
								<div class="flex items-center gap-4">
									<div class="h-12 w-12 rounded-2xl font-black text-white shadow-xl flex items-center justify-center transition-transform group-hover:scale-105"
										style="background: linear-gradient(135deg, var(--brand), #7c3aed)">
										{user.username[0].toUpperCase()}
									</div>
									<div class="flex flex-col gap-0.5">
										<span class="text-base font-black tracking-tight" style="color: var(--text-primary)">{user.username}</span>
										{#if user.is_initial_admin}
											<span class="text-[9px] font-black uppercase tracking-[0.2em] text-amber-500 bg-amber-500/10 px-2 py-0.5 rounded-full w-fit">Superadmin</span>
										{/if}
									</div>
								</div>
							</td>
							<td class="px-8 py-6">
								<div class="flex flex-col gap-0.5">
									<span class="text-xs font-black tracking-tight" style="color: var(--text-secondary)">
										{roles.find(r => r.id === user.role_id)?.name || 'Unknown Protocol'}
									</span>
									<span class="text-[9px] font-bold opacity-30 uppercase tracking-widest leading-none">Security Scope</span>
								</div>
							</td>
							<td class="px-8 py-6">
								{#if user.active_status}
									<span class="inline-flex items-center gap-2 px-3 py-1 rounded-xl text-[10px] font-black uppercase tracking-widest bg-emerald-500/10 text-emerald-500 border border-emerald-500/20 shadow-sm shadow-emerald-500/5">
										<span class="h-1.5 w-1.5 rounded-full bg-emerald-500 animate-pulse"></span> Authorized
									</span>
								{:else}
									<span class="inline-flex items-center gap-2 px-3 py-1 rounded-xl text-[10px] font-black uppercase tracking-widest bg-slate-500/10 text-slate-500 border border-slate-500/20 shadow-sm">
										<span class="h-1.5 w-1.5 rounded-full bg-slate-400"></span> Restricted
									</span>
								{/if}
							</td>
							<td class="px-8 py-6 text-right">
								<button 
									onclick={() => handleDeleteUser(user.id)}
									disabled={user.is_initial_admin}
									class="h-10 w-10 rounded-xl bg-[var(--surface-alt)] flex items-center justify-center transition-all hover:bg-red-500 hover:text-white hover:shadow-lg hover:shadow-red-500/30 disabled:opacity-20 disabled:cursor-not-allowed group/btn"
									title={user.is_initial_admin ? "Cannot delete initial admin" : "Purge User Identity"}
								>
									<span class="text-sm group-hover/btn:scale-110 transition-transform">🗑</span>
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<Modal bind:open={showCreateModal} title="Provision New Identity Node">
	<form onsubmit={(e) => { e.preventDefault(); handleCreateUser(); }} class="flex flex-col gap-8 p-2">
		<div class="space-y-6">
			<Input id="newUsername" label="Identity Identifier" bind:value={newUsername} placeholder="e.g. architect_prime" required class="font-bold" />
			<Input id="newPassword" label="Secure Access Key" type="password" bind:value={newPassword} placeholder="••••••••" required class="font-bold" />
			
			<div class="flex flex-col gap-2">
				<label for="newRoleId" class="px-1 text-[10px] font-black uppercase tracking-[0.15em] opacity-40" style="color: var(--text-primary)">Security Clearance <span class="text-[var(--brand)] ml-1">●</span></label>
				<div class="relative">
					<select 
						id="newRoleId"
						bind:value={newRoleId}
						class="h-12 w-full rounded-2xl px-4 text-sm font-bold bg-[var(--surface-alt)] border-2 border-transparent transition-all duration-200 focus:outline-none focus:ring-4 focus:ring-[var(--brand)]/10 focus:border-[var(--brand)]/20 appearance-none"
						style="color: var(--text-primary)"
					>
						<option value="">Select operational scope...</option>
						{#each roles as role}
							<option value={role.id}>{role.name}</option>
						{/each}
					</select>
					<div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none opacity-40 font-black text-xs">▼</div>
				</div>
			</div>
		</div>

		<div class="flex justify-end gap-6 pt-4 border-t border-[var(--border)]">
			<button 
				type="button"
				onclick={() => (showCreateModal = false)}
				class="text-[10px] font-black uppercase tracking-widest opacity-40 hover:opacity-100 transition-opacity"
			>Discard Session</button>
			<Button type="submit" loading={creating} class="w-48 shadow-[var(--brand)]/20">Finalize Identity</Button>
		</div>
	</form>
</Modal>

<style>
	.registry-card {
		background: rgba(var(--surface-rgb), 0.7);
		backdrop-filter: blur(40px);
		border-radius: 2.5rem;
		box-shadow: var(--shadow-xl);
		border: 1px solid var(--border);
	}
	
	:global(.dark) .registry-card {
		background: rgba(var(--surface-rgb), 0.4);
	}
</style>
