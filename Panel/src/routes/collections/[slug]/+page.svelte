<script lang="ts">
	import { page } from '$app/state';
	import { collectionsApi } from '$lib/api/collections';
	import { toast } from '$lib/stores/toast';
	import { collectionsStore } from '$lib/stores/collections';
	import Badge from '$lib/shared/Badge.svelte';
	import Button from '$lib/shared/Button.svelte';
	import Spinner from '$lib/shared/Spinner.svelte';
	import ConfirmDialog from '$lib/shared/ConfirmDialog.svelte';
	import { goto } from '$app/navigation';
	import type { Collection } from '$lib/types/collection';

	let slug = $derived(page.params.slug as string);
	let collection = $state<Collection | null>(null);
	let loading = $state(true);
	let deleteOpen = $state(false);
	let deleteLoading = $state(false);

	$effect(() => {
		loading = true;
		collectionsApi
			.get(slug)
			.then((c) => {
				collection = c;
				loading = false;
			})
			.catch(() => {
				toast.error('Collection not found');
				loading = false;
			});
	});

	async function handleDelete() {
		deleteLoading = true;
		try {
			await collectionsApi.delete(slug);
			collectionsStore.removeLocal(slug);
			toast.success('Collection deleted');
			goto('/collections');
		} catch (err: any) {
			toast.error(err.message ?? 'Failed to delete');
		} finally {
			deleteLoading = false;
		}
	}
</script>

<svelte:head>
	<title>{collection?.name ?? slug} — GoHeadless CMS</title>
</svelte:head>

{#if loading}
	<div class="py-24 flex justify-center"><Spinner /></div>
{:else if collection}
	<div class="gap-6 animate-fade-in max-w-3xl flex flex-col">
		<!-- Header -->
		<div class="flex items-start justify-between">
			<div class="gap-1 flex flex-col">
				<h1 class="text-2xl font-bold" style="color: var(--text-primary)">{collection.name}</h1>
				<code class="text-sm font-mono" style="color: var(--text-muted)">/{collection.slug}</code>
				{#if collection.description}
					<p class="text-sm mt-1" style="color: var(--text-secondary)">{collection.description}</p>
				{/if}
			</div>
			<div class="gap-2 flex">
				<a
					href="/content/{slug}"
					class="h-9 gap-1.5 rounded-lg px-4 text-sm font-medium text-white inline-flex items-center transition hover:opacity-90"
					style="background: var(--brand)">Open Content ↗</a
				>
				<Button variant="danger" size="sm" onclick={() => (deleteOpen = true)}>🗑 Delete</Button>
			</div>
		</div>

		<!-- Schema fields table -->
		<div class="card overflow-hidden">
			<div class="card-header">
				<h2 class="text-base font-semibold" style="color: var(--text-primary)">
					Schema — {collection.fields.length} field{collection.fields.length !== 1 ? 's' : ''}
				</h2>
			</div>
			<div class="overflow-x-auto">
				<table class="text-sm w-full">
					<thead>
						<tr style="border-bottom: 1px solid var(--border); background: var(--surface-alt)">
							{#each ['Key', 'Label', 'Type', 'Required', 'Unique'] as col}
								<th
									class="px-4 py-3 text-xs font-semibold tracking-wide text-left uppercase"
									style="color: var(--text-muted)">{col}</th
								>
							{/each}
						</tr>
					</thead>
					<tbody>
						{#each collection.fields as field}
							<tr
								class="border-b transition last:border-0 hover:bg-[--surface-alt]"
								style="border-color: var(--border)"
							>
								<td class="px-4 py-3">
									<code class="text-xs font-mono" style="color: var(--text-primary)"
										>{field.key}</code
									>
								</td>
								<td class="px-4 py-3 text-sm" style="color: var(--text-primary)">{field.label}</td>
								<td class="px-4 py-3"><Badge type={field.type} /></td>
								<td class="px-4 py-3">
									<span class={field.required ? 'text-green-500' : 'text-gray-400'}
										>{field.required ? '✓' : '—'}</span
									>
								</td>
								<td class="px-4 py-3">
									<span class={field.unique ? 'text-green-500' : 'text-gray-400'}
										>{field.unique ? '✓' : '—'}</span
									>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
{/if}

<ConfirmDialog
	bind:open={deleteOpen}
	title="Delete collection"
	message="Delete '{collection?.name}'? This cannot be undone."
	confirmLabel="Delete"
	loading={deleteLoading}
	onconfirm={handleDelete}
/>
