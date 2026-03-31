<script lang="ts">
	import { goto } from '$app/navigation';
	import { collectionsApi } from '$lib/api/collections';
	import { collectionsStore } from '$lib/stores/collections';
	import { toast } from '$lib/stores/toast';
	import CollectionForm from '$lib/features/collections/CollectionForm.svelte';
	import type { Collection } from '$lib/types/collection';

	let loading = $state(false);

	async function handleSubmit(collection: Omit<Collection, 'id'>) {
		loading = true;
		try {
			const created = await collectionsApi.create(collection);
			collectionsStore.addLocal(created);
			toast.success(`Collection "${created.name}" created!`);
			goto('/collections');
		} catch (err: any) {
			toast.error(err.message ?? 'Failed to create collection');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>New Collection — GoHeadless CMS</title>
</svelte:head>

<div class="flex flex-col gap-6 animate-fade-in max-w-4xl">
	<div>
		<h1 class="text-2xl font-bold" style="color: var(--text-primary)">New Collection</h1>
		<p class="text-sm mt-1" style="color: var(--text-muted)">
			Define a new dynamic collection with a custom schema.
		</p>
	</div>

	<CollectionForm {loading} onsubmit={handleSubmit} />
</div>
