<script lang="ts">
	import type { Collection, Field } from '$lib/types/collection';
	import FieldEditor from './FieldEditor.svelte';
	import Button from '$lib/shared/Button.svelte';

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
			{ key: '', label: '', type: 'text', required: false, unique: false }
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
		if (!/^[a-z0-9-]+$/.test(slug)) e.slug = 'Slug must be lowercase letters, numbers, and hyphens only';
		if (fields.length === 0) e.fields = 'Add at least one field';
		const emptyKey = fields.some((f) => !f.key.trim());
		if (emptyKey) e.fields = 'All fields must have a key';
		errors = e;
		return Object.keys(e).length === 0;
	}

	function handleSubmit(e: Event) {
		e.preventDefault();
		if (!validate()) return;
		onsubmit?.({ name, slug, description, fields });
	}
</script>

<form onsubmit={handleSubmit} class="flex flex-col gap-8">
	<!-- Metadata section -->
	<div class="card">
		<div class="card-header">
			<h2 class="text-base font-semibold" style="color: var(--text-primary)">Collection Info</h2>
			<p class="text-sm mt-0.5" style="color: var(--text-muted)">
				Define the name and slug for this collection.
			</p>
		</div>
		<div class="card-body grid grid-cols-1 gap-5 sm:grid-cols-2">
			<div class="flex flex-col gap-1.5">
				<label class="text-sm font-medium" style="color: var(--text-primary)">
					Name <span class="text-red-500">*</span>
				</label>
				<input
					bind:value={name}
					placeholder="Blog Posts"
					class="h-9 rounded-lg border px-3 text-sm w-full"
					style="background: var(--surface); border-color: {errors.name ? '#ef4444' : 'var(--border)'}; color: var(--text-primary)"
				/>
				{#if errors.name}<p class="text-xs text-red-500">{errors.name}</p>{/if}
			</div>

			<div class="flex flex-col gap-1.5">
				<label class="text-sm font-medium" style="color: var(--text-primary)">
					Slug <span class="text-red-500">*</span>
				</label>
				<div class="relative">
					<span
						class="absolute left-3 top-1/2 -translate-y-1/2 text-sm font-mono"
						style="color: var(--text-muted)"
					>/</span>
					<input
						bind:value={slug}
						oninput={() => (slugEdited = true)}
						placeholder="blog-posts"
						class="h-9 rounded-lg border pl-6 pr-3 text-sm font-mono w-full"
						style="background: var(--surface); border-color: {errors.slug ? '#ef4444' : 'var(--border)'}; color: var(--text-primary)"
					/>
				</div>
				{#if errors.slug}<p class="text-xs text-red-500">{errors.slug}</p>{/if}
			</div>

			<div class="col-span-full flex flex-col gap-1.5">
				<label class="text-sm font-medium" style="color: var(--text-primary)">Description</label>
				<textarea
					bind:value={description}
					placeholder="Optional description of this collection…"
					rows="2"
					class="rounded-lg border px-3 py-2 text-sm w-full resize-none"
					style="background: var(--surface); border-color: var(--border); color: var(--text-primary)"
				></textarea>
			</div>
		</div>
	</div>

	<!-- Fields section -->
	<div class="card">
		<div class="card-header flex items-center justify-between">
			<div>
				<h2 class="text-base font-semibold" style="color: var(--text-primary)">Schema Fields</h2>
				<p class="text-sm mt-0.5" style="color: var(--text-muted)">
					Define the fields for this collection's data structure.
				</p>
			</div>
			<Button variant="primary" size="sm" onclick={addField} type="button">+ Add Field</Button>
		</div>

		<div class="card-body flex flex-col gap-3">
			{#if errors.fields}
				<p class="rounded-lg px-3 py-2 text-sm text-red-700" style="background: #fef2f2">{errors.fields}</p>
			{/if}

			{#if fields.length === 0}
				<div class="flex flex-col items-center gap-3 py-12 text-center">
					<span class="text-4xl">🏗</span>
					<p class="text-sm" style="color: var(--text-muted)">
						No fields yet. Click "Add Field" to define your schema.
					</p>
				</div>
			{:else}
				{#each fields as field, i}
					<FieldEditor
						bind:field={fields[i]}
						onremove={() => removeField(i)}
					/>
				{/each}
			{/if}
		</div>
	</div>

	<!-- Submit -->
	<div class="flex items-center justify-end gap-3">
		<a href="/collections" class="inline-flex h-9 items-center px-4 text-sm rounded-lg transition hover:bg-[--surface-alt]" style="color: var(--text-secondary)">
			Cancel
		</a>
		<Button type="submit" variant="primary" {loading}>
			{loading ? 'Saving…' : '💾 Save Collection'}
		</Button>
	</div>
</form>
