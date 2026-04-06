<script lang="ts">
	import type { Collection } from '$lib/types/collection';
	import Badge from '$lib/shared/Badge.svelte';
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
				class="gap-2 rounded-xl px-5 py-2.5 text-sm font-semibold text-white hover:-translate-y-0.5 shadow-lg hover:shadow-xl inline-flex items-center shadow-[--brand]/20 transition-all hover:shadow-[--brand]/30"
				style="background: var(--brand)"
			>
				+ New Collection
			</a>
		{/snippet}
	</EmptyState>
{:else}
	<div class="gap-6 sm:grid-cols-2 xl:grid-cols-3 grid grid-cols-1">
		{#each collections as col}
			<div
				class="group hover:-translate-y-1 hover:shadow-xl relative flex flex-col overflow-hidden rounded-[1.25rem] border transition-all duration-300"
				style="background: var(--surface); border-color: var(--border); box-shadow: var(--shadow-sm);"
			>
				<!-- Subtle branded glow on hover -->
				<div
					class="inset-0 pointer-events-none absolute opacity-0 transition-opacity duration-300 group-hover:opacity-100"
					style="background: linear-gradient(135deg, var(--brand-light) 0%, transparent 60%);"
				></div>

				<div class="p-6 relative z-10 flex flex-1 flex-col">
					<!-- Top row: Icon and Actions -->
					<div class="mb-5 gap-4 flex items-start justify-between">
						<div
							class="h-12 w-12 flex flex-shrink-0 items-center justify-center rounded-[0.85rem] border transition-transform duration-300 group-hover:scale-105"
							style="background: var(--surface-alt); border-color: var(--border); color: var(--brand)"
						>
							<svg
								class="h-6 w-6"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
								stroke-width="1.75"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M2.25 12.75V12A2.25 2.25 0 0 1 4.5 9.75h15A2.25 2.25 0 0 1 21.75 12v.75m-8.69-6.44l-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z"
								/>
							</svg>
						</div>

						<div class="flex items-center gap-1 opacity-0 transition-all group-hover:opacity-100">
							<a
								href="/collections/{col.slug}/edit"
								class="h-8 w-8 rounded-lg text-blue-400 hover:bg-blue-50 hover:text-blue-500 dark:hover:bg-blue-500/10 flex items-center justify-center transition-all"
								title="Edit schema"
							>
								<svg class="h-4.5 w-4.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
									/>
								</svg>
							</a>

							<button
								type="button"
								onclick={() => confirmDelete(col.slug)}
								class="h-8 w-8 rounded-lg text-red-400 hover:bg-red-50 hover:text-red-500 dark:hover:bg-red-500/10 flex items-center justify-center transition-all"
								title="Delete collection"
							>
								<svg class="h-4.5 w-4.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M19 7l-.867 12.142A2 2 0 0 1 16.138 21H7.862a2 2 0 0 1-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v3M4 7h16"
									/>
								</svg>
							</button>
						</div>
					</div>

					<!-- Content -->
					<div class="mb-2 gap-1.5 flex flex-col">
						<h3
							class="font-bold tracking-tight truncate text-[1.15rem]"
							style="color: var(--text-primary)"
						>
							{col.name}
						</h3>
						<div class="gap-2.5 flex items-center">
							<code
								class="rounded-md px-2 py-0.5 text-xs font-mono font-medium truncate"
								style="background: var(--surface-alt); color: var(--text-muted)">/{col.slug}</code
							>
							<span class="text-xs font-medium" style="color: var(--text-muted)">
								• {col.fields.length} schema block{col.fields.length !== 1 ? 's' : ''}
							</span>
						</div>
					</div>

					{#if col.description}
						<p class="mt-1 text-sm line-clamp-2" style="color: var(--text-secondary)">
							{col.description}
						</p>
					{/if}

					<!-- Fields preview -->
					<div class="pt-6 gap-1.5 mt-auto flex flex-wrap">
						{#each col.fields.slice(0, 4) as field}
							<Badge type={field.type} />
						{/each}
						{#if col.fields.length > 4}
							<span
								class="rounded-lg px-2.5 py-0.5 font-semibold inline-flex items-center text-[0.7rem]"
								style="background: var(--surface-alt); color: var(--text-muted)"
								>+{col.fields.length - 4}</span
							>
						{/if}
						{#if col.fields.length === 0}
							<span class="text-xs italic" style="color: var(--text-muted)"
								>Unconfigured schema</span
							>
						{/if}
					</div>
				</div>

				<!-- Floating Action Footer -->
				<div
					class="p-2 relative z-10 border-t transition-colors duration-300"
					style="border-color: var(--border); background: transparent"
				>
					<a
						href="/content/{col.slug}"
						class="gap-2 rounded-xl py-3 text-sm font-semibold flex w-full items-center justify-center transition-all hover:bg-[--brand-light]"
						style="color: var(--brand)"
					>
						Manage Content
						<svg
							class="h-4 w-4 group-hover:translate-x-1 transition-transform duration-300"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2.5"
								d="M13.5 4.5 21 12m0 0-7.5 7.5M21 12H3"
							/>
						</svg>
					</a>
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
