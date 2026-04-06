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
		initial ??
			Object.fromEntries(
				fields.map((f) => [
					f.key,
					f.default_value ?? (f.type === 'bool' || f.type === 'toggle' ? false : '')
				])
			)
	);

	import { goto } from '$app/navigation';
	import { uploadApi } from '$lib/api/upload';

	function handleSubmit(e: Event) {
		e.preventDefault();
		onsubmit?.(data);
	}

	function extractFilePaths(fieldsList: Field[], dataObj: any): string[] {
		if (!dataObj) return [];
		let paths: string[] = [];

		for (const field of fieldsList) {
			const val = dataObj[field.key];
			if (!val) continue;

			if (field.type === 'file' || field.type === 'image') {
				if (field.is_array && Array.isArray(val)) {
					paths.push(...val);
				} else if (typeof val === 'string') {
					paths.push(val);
				}
			} else if (field.fields && (field.type === 'group' || field.type === 'repeater')) {
				if (Array.isArray(val)) {
					for (const item of val) {
						paths.push(...extractFilePaths(field.fields, item));
					}
				} else {
					paths.push(...extractFilePaths(field.fields, val));
				}
			}
		}
		const filtered = paths
			.filter((p) => typeof p === 'string' && (p.startsWith('/uploads') || p.startsWith('\\uploads')))
			.map((p) => p.replace(/\\/g, '/'));
		console.log('extractFilePaths -> Extracted:', paths, 'Filtered:', filtered);
		return filtered;
	}

	let isCanceling = $state(false);

	async function handleCancel() {
		console.log('handleCancel Triggered. initial:', initial, 'data:', data);
		// Only remove uploaded files if this is a "Create Record" form (no initial data)
		if (!initial && data) {
			const pathsToDelete = extractFilePaths(fields, data);
			console.log('handleCancel -> pathsToDelete:', pathsToDelete);

			if (pathsToDelete.length > 0) {
				// Prevent users from clicking multiple times
				isCanceling = true;
				// Delete all unused files silently
				await Promise.all(
					pathsToDelete.map((path) =>
						uploadApi.delete(path).catch((err) => console.log('Delete fail:', err))
					)
				);
				isCanceling = false;
			}
		}

		goto(`/content/${collectionSlug}`);
	}
</script>

<form onsubmit={handleSubmit} class="gap-6 flex flex-col">
	<div class="card">
		<div class="card-body gap-5 flex flex-col">
			{#each fields as field}
				<FieldRenderer {field} bind:data />
			{/each}

			{#if fields.length === 0}
				<p class="py-8 text-sm text-center" style="color: var(--text-muted)">
					This collection has no fields defined yet.
					<a href="/collections" class="text-[--brand] hover:underline">Manage collections ↗</a>
				</p>
			{/if}
		</div>
	</div>

	<div class="gap-3 flex items-center justify-end">
		<button
			type="button"
			onclick={handleCancel}
			class="h-9 px-4 text-sm rounded-lg inline-flex items-center transition hover:bg-[--surface-alt]"
			style="color: var(--text-secondary)"
			disabled={loading || isCanceling}
		>
			Cancel
		</button>
		<Button type="submit" variant="primary" {loading}>
			{loading ? 'Saving…' : '💾 Save Record'}
		</Button>
	</div>
</form>
