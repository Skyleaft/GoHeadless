<script lang="ts">
	import type { Collection, Field } from '$lib/types/collection';
	import FieldEditor from './FieldEditor.svelte';
	import Button from '$lib/shared/Button.svelte';
	import Input from '$lib/shared/Input.svelte';
	import { adminApi, type Role } from '$lib/api/admin';
	import { onMount } from 'svelte';

	interface Props {
		initial?: Partial<Collection>;
		loading?: boolean;
		isEditMode?: boolean;
		onsubmit?: (collection: Omit<Collection, 'id'>) => void;
	}

	let { initial, loading = false, isEditMode = false, onsubmit }: Props = $props();

	let name = $state(initial?.name ?? '');
	let slug = $state(initial?.slug ?? '');
	let description = $state(initial?.description ?? '');
	let fields = $state<Field[]>(initial?.fields ?? []);
	let access = $state(
		initial?.access ?? {
			is_public: false,
			allowed_roles: [],
			crud_policy: { create: [], read: [], update: [], delete: [] }
		}
	);
	let slugEdited = $state(!!initial?.slug);
	let roles = $state<Role[]>([]);

	// Search state for each action
	let searchTerms = $state({ create: '', read: '', update: '', delete: '' });
	let openDropdown = $state<string | null>(null);

	onMount(async () => {
		try { roles = await adminApi.getRoles(); } catch {}
	});

	// Auto-generate slug from name
	$effect(() => {
		if (!slugEdited && name) {
			slug = name
				.toLowerCase()
				.replace(/\s+/g, '-')
				.replace(/[^a-z0-9-]/g, '');
		}
	});

	function addField() {
		fields = [
			...fields,
			{
				key: '',
				label: '',
				type: 'text',
				required: false,
				unique: false,
				is_array: false,
				array_config: { unique_items: false },
				searchable: true,
				internal: false
			}
		];
	}

	function removeField(i: number) {
		fields = fields.filter((_, idx) => idx !== i);
	}

	let errors = $state<{ name?: string; slug?: string; fields?: string }>({});

	function validate() {
		const e: typeof errors = {};
		if (!name.trim()) e.name = 'Name is required';
		if (!slug.trim()) e.slug = 'Slug is required';
		if (!/^[a-z0-9-]+$/.test(slug))
			e.slug = 'Slug must be lowercase letters, numbers, and hyphens only';
		if (fields.length === 0) e.fields = 'Add at least one field';
		const emptyKey = fields.some((f) => !f.key.trim());
		if (emptyKey) e.fields = 'All fields must have a key';
		errors = e;
		return Object.keys(e).length === 0;
	}

	function handleSubmit(e: Event) {
		e.preventDefault();
		if (!validate()) return;
		onsubmit?.({ name, slug, description, fields, access });
	}

	function getPolicyForAction(action: string): string[] {
		return access.crud_policy?.[action as keyof typeof access.crud_policy] ?? [];
	}

	function toggleRole(action: string, roleId: string) {
		if (!access.crud_policy) {
			access.crud_policy = { create: [], read: [], update: [], delete: [] };
		}
		const key = action as keyof typeof access.crud_policy;
		const policy = access.crud_policy[key];
		if (policy.includes(roleId)) {
			access.crud_policy = { ...access.crud_policy, [key]: policy.filter((id) => id !== roleId) };
		} else {
			access.crud_policy = { ...access.crud_policy, [key]: [...policy, roleId] };
		}
	}

	function selectAllRoles(action: string) {
		if (!access.crud_policy) {
			access.crud_policy = { create: [], read: [], update: [], delete: [] };
		}
		const key = action as keyof typeof access.crud_policy;
		access.crud_policy = { ...access.crud_policy, [key]: roles.map((r) => r.id ?? '') };
	}

	function clearAllRoles(action: string) {
		if (!access.crud_policy) {
			access.crud_policy = { create: [], read: [], update: [], delete: [] };
		}
		const key = action as keyof typeof access.crud_policy;
		access.crud_policy = { ...access.crud_policy, [key]: [] };
	}

	function getFilteredRoles(action: string): Role[] {
		const term = searchTerms[action as keyof typeof searchTerms]?.toLowerCase() ?? '';
		if (!term) return roles;
		return roles.filter((r) => r.name.toLowerCase().includes(term));
	}

	function getRoleNameById(id: string): string {
		return roles.find((r) => r.id === id)?.name ?? id;
	}

	const actionConfigs = [
		{ key: 'create', icon: 'M12 4.5v15m7.5-7.5h-15', color: '#10b981' },
		{ key: 'read', icon: 'M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z', color: '#3b82f6' },
		{ key: 'update', icon: 'm16.862 4.487 1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125', color: '#f59e0b' },
		{ key: 'delete', icon: 'm14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0', color: '#ef4444' }
	];
	
	// 🔥 click outside + focus leave action
	function dropdown(node: HTMLElement, key: string) {
		const handleClick = (e: MouseEvent) => {
			if (!node.contains(e.target as Node)) {
				if (openDropdown === key) openDropdown = null;
			}
		};

		const handleFocusOut = (e: FocusEvent) => {
			const next = e.relatedTarget as Node | null;
			if (!node.contains(next)) {
				if (openDropdown === key) openDropdown = null;
			}
		};

		const handleKey = (e: KeyboardEvent) => {
			if (e.key === 'Escape' && openDropdown === key) {
				openDropdown = null;
			}
		};

		document.addEventListener('click', handleClick);
		node.addEventListener('focusout', handleFocusOut);
		document.addEventListener('keydown', handleKey);

		return {
			destroy() {
				document.removeEventListener('click', handleClick);
				node.removeEventListener('focusout', handleFocusOut);
				document.removeEventListener('keydown', handleKey);
			}
		};
	}
</script>

<form onsubmit={handleSubmit} class="gap-10 flex flex-col">
	<!-- Metadata section -->
	<div class="card overflow-hidden">
		<div class="p-8 border-b border-[--border] bg-[--surface-alt]/20">
			<h2 class="text-xl font-black tracking-tight" style="color: var(--text-primary)">
				Architect Identity
			</h2>
			<p
				class="text-xs font-bold mt-1 tracking-wide uppercase opacity-40"
				style="color: var(--text-muted)"
			>
				System-level identifiers and indexing schemas.
			</p>
		</div>
		<div class="p-8 gap-8 sm:grid-cols-2 grid grid-cols-1">
			<Input
				label="Natural Descriptor"
				placeholder="e.g., Blog Articles"
				bind:value={name}
				error={errors.name}
				required
				disabled={isEditMode}
			/>

			<div class="gap-1.5 flex flex-col">
				<label
					class="font-black px-1 text-[10px] tracking-[0.15em] uppercase opacity-50"
					style="color: var(--text-primary)"
				>
					URI Access Link <span class="ml-1 text-[var(--brand)]">●</span>
				</label>
				<div class="relative">
					<span
						class="left-4 text-sm font-black absolute top-1/2 -translate-y-1/2 opacity-30"
						style="color: var(--text-muted)">/</span
					>
					<input
						bind:value={slug}
						oninput={() => (slugEdited = true)}
						placeholder="blog-articles"
						disabled={isEditMode}
						class="h-12 rounded-2xl pl-8 pr-4 text-sm font-bold w-full border-2 transition-all duration-200 focus:ring-4 focus:ring-[var(--brand)]/10 focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed"
						style="background: var(--surface-alt); border-color: {errors.slug
							? '#ef4444'
							: 'transparent'}; color: var(--text-primary)"
					/>
				</div>
				{#if errors.slug}<p class="text-xs font-bold text-red-500 px-1">{errors.slug}</p>{/if}
			</div>

			<div class="col-span-full">
				<Input
					label="Schema Purpose"
					placeholder="Describe the context of this data structure..."
					bind:value={description}
					type="textarea"
				/>
			</div>
		</div>
	</div>

	<!-- Fields section -->
	<div class="card overflow-hidden">
		<div
			class="p-8 flex items-center justify-between border-b border-[--border] bg-[--surface-alt]/20"
		>
			<div>
				<h2 class="text-xl font-black tracking-tight" style="color: var(--text-primary)">
					Data Blueprint
				</h2>
				<p
					class="text-xs font-bold mt-1 tracking-wide uppercase opacity-40"
					style="color: var(--text-muted)"
				>
					Define granular field types and validation rules.
				</p>
			</div>
			<Button variant="primary" onclick={addField} type="button" class="shadow-[var(--brand)]/40">
				Add Component
			</Button>
		</div>

		<div class="p-8 gap-4 flex flex-col">
			{#if errors.fields}
				<div
					class="rounded-2xl px-4 py-3 text-xs font-black tracking-widest text-red-500 bg-red-500/5 border-red-500/20 mb-4 border uppercase"
				>
					{errors.fields}
				</div>
			{/if}

			{#if fields.length === 0}
				<div
					class="gap-4 py-16 rounded-3xl flex flex-col items-center border-2 border-dashed border-[--border] text-center opacity-40"
				>
					<span class="text-5xl grayscale">🧊</span>
					<div class="space-y-1">
						<p class="text-base font-black tracking-tighter" style="color: var(--text-primary)">
							Zero Nodes Detected
						</p>
						<p class="text-xs font-bold tracking-widest uppercase" style="color: var(--text-muted)">
							Initialize your blueprint above
						</p>
					</div>
				</div>
			{:else}
				<div class="gap-4 grid grid-cols-1">
					{#each fields as field, i (i)}
						<div class="animate-fade-in">
							<FieldEditor bind:field={fields[i]} onremove={() => removeField(i)} />
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>

	<!-- Access Control section -->
	<div class="card">
		<div class="p-8 border-b border-[--border] bg-[--surface-alt]/20">
			<h2 class="text-xl font-black tracking-tight" style="color: var(--text-primary)">
				Security & Exposure
			</h2>
			<p
				class="text-xs font-bold mt-1 tracking-wide uppercase opacity-40"
				style="color: var(--text-muted)"
			>
				RBAC matrices and public endpoint accessibility.
			</p>
		</div>
		<div class="p-8 gap-6 flex flex-col">
			<!-- Public toggle -->
			<div
				class="p-6 shadow-inner hover:shadow-md group flex items-center justify-between rounded-[2rem] bg-[var(--bg)] transition-all"
			>
				<div class="gap-1 flex flex-col">
					<label
						for="isPublic"
						class="text-sm font-black tracking-tight uppercase"
						style="color: var(--text-primary)">Public Exposure (GET)</label
					>
					<p class="text-xs font-bold tracking-tight opacity-40" style="color: var(--text-muted)">
						Allow anonymous read-access to this collection's content.
					</p>
				</div>
				<label class="pr-2 relative inline-flex cursor-pointer items-center">
					<input
						id="isPublic"
						type="checkbox"
						bind:checked={access.is_public}
						class="peer sr-only"
					/>
					<div
						class="peer h-7 w-12 bg-slate-300 dark:bg-slate-700 after:h-5 after:w-5 after:bg-white after:shadow-md peer-checked:after:translate-x-5 rounded-full transition-all peer-checked:bg-[var(--brand)] peer-focus:ring-4 peer-focus:ring-[var(--brand)]/20 after:absolute after:top-[4px] after:left-[4px] after:rounded-full after:transition-all after:content-['']"
					></div>
				</label>
			</div>

			<!-- Permission Matrix -->
			<div
				class="rounded-3xl border border-[var(--border)]/50 bg-[var(--bg)]/40"
			>
				<div class="px-6 py-4 bg-[var(--surface-alt)]/30">
					<div class="flex items-center gap-3">
						<div
							class="h-8 w-8 rounded-xl flex items-center justify-center bg-[var(--brand)]/10"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-4 w-4"
								style="color: var(--brand)"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z"
								/>
							</svg>
						</div>
						<div>
							<h3 class="text-sm font-black tracking-tight uppercase" style="color: var(--text-primary)">
								Permission Matrix
							</h3>
							<p class="text-[10px] font-bold tracking-wider uppercase opacity-40" style="color: var(--text-muted)">
								Search and assign roles to each action
							</p>
						</div>
					</div>
				</div>

				<div class="p-6 grid grid-cols-1 lg:grid-cols-2 gap-6">
					{#each actionConfigs as actionDef (actionDef.key)}
						{@const selectedRoles = getPolicyForAction(actionDef.key)}
						{@const filteredRoles = getFilteredRoles(actionDef.key)}
						{@const isOpen = openDropdown === actionDef.key}
						<div 
						use:dropdown={actionDef.key} class="flex flex-col gap-3 relative"
							class:z-[50]={isOpen}>
							<!-- Header -->
							<div class="flex items-center justify-between">
								<div class="flex items-center gap-2.5">
									<div
										class="h-7 w-7 rounded-lg flex items-center justify-center"
										style="background: {actionDef.color}15"
									>
										<svg
											xmlns="http://www.w3.org/2000/svg"
											class="h-3.5 w-3.5"
											style="color: {actionDef.color}"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											stroke-width="2.5"
											stroke-linecap="round"
											stroke-linejoin="round"
										>
											<path d={actionDef.icon} />
										</svg>
									</div>
									<span
										class="text-xs font-black tracking-[0.15em] uppercase"
										style="color: {actionDef.color}"
									>
										{actionDef.key.toUpperCase()}
									</span>
								</div>
								<div class="flex items-center gap-2">
									<span class="text-[10px] font-bold opacity-40 uppercase" style="color: var(--text-muted)">
										{selectedRoles.length} role{selectedRoles.length !== 1 ? 's' : ''}
									</span>
									{#if selectedRoles.length > 0}
										<button
											type="button"
											class="text-[10px] font-bold uppercase tracking-wide transition-colors hover:opacity-100 opacity-50"
											style="color: var(--brand)"
											onclick={() => clearAllRoles(actionDef.key)}
										>
											Clear
										</button>
									{/if}
								</div>
							</div>

							<!-- Selected roles chips -->
							<div class="flex flex-wrap gap-1.5 min-h-[2.5rem]">
								{#if selectedRoles.length === 0}
									<span
										class="text-xs italic opacity-40 py-1"
										style="color: var(--text-muted)"
									>
										No roles assigned
									</span>
								{:else}
									{#each selectedRoles as roleId}
										<span
											class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs font-semibold transition-all"
											style="background: {actionDef.color}15; color: {actionDef.color}"
										>
											{getRoleNameById(roleId)}
											<button
												type="button"
												class="ml-0.5 opacity-60 hover:opacity-100 transition-opacity"
												onclick={() => toggleRole(actionDef.key, roleId)}
											>
												<svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
													<path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
												</svg>
											</button>
										</span>
									{/each}
								{/if}
							</div>

							<!-- Search and dropdown -->
							<div class="relative">
								<div class="relative">
									<svg
										xmlns="http://www.w3.org/2000/svg"
										class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 opacity-30"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
									>
										<circle cx="11" cy="11" r="8" />
										<path d="m21 21-4.35-4.35" />
									</svg>
									<input
										type="text"
										placeholder="Search roles..."
										value={searchTerms[actionDef.key as keyof typeof searchTerms] ?? ''}
										oninput={(e) => {
											searchTerms = { ...searchTerms, [actionDef.key]: (e.target as HTMLInputElement).value };
											openDropdown = actionDef.key;
										}}
										onfocus={() => { openDropdown = actionDef.key; }}
										class="w-full h-10 pl-9 pr-4 rounded-xl text-xs font-medium border transition-all focus:outline-none"
										style="background: var(--surface-alt); border-color: var(--border); color: var(--text-primary)"
										autocomplete="off"
									/>
								</div>

								{#if isOpen && roles.length > 0}
									<div
										class=" mt-1.5 w-full max-h-48 overflow-auto rounded-xl border shadow-2xl"
										style="background: var(--bg); border-color: var(--border)"
									>
										{#if filteredRoles.length === 0}
											<div class="px-3 py-3 text-xs text-center opacity-50" style="color: var(--text-muted)">
												No roles found
											</div>
										{:else}
											{#if filteredRoles.length > 5}
												<button
													type="button"
													class="w-full px-3 py-2 text-xs font-bold uppercase tracking-wide text-left transition-colors border-b"
													style="color: var(--brand); border-color: var(--border)/30; background: var(--surface-alt)/30"
													onclick={() => selectAllRoles(actionDef.key)}
												>
													+ Select all ({filteredRoles.length})
												</button>
											{/if}
											{#each filteredRoles as role}
												{@const isSelected = selectedRoles.includes(role.id ?? '')}
												<label
													class="flex items-center gap-2.5 px-3 py-2 cursor-pointer transition-colors hover:bg-[var(--surface-alt)]/50"
													style="border-bottom: 1px solid var(--border)/10"
												>
													<div
														class="h-4 w-4 rounded-md border-2 flex items-center justify-center flex-shrink-0 transition-all"
														style="border-color: {isSelected ? 'var(--brand)' : 'var(--border)'}; background: {isSelected ? 'var(--brand)' : 'transparent'}"
													>
														{#if isSelected}
															<svg class="h-3 w-3 text-white" viewBox="0 0 20 20" fill="currentColor">
																<path fill-rule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clip-rule="evenodd" />
															</svg>
														{/if}
													</div>
													<input
														type="checkbox"
														class="sr-only"
														checked={isSelected}
														onchange={() => toggleRole(actionDef.key, role.id ?? '')}
													/>
													<span class="text-xs font-medium truncate" style="color: var(--text-secondary)">
														{role.name}
													</span>
												</label>
											{/each}
										{/if}
									</div>
								{/if}
							</div>
						</div>
					{/each}
				</div>

				{#if roles.length === 0}
					<div class="px-6 py-8 text-center">
						<p class="text-xs font-medium italic opacity-40" style="color: var(--text-muted)">
							No roles configured yet. Create roles first to assign permissions.
						</p>
					</div>
				{/if}
			</div>
		</div>
	</div>

	<!-- Submit -->
	<div class="gap-6 pb-12 pt-4 px-4 flex items-center justify-end">
		<a
			href="/collections"
			class="text-sm font-black tracking-widest uppercase opacity-40 transition-opacity hover:opacity-100"
			style="color: var(--text-secondary)"
		>
			Discard Suite
		</a>
		<Button type="submit" variant="primary" {loading} class="w-64">
			{#if loading}
				Archiving...
			{:else}
				Finalize & Build
			{/if}
		</Button>
	</div>
</form>

<style>
	.card {
		background: rgba(var(--surface-rgb), 0.7);
		backdrop-filter: blur(40px);
		border-radius: 2.5rem;
		box-shadow: var(--shadow-xl);
		border: 1px solid var(--border);
	}
</style>