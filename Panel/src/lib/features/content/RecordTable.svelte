<script lang="ts">
	import type { Field, ContentRecord } from '$lib/types/collection';
	import { isStructuralType } from '$lib/types/collection';
	import { toast } from '$lib/stores/toast';
	import Badge from '$lib/shared/Badge.svelte';
	import Button from '$lib/shared/Button.svelte';
	import EmptyState from '$lib/shared/EmptyState.svelte';
	import ConfirmDialog from '$lib/shared/ConfirmDialog.svelte';

	interface Props {
		fields: Field[];
		records: ContentRecord[];
		collectionSlug: string;
		ondelete?: (id: string) => void;
		loading?: boolean;
		canUpdate?: boolean;
		canDelete?: boolean;
	}

	let {
		fields,
		records,
		collectionSlug,
		ondelete,
		loading = false,
		canUpdate = true,
		canDelete = true
	}: Props = $props();

	let deleteTarget = $state<string | null>(null);
	let deleteLoading = $state(false);

	// Flatten leaf fields (skip pure layout containers)
	let displayFields = $derived(
		fields.filter((f) => !['section', 'tabs', 'grid'].includes(f.type)).slice(0, 8)
	);

	function formatValue(val: unknown): string {
		if (val === null || val === undefined) return '—';
		if (typeof val === 'boolean') return val ? '✓ Yes' : '✗ No';
		if (typeof val === 'object') return JSON.stringify(val).slice(0, 60) + '…';
		const s = String(val);
		return s.length > 80 ? s.slice(0, 80) + '…' : s;
	}

	function getRecordId(record: ContentRecord): string {
		return (record._id as string) ?? (record.id as string) ?? '';
	}

	async function handleDelete() {
		if (!deleteTarget) return;
		deleteLoading = true;
		await ondelete?.(deleteTarget);
		deleteLoading = false;
		deleteTarget = null;
	}
</script>

{#if records.length === 0 && !loading}
	<EmptyState
		icon="📄"
		title="No records yet"
		description="Create the first record for this collection."
	>
		{#snippet children()}
			<a
				href="/content/{collectionSlug}/new"
				class="inline-flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium text-white transition-all hover:opacity-90"
				style="background: var(--brand)"
			>+ New Record</a>
		{/snippet}
	</EmptyState>
{:else}
	<div class="card overflow-hidden">
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead>
					<tr style="border-bottom: 1px solid var(--border); background: var(--surface-alt)">
						<th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide" style="color: var(--text-muted)">
							ID
						</th>
						{#each displayFields as field}
							<th class="px-4 py-3 text-left text-xs font-semibold uppercase tracking-wide whitespace-nowrap" style="color: var(--text-muted)">
								<div class="flex items-center gap-2">
									{field.label || field.key}
									<Badge type={field.type} />
								</div>
							</th>
						{/each}
						<th class="px-4 py-3 text-right text-xs font-semibold uppercase tracking-wide" style="color: var(--text-muted)">
							Actions
						</th>
					</tr>
				</thead>
				<tbody>
					{#each records as record}
						{@const id = getRecordId(record)}
						<tr
							class="transition-colors border-b last:border-0"
							style="border-color: var(--border)"
							onmouseenter={(e) => (e.currentTarget as HTMLElement).style.background = 'var(--surface-alt)'}
							onmouseleave={(e) => (e.currentTarget as HTMLElement).style.background = 'transparent'}
						>
							<!-- ID -->
							<td class="px-4 py-3">
								<code
									class="rounded px-1.5 py-0.5 text-xs font-mono"
									style="background: var(--surface-alt); color: var(--text-muted)"
								>{id.slice(-8)}</code>
							</td>

							<!-- Fields -->
							{#each displayFields as field}
								<td class="px-4 py-3 max-w-xs" style="color: var(--text-primary)">
									{#if field.type === 'image' && record[field.key]}
										<img
											src="/uploads{record[field.key]}"
											alt={field.label}
											class="h-8 w-8 rounded object-cover"
										/>
									{:else if field.type === 'bool' || field.type === 'toggle'}
										<span class={record[field.key] ? 'text-green-500' : 'text-red-400'}>
											{record[field.key] ? '✓ Yes' : '✗ No'}
										</span>
									{:else}
										<span class="block truncate">{formatValue(record[field.key])}</span>
									{/if}
								</td>
							{/each}

							<!-- Actions -->
							<td class="px-4 py-3 text-right">
								<div class="flex items-center justify-end gap-2">
									{#if canUpdate}
										<a
											href="/content/{collectionSlug}/{id}"
											class="inline-flex h-7 items-center gap-1 rounded-lg px-2.5 text-xs font-medium transition hover:bg-[--surface-hover]"
											style="color: var(--brand)"
										>Edit</a>
									{/if}
									{#if canDelete}
										<button
											onclick={() => (deleteTarget = id)}
											class="flex h-7 w-7 items-center justify-center rounded-lg text-red-400 transition hover:bg-red-50 hover:text-red-600"
											title="Delete record"
										>🗑</button>
									{/if}
									{#if !canUpdate && !canDelete}
										<span class="text-xs text-[--text-muted]">Read only</span>
									{/if}
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Footer -->
		<div
			class="flex items-center justify-between px-4 py-3 border-t text-xs"
			style="border-color: var(--border); color: var(--text-muted); background: var(--surface-alt)"
		>
			<span>{records.length} record{records.length !== 1 ? 's' : ''} total</span>
			{#if fields.length > 8}
				<span>{fields.length - displayFields.length} columns hidden</span>
			{/if}
		</div>
	</div>
{/if}

<ConfirmDialog
	open={!!deleteTarget}
	title="Delete record"
	message="This will permanently delete this record. This action cannot be undone."
	confirmLabel="Delete"
	loading={deleteLoading}
	onconfirm={handleDelete}
	oncancel={() => (deleteTarget = null)}
/>
