<script lang="ts">
	import type { Collection } from '$lib/types/collection';
	import Badge from '$lib/shared/Badge.svelte';
	import Button from '$lib/shared/Button.svelte';
	import ConfirmDialog from '$lib/shared/ConfirmDialog.svelte';
	import EmptyState from '$lib/shared/EmptyState.svelte';

	interface Props {
		collections: Collection[];
		ondelete?: (slug: string) => void;
		loading?: boolean;
	}

	let { collections, ondelete, loading = false }: Props = $props();

	let deleteTarget = $state<string | null>(null);
	let deleteLoading = $state(false);

	function confirmDelete(slug: string) {
		deleteTarget = slug;
	}

	async function handleDelete() {
		if (!deleteTarget) return;
		deleteLoading = true;
		await ondelete?.(deleteTarget);
		deleteLoading = false;
		deleteTarget = null;
	}
</script>

{#if collections.length === 0 && !loading}
	<EmptyState
		icon="🗂"
		title="No collections yet"
		description="Create your first collection to start managing dynamic content."
	>
		{#snippet children()}
			<a
				href="/collections/new"
				class="gap-2 rounded-lg px-4 py-2 text-sm font-medium text-white inline-flex items-center transition-all hover:opacity-90"
				style="background: var(--brand)"
			>
				+ New Collection
			</a>
		{/snippet}
	</EmptyState>
{:else}
	<div class="gap-4 sm:grid-cols-2 xl:grid-cols-3 grid grid-cols-1">
		{#each collections as col}
			<div
				class="card group gap-4 p-5 relative flex flex-col transition-all hover:border-[--brand]"
			>
				<!-- Header -->
				<div class="gap-2 flex items-start justify-between">
					<div class="gap-1 flex flex-col overflow-hidden">
						<h3 class="text-base font-semibold truncate" style="color: var(--text-primary)">
							{col.name}
						</h3>
						<code
							class="rounded px-1.5 py-0.5 text-xs font-mono truncate"
							style="background: var(--surface-alt); color: var(--text-muted)">/{col.slug}</code
						>
					</div>
					<div
						class="h-10 w-10 rounded-xl text-lg flex flex-shrink-0 items-center justify-center"
						style="background: var(--brand-light); color: var(--brand)"
					>
						🗂
					</div>
				</div>

				<!-- Description -->
				{#if col.description}
					<p class="text-sm line-clamp-2" style="color: var(--text-secondary)">{col.description}</p>
				{/if}

				<!-- Fields preview -->
				<div class="gap-1.5 flex flex-wrap">
					{#each col.fields.slice(0, 5) as field}
						<Badge type={field.type} />
					{/each}
					{#if col.fields.length > 5}
						<span
							class="rounded-md px-2 py-0.5 text-xs inline-flex items-center"
							style="background: var(--surface-alt); color: var(--text-muted)"
							>+{col.fields.length - 5} more</span
						>
					{/if}
					{#if col.fields.length === 0}
						<span class="text-xs" style="color: var(--text-muted)">No fields defined</span>
					{/if}
				</div>

				<!-- Footer: Stats + Actions -->
				<div
					class="pt-4 flex items-center justify-between border-t"
					style="border-color: var(--border)"
				>
					<span class="text-xs" style="color: var(--text-muted)">
						{col.fields.length} field{col.fields.length !== 1 ? 's' : ''}
					</span>
					<div class="gap-2 flex items-center">
						<Button variant="ghost" size="sm" onclick={() => confirmDelete(col.slug)}>
							🗑 Delete
						</Button>
						<a
							href="/content/{col.slug}"
							class="gap-1 rounded-lg px-3 py-1.5 text-xs font-medium text-white inline-flex items-center transition-all hover:opacity-90"
							style="background: var(--brand)"
						>
							Open ↗
						</a>
					</div>
				</div>
			</div>
		{/each}
	</div>
{/if}

<ConfirmDialog
	open={!!deleteTarget}
	title="Delete collection"
	message="This will permanently delete the '{deleteTarget}' collection definition. Records in this collection will NOT be removed from the database."
	confirmLabel="Delete"
	loading={deleteLoading}
	onconfirm={handleDelete}
	oncancel={() => (deleteTarget = null)}
/>
