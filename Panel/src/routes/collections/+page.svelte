<script lang="ts">
	import { collectionsStore } from '$lib/stores/collections';
	import { collectionsApi } from '$lib/api/collections';
	import { toast } from '$lib/stores/toast';
	import CollectionList from '$lib/features/collections/CollectionList.svelte';
	import Button from '$lib/shared/Button.svelte';
	import Spinner from '$lib/shared/Spinner.svelte';
	import type { Collection } from '$lib/types/collection';

	let collections = $derived($collectionsStore);
	let loading = $state(true);

	$effect(() => {
		collectionsStore.load(true).then(() => (loading = false));
	});

	async function handleDelete(slug: string) {
		try {
			await collectionsApi.delete(slug);
			collectionsStore.removeLocal(slug);
			toast.success(`Collection "${slug}" deleted`);
		} catch (err: any) {
			toast.error(err.message ?? 'Failed to delete collection');
		}
	}
</script>

<svelte:head>
	<title>Collections — GoHeadless CMS</title>
</svelte:head>

<div class="flex flex-col gap-6 animate-fade-in">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold" style="color: var(--text-primary)">Collections</h1>
			<p class="text-sm mt-1" style="color: var(--text-muted)">
				{collections.length} collection{collections.length !== 1 ? 's' : ''} defined
			</p>
		</div>
		<Button variant="primary" onclick={() => (window.location.href = '/collections/new')}>
			+ New Collection
		</Button>
	</div>

	{#if loading}
		<div class="flex justify-center py-20"><Spinner /></div>
	{:else}
		<CollectionList {collections} ondelete={handleDelete} />
	{/if}
</div>
