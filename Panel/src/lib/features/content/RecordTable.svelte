<script lang="ts">
	import type { Field, ContentRecord } from '$lib/types/collection';
	import { isStructuralType } from '$lib/types/collection';
	import { toast } from '$lib/stores/toast';
	import Badge from '$lib/shared/Badge.svelte';
	import Button from '$lib/shared/Button.svelte';
	import EmptyState from '$lib/shared/EmptyState.svelte';
	import ConfirmDialog from '$lib/shared/ConfirmDialog.svelte';
	import { getFileUrl } from '$lib/utils/path';

	interface Props {
		fields: Field[];
		records: ContentRecord[];
		collectionSlug: string;
		ondelete?: (id: string) => void;
		loading?: boolean;
		canUpdate?: boolean;
		canDelete?: boolean;
		/** Current `sort` query value, e.g. `title` or `-created_at` */
		sortParam?: string | null;
		onSortClick?: (fieldKey: string) => void;
		/** Total matching rows (when paginated); falls back to records.length */
		totalCount?: number;
	}

	let {
		fields,
		records,
		collectionSlug,
		ondelete,
		loading = false,
		canUpdate = true,
		canDelete = true,
		sortParam = null,
		onSortClick,
		totalCount
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
				class="gap-2 rounded-lg px-4 py-2 text-sm font-medium text-white inline-flex items-center transition-all hover:opacity-90"
				style="background: var(--brand)">+ New Record</a
			>
		{/snippet}
	</EmptyState>
{:else}
	<div class="card overflow-hidden">
		<div class="overflow-x-auto">
			<table class="text-sm w-full">
				<thead>
					<tr style="border-bottom: 1px solid var(--border); background: var(--surface-alt)">
						<th
							class="px-4 py-3 text-xs font-semibold tracking-wide text-left uppercase"
							style="color: var(--text-muted)"
						>
							ID
						</th>
						{#each displayFields as field}
							<th
								class="px-4 py-3 text-xs font-semibold tracking-wide text-left whitespace-nowrap uppercase"
								style="color: var(--text-muted)"
							>
								{#if onSortClick}
									<button
										type="button"
										class="gap-2 max-w-full flex items-center text-left transition hover:opacity-90"
										style="color: var(--text-muted)"
										onclick={() => onSortClick?.(field.key)}
									>
										<span class="truncate">{field.label || field.key}</span>
										<Badge type={field.type} />
										<span class="shrink-0 text-[10px] leading-none" aria-hidden="true">
											{#if sortParam === field.key}
												<span style="color: var(--brand)">▲</span>
											{:else if sortParam === `-${field.key}`}
												<span style="color: var(--brand)">▼</span>
											{:else}
												<span class="opacity-35">⇅</span>
											{/if}
										</span>
									</button>
								{:else}
									<div class="gap-2 flex items-center">
										{field.label || field.key}
										<Badge type={field.type} />
									</div>
								{/if}
							</th>
						{/each}
						<th
							class="px-4 py-3 text-xs font-semibold tracking-wide text-right uppercase"
							style="color: var(--text-muted)"
						>
							Actions
						</th>
					</tr>
				</thead>
				<tbody>
					{#each records as record}
						{@const id = getRecordId(record)}
						<tr
							class="border-b transition-colors last:border-0"
							style="border-color: var(--border)"
							onmouseenter={(e) =>
								((e.currentTarget as HTMLElement).style.background = 'var(--surface-alt)')}
							onmouseleave={(e) =>
								((e.currentTarget as HTMLElement).style.background = 'transparent')}
						>
							<!-- ID -->
							<td class="px-4 py-3">
								<code
									class="rounded px-1.5 py-0.5 text-xs font-mono"
									style="background: var(--surface-alt); color: var(--text-muted)"
									>{id.slice(-8)}</code
								>
							</td>

							<!-- Fields -->
							{#each displayFields as field}
								<td class="px-4 py-3 max-w-xs" style="color: var(--text-primary)">
									{#if field.type === 'image' && record[field.key]}
										<img
											src={getFileUrl(record[field.key] as string)}
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
								<div class="gap-2 flex items-center justify-end">
									{#if canUpdate}
										<a
											href="/content/{collectionSlug}/{id}"
											class="h-7 gap-1 rounded-lg px-2.5 text-xs font-medium inline-flex items-center transition hover:bg-[--surface-hover]"
											style="color: var(--brand)">Edit</a
										>
									{/if}
									{#if canDelete}
										<button
											onclick={() => (deleteTarget = id)}
											class="h-7 w-7 rounded-lg text-red-400 hover:bg-red-50 hover:text-red-600 flex items-center justify-center transition"
											title="Delete record">🗑</button
										>
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
			class="px-4 py-3 text-xs flex items-center justify-between border-t"
			style="border-color: var(--border); color: var(--text-muted); background: var(--surface-alt)"
		>
			<span
				>{totalCount ?? records.length} record{(totalCount ?? records.length) !== 1 ? 's' : ''} total</span
			>
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
