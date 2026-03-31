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
	let initialRoleName = $state(''); // For breadcrumb

	const actions = ['create', 'read', 'update', 'delete'] as const;

	onMount(async () => {
		const id = page.params.id;
		try {
			const [r, colls] = await Promise.all([
				adminApi.getRoles().then((list) => list.find((x) => x.id === id)),
				collectionsApi.list()
			]);
			if (r) {
				role = JSON.parse(JSON.stringify(r));
				initialRoleName = role?.name || '';
			}
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

<div class="max-w-6xl animate-fade-in gap-10 py-6 px-2 mx-auto flex w-full flex-col">
	<div class="relative flex items-center justify-between">
		<div class="gap-6 flex items-center">
			<a
				href="/admin/roles"
				class="h-12 w-12 hover:text-white shadow-sm flex items-center justify-center rounded-full border border-[var(--border)] bg-[var(--surface-alt)] transition-all hover:scale-110 hover:bg-[var(--brand)]"
				title="Back to Security Protocols"
			>
				<span class="text-xl">←</span>
			</a>
			<div class="pl-6 relative">
				<div
					class="left-0 w-1 h-12 shadow-lg absolute top-1/2 -translate-y-1/2 rounded-full bg-[var(--brand)] shadow-[var(--brand)]/40"
				></div>
				<h1 class="text-4xl font-black tracking-tighter" style="color: var(--text-primary)">
					{loading ? 'Decrypting Protocol...' : initialRoleName}
				</h1>
				<p
					class="text-sm font-bold mt-1 tracking-tight opacity-50"
					style="color: var(--text-muted)"
				>
					Configure operational boundaries and access matrices.
				</p>
			</div>
		</div>
		<Button
			variant="primary"
			size="lg"
			onclick={handleSave}
			loading={saving}
			disabled={loading}
			class="w-56 shadow-[var(--brand)]/20"
		>
			Commit Protocol
		</Button>
	</div>

	{#if loading}
		<div class="py-32 gap-6 flex flex-col items-center justify-center opacity-40">
			<div class="h-16 w-16 relative">
				<div class="inset-0 rounded-2xl animate-pulse absolute bg-[var(--brand)]/20"></div>
				<div class="inset-0 absolute flex items-center justify-center">
					<div
						class="h-8 w-8 animate-spin rounded-full border-4 border-[var(--brand)] border-t-transparent shadow-[var(--brand)]/20"
					></div>
				</div>
			</div>
			<p class="font-black text-[10px] tracking-[0.3em] uppercase">Accessing Matrix...</p>
		</div>
	{:else if role}
		<div class="lg:grid-cols-12 gap-8 grid grid-cols-1 items-start">
			<!-- Basic Info -->
			<div class="lg:col-span-4 gap-8 flex flex-col">
				<div class="card p-8 gap-6 flex flex-col">
					<div>
						<h2 class="text-xl font-black tracking-tight" style="color: var(--text-primary)">
							Protocol Metadata
						</h2>
						<p
							class="font-black tracking-widest mt-1 text-[10px] uppercase opacity-40"
							style="color: var(--text-muted)"
						>
							Identity & Definition
						</p>
					</div>

					<Input
						id="roleName"
						label="Protocol Title"
						bind:value={role.name}
						required
						class="font-bold"
					/>

					<div class="gap-2 flex flex-col">
						<label
							for="roleDesc"
							class="px-1 font-black text-[10px] tracking-[0.15em] uppercase opacity-40"
							style="color: var(--text-primary)">Operational Boundaries</label
						>
						<textarea
							id="roleDesc"
							bind:value={role.description}
							rows="4"
							class="rounded-2xl p-4 text-sm font-bold w-full border-2 border-transparent bg-[var(--surface-alt)] transition-all focus:border-[var(--brand)]/20 focus:ring-4 focus:ring-[var(--brand)]/10 focus:outline-none"
							style="color: var(--text-primary)"
							placeholder="Describe the operational scope of this protocol..."
						></textarea>
					</div>
				</div>

				<div class="card p-8 to-purple-500/5 group bg-gradient-to-br from-[var(--brand)]/5">
					<div class="gap-2 relative z-10 flex flex-col">
						<span class="text-2xl w-fit opacity-40 transition-transform group-hover:scale-110"
							>🧩</span
						>
						<h3
							class="text-sm font-black tracking-widest uppercase"
							style="color: var(--text-primary)"
						>
							RBAC Inheritance
						</h3>
						<p
							class="text-xs font-bold leading-relaxed opacity-40"
							style="color: var(--text-muted)"
						>
							Permissions defined in this matrix are strictly enforced at the API gateway layer.
						</p>
					</div>
				</div>
			</div>

			<!-- Permission Matrix -->
			<div class="lg:col-span-8">
				<div class="card overflow-hidden">
					<div class="p-8 border-b border-[var(--border)] bg-[var(--surface-alt)]/20">
						<h2 class="text-xl font-black tracking-tight" style="color: var(--text-primary)">
							Security Matrix
						</h2>
						<p
							class="font-black tracking-widest mt-1 text-[10px] uppercase opacity-40"
							style="color: var(--text-muted)"
						>
							CRUD Authorization Nodes
						</p>
					</div>
					<table class="text-sm w-full border-collapse text-left">
						<thead>
							<tr class="border-b border-[var(--border)] bg-[var(--surface-alt)]/30">
								<th
									class="px-8 py-5 font-black text-[10px] tracking-[0.2em] uppercase opacity-40"
									style="color: var(--text-muted)">Identity Node</th
								>
								{#each actions as action}
									<th
										class="px-8 py-5 font-black text-center text-[10px] tracking-[0.2em] uppercase opacity-40"
										style="color: var(--text-muted)">{action}</th
									>
								{/each}
							</tr>
						</thead>
						<tbody class="divide-y divide-[var(--border)]">
							{#each collections as coll}
								<tr class="group transition-colors hover:bg-[var(--surface-alt)]/20">
									<td class="px-8 py-6">
										<div class="gap-4 flex items-center">
											<div
												class="h-10 w-10 rounded-xl font-black shadow-inner flex items-center justify-center bg-[var(--brand)]/10 text-[var(--brand)] transition-transform group-hover:scale-105"
											>
												{coll.name[0].toUpperCase()}
											</div>
											<div class="flex flex-col">
												<span
													class="text-base font-black tracking-tight"
													style="color: var(--text-primary)">{coll.name}</span
												>
												<span
													class="font-black tracking-widest text-[9px] leading-none uppercase opacity-30"
													>{coll.slug}</span
												>
											</div>
										</div>
									</td>
									{#each actions as action}
										<td class="px-8 py-6 text-center">
											<button
												onclick={() => togglePermission(coll.slug, action)}
												title={`${hasPermission(coll.slug, action) ? 'Revoke' : 'Grant'} ${action} rights`}
												class="h-9 w-9 rounded-xl inline-flex items-center justify-center border-2 transition-all duration-200 active:scale-90
													{hasPermission(coll.slug, action)
													? 'text-white shadow-lg border-[var(--brand)] bg-[var(--brand)] shadow-[var(--brand)]/20'
													: 'border-transparent bg-[var(--surface-alt)] hover:border-[var(--brand)]/40 hover:bg-[var(--surface-alt)]/80'}"
											>
												{#if hasPermission(coll.slug, action)}
													<span class="text-base font-black">✓</span>
												{:else}
													<span
														class="text-xs font-black opacity-20 transition-opacity group-hover:opacity-40"
														>○</span
													>
												{/if}
											</button>
										</td>
									{/each}
								</tr>
							{:else}
								<tr>
									<td colspan={actions.length + 1} class="px-8 py-24 text-center opacity-40">
										<div class="flex flex-col items-center gap-4">
											<span class="text-4xl">📁</span>
											<div>
												<p
													class="text-base font-black tracking-tight"
													style="color: var(--text-primary)"
												>
													Zero Entities Detected
												</p>
												<p class="text-xs font-bold uppercase tracking-widest leading-relaxed">
													Initialize a collection to map permissions
												</p>
											</div>
										</div>
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

<style>
	.card {
		background: rgba(var(--surface-rgb), 0.7);
		backdrop-filter: blur(40px);
		border-radius: 2.5rem;
		box-shadow: var(--shadow-xl);
		border: 1px solid var(--border);
	}

	:global(.dark) .card {
		background: rgba(var(--surface-rgb), 0.4);
	}
</style>
