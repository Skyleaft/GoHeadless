<script lang="ts">
	import type { PageData } from './$types';
	import { goto } from '$app/navigation';
	import { collectionsStore } from '$lib/stores/collections';
	import { collectionsApi } from '$lib/api/collections';
	import CollectionForm from '$lib/features/collections/CollectionForm.svelte';
	import Spinner from '$lib/shared/Spinner.svelte';
	import type { Collection } from '$lib/types/collection';
	import { toast } from '$lib/stores/toast';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	let collection = $derived($collectionsStore.find((c) => c.slug === data.slug));
	let loading = $state(false);

	$effect(() => {
		if (!collection) {
			collectionsStore.load();
		}
	});

	async function handleSubmit(payload: Omit<Collection, 'id'>) {
		loading = true;
		try {
			const updated = await collectionsApi.update(data.slug, payload);
			collectionsStore.updateLocal(data.slug, updated);
			toast.success('Collection updated successfully');
			goto('/collections');
		} catch (err: any) {
			console.error(err);
			toast.error(err.message || 'Failed to update collection');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Edit Collection - GoHeadless</title>
</svelte:head>

<div class="max-w-5xl mx-auto">
	<div class="mb-8 gap-4 flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-black tracking-tight" style="color: var(--text-primary)">
				Edit Schema
			</h1>
			<p class="mt-2 text-sm font-medium opacity-60" style="color: var(--text-secondary)">
				Update fields and access control for "{@html collection?.name || data.slug}"
			</p>
		</div>
		<button
			onclick={() => goto('/collections')}
			class="gap-2 rounded-xl px-4 py-2 font-bold shadow-sm hover:bg-gray-50/50 flex items-center border transition-all"
			style="background: var(--surface); border-color: var(--border); color: var(--text-secondary)"
		>
			Cancel
		</button>
	</div>

	{#if !collection}
		<div class="py-20 flex items-center justify-center">
			<Spinner />
		</div>
	{:else}
		<CollectionForm initial={collection} isEditMode={true} {loading} onsubmit={handleSubmit} />
	{/if}
</div>
