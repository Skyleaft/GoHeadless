<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { adminApi, type Role } from '$lib/api/admin';
	import { collectionsApi } from '$lib/api/collections';
	import { toast } from '$lib/stores/toast';
	import Button from '$lib/shared/Button.svelte';
	import Input from '$lib/shared/Input.svelte';
	import type { Collection } from '$lib/types/collection';

	let role = $state<Role | null>(null);
	let collections = $state<Collection[]>([]);
	let loading = $state(true);
	let saving = $state(false);

	const actions = ['create', 'read', 'update', 'delete'] as const;

	onMount(async () => {
		const id = page.params.id;
		try {
			const [r, colls] = await Promise.all([
				adminApi.getRoles().then((list) => list.find((x) => x.id === id)),
				collectionsApi.list()
			]);
			if (r) role = JSON.parse(JSON.stringify(r)); // Clone to avoid direct store mutation if any
			collections = colls;
		} catch (err: any) {
			toast.error(err.message || 'Failed to load role details');
		} finally {
			loading = false;
		}
	});

	function hasPermission(collSlug: string, action: string) {
		if (!role) return false;
		const perm = role.permissions?.find((p) => p.collection_slug === collSlug);
		return perm?.actions.includes(action) || false;
	}

	function togglePermission(collSlug: string, action: string) {
		if (!role) return;
		if (!role.permissions) role.permissions = [];

		let permIndex = role.permissions.findIndex((p) => p.collection_slug === collSlug);

		if (permIndex === -1) {
			role.permissions.push({ collection_slug: collSlug, actions: [action] });
		} else {
			const perm = role.permissions[permIndex];
			if (perm.actions.includes(action)) {
				perm.actions = perm.actions.filter((a) => a !== action);
				if (perm.actions.length === 0) {
					role.permissions = role.permissions.filter((_, i) => i !== permIndex);
				}
			} else {
				perm.actions.push(action);
			}
		}
	}

	async function handleSave() {
		if (!role || !role.id) return;
		saving = true;
		try {
			await adminApi.updateRole(role.id, role);
			toast.success('Permissions updated successfully');
		} catch (err: any) {
			toast.error(err.message || 'Failed to update role');
		} finally {
			saving = false;
		}
	}
</script>

<div class="gap-6 animate-fade-in pb-20 flex flex-col">
	<div class="gap-4 flex items-center">
		<a
			href="/admin/roles"
			class="h-10 w-10 flex items-center justify-center rounded-full transition-colors hover:bg-[--surface-alt]"
			title="Back to Roles">←</a
		>
		<div class="flex-1">
			<h1 class="text-2xl font-bold">{loading ? 'Loading Role...' : role?.name}</h1>
			<p class="text-sm text-[--text-muted]">Edit permissions and scope</p>
		</div>
		<Button onclick={handleSave} loading={saving} disabled={loading}>Save Permissions</Button>
	</div>

	{#if loading}
		<div class="py-12 flex justify-center">
			<div
				class="h-8 w-8 animate-spin rounded-full border-4 border-[--brand] border-t-transparent"
			></div>
		</div>
	{:else if role}
		<div class="lg:grid-cols-3 gap-6 grid grid-cols-1">
			<!-- Basic Info -->
			<div class="lg:col-span-1">
				<div class="card p-6 flex flex-col gap-4">
					<h2 class="text-sm font-bold uppercase tracking-wider text-[--text-muted]">Role Details</h2>
					<Input id="roleName" label="Name" bind:value={role.name} required />
					<div class="flex flex-col gap-1.5">
						<label for="roleDesc" class="text-sm font-medium">Description</label>
						<textarea
							id="roleDesc"
							bind:value={role.description}
							rows="3"
							class="w-full rounded-lg border border-[--border] bg-[--surface] p-3 text-sm focus:border-[--brand] focus:ring-1 focus:ring-[--brand] outline-none"
							placeholder="Describe what this role is for..."
						></textarea>
					</div>
				</div>
			</div>

			<!-- Permission Matrix -->
			<div class="lg:col-span-2">
				<div class="card overflow-hidden">
					<div class="p-4 border-b border-[--border] bg-[--surface-alt]">
						<h2 class="text-sm font-bold tracking-wider text-[--text-muted] uppercase">
							Collection Permissions
						</h2>
					</div>
					<table class="text-sm w-full text-left">
						<thead>
							<tr class="border-b border-[--border] bg-[--surface-alt]/30">
								<th class="px-6 py-3 font-semibold text-xs text-[--text-secondary]">Collection</th>
								{#each actions as action}
									<th
										class="px-6 py-3 font-semibold text-xs text-center text-[--text-secondary] capitalize"
										>{action}</th
									>
								{/each}
							</tr>
						</thead>
						<tbody class="divide-y divide-[--border]">
							{#each collections as coll}
								<tr class="transition-colors hover:bg-[--surface-alt]/30">
									<td class="px-6 py-4">
										<div class="flex flex-col">
											<span class="font-medium">{coll.name}</span>
											<span class="text-[10px] text-[--text-muted]">{coll.slug}</span>
										</div>
									</td>
									{#each actions as action}
										<td class="px-6 py-4 text-center">
											<button
												onclick={() => togglePermission(coll.slug, action)}
												class="h-6 w-6 rounded inline-flex items-center justify-center border transition-all
													{hasPermission(coll.slug, action)
													? 'text-white shadow-sm border-[--brand] bg-[--brand]'
													: 'border-[--border] bg-[--surface] hover:border-[--brand]'}"
											>
												{#if hasPermission(coll.slug, action)}
													✓
												{/if}
											</button>
										</td>
									{/each}
								</tr>
							{:else}
								<tr>
									<td
										colspan={actions.length + 1}
										class="px-6 py-12 text-center text-[--text-muted]"
									>
										No collections found. Create a collection first to define permissions.
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		</div>
	{/if}
</div>
