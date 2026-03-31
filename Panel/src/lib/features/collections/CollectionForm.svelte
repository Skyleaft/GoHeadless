<script lang="ts">
	import type { Collection, Field } from '$lib/types/collection';
	import FieldEditor from './FieldEditor.svelte';
	import Button from '$lib/shared/Button.svelte';
	import Input from '$lib/shared/Input.svelte';

	interface Props {
		initial?: Partial<Collection>;
		loading?: boolean;
		onsubmit?: (collection: Omit<Collection, 'id'>) => void;
	}

	let { initial, loading = false, onsubmit }: Props = $props();

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
				array_config: { unique_items: false }
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
						class="h-12 rounded-2xl pl-8 pr-4 text-sm font-bold w-full border-2 transition-all duration-200 focus:ring-4 focus:ring-[var(--brand)]/10 focus:outline-none"
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
	<div class="card overflow-hidden">
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
		<div class="p-8 gap-8 flex flex-col">
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

			<div class="gap-6 sm:grid-cols-2 lg:grid-cols-4 grid grid-cols-1">
				{#each ['create', 'read', 'update', 'delete'] as action}
					<div
						class="gap-4 p-5 rounded-3xl flex flex-col border-2 border-transparent bg-[var(--surface-alt)]/30 transition-all hover:border-[var(--brand)]/20 hover:bg-[var(--surface-alt)]/50"
					>
						<div class="space-y-1">
							<span
								class="font-black text-[10px] tracking-[0.2em] uppercase"
								style="color: var(--brand)">{action}</span
							>
							<p class="text-xs font-bold tracking-widest uppercase opacity-40">{action} Rights</p>
						</div>
						<input
							type="text"
							placeholder="Roles..."
							value={access.crud_policy?.[action as keyof typeof access.crud_policy]?.join(', ') ??
								''}
							oninput={(e) => {
								const val = (e.target as HTMLInputElement).value;
								const roles = val
									.split(',')
									.map((r) => r.trim())
									.filter(Boolean);
								if (!access.crud_policy)
									access.crud_policy = { create: [], read: [], update: [], delete: [] };
								// @ts-ignore
								access.crud_policy[action] = roles;
							}}
							class="h-10 px-4 text-xs font-black tracking-tight rounded-xl w-full border-none bg-[var(--bg)] focus:ring-2 focus:ring-[var(--brand)]/40"
							style="color: var(--text-primary)"
						/>
					</div>
				{/each}
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
