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

<div class="max-w-4xl animate-fade-in gap-10 py-4 mx-auto flex w-full flex-col">
	<div class="px-2 relative">
		<div
			class="-left-4 w-1 h-12 shadow-lg absolute top-1/2 -translate-y-1/2 rounded-full bg-[--brand] shadow-[--brand]/40"
		></div>
		<h1 class="text-4xl font-black tracking-tighter" style="color: var(--text-primary)">
			New Collection
		</h1>
		<p class="text-sm font-bold mt-2 tracking-tight opacity-50" style="color: var(--text-muted)">
			Architect your data structure with dynamic fields and access policies.
		</p>
	</div>

	<CollectionForm {loading} onsubmit={handleSubmit} />
</div>
