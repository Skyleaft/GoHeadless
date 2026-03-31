<script lang="ts">
	import type { Field, ContentRecord } from '$lib/types/collection';
	import FieldRenderer from './FieldRenderer.svelte';
	import Button from '$lib/shared/Button.svelte';

	interface Props {
		fields: Field[];
		initial?: ContentRecord;
		loading?: boolean;
		collectionSlug: string;
		onsubmit?: (data: ContentRecord) => void;
	}

	let { fields, initial, loading = false, collectionSlug, onsubmit }: Props = $props();

	// Initialize record data, pre-filling defaults
	let data = $state<ContentRecord>(
		initial ?? Object.fromEntries(
			fields.map((f) => [f.key, f.default_value ?? (f.type === 'bool' || f.type === 'toggle' ? false : '')])
		)
	);

	function handleSubmit(e: Event) {
		e.preventDefault();
		onsubmit?.(data);
	}
</script>

<form onsubmit={handleSubmit} class="flex flex-col gap-6">
	<div class="card">
		<div class="card-body flex flex-col gap-5">
			{#each fields as field}
				<FieldRenderer {field} bind:data />
			{/each}

			{#if fields.length === 0}
				<p class="py-8 text-center text-sm" style="color: var(--text-muted)">
					This collection has no fields defined yet.
					<a href="/collections" class="text-[--brand] hover:underline">Manage collections ↗</a>
				</p>
			{/if}
		</div>
	</div>

	<div class="flex items-center justify-end gap-3">
		<a
			href="/content/{collectionSlug}"
			class="inline-flex h-9 items-center px-4 text-sm rounded-lg transition hover:bg-[--surface-alt]"
			style="color: var(--text-secondary)"
		>Cancel</a>
		<Button type="submit" variant="primary" {loading}>
			{loading ? 'Saving…' : '💾 Save Record'}
		</Button>
	</div>
</form>
