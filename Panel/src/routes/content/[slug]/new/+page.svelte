<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { collectionsApi } from '$lib/api/collections';
	import { contentApi } from '$lib/api/content';
	import { toast } from '$lib/stores/toast';
	import RecordForm from '$lib/features/content/RecordForm.svelte';
	import Spinner from '$lib/shared/Spinner.svelte';
	import type { Collection } from '$lib/types/collection';

	let slug = $derived(page.params.slug as string);
	let collection = $state<Collection | null>(null);
	let loading = $state(true);
	let saving = $state(false);

	$effect(() => {
		loading = true;
		collectionsApi.get(slug).then((c) => {
			collection = c;
			loading = false;
		}).catch((err) => {
			toast.error(err.message ?? 'Failed to load collection');
			loading = false;
		});
	});

	async function handleSubmit(data: Record<string, unknown>) {
		saving = true;
		try {
			const result = await contentApi.create(slug, data);
			toast.success('Record created successfully!');
			goto(`/content/${slug}`);
		} catch (err: any) {
			toast.error(err.message ?? 'Failed to create record');
		} finally {
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>New Record — {collection?.name ?? slug} — GoHeadless CMS</title>
</svelte:head>

{#if loading}
	<div class="flex justify-center py-24"><Spinner /></div>
{:else if collection}
	<div class="flex flex-col gap-6 animate-fade-in max-w-2xl">
		<div>
			<div class="flex items-center gap-2 text-sm mb-2">
				<a href="/content/{slug}" class="transition hover:text-[--brand]" style="color: var(--text-muted)">
					← {collection.name}
				</a>
			</div>
			<h1 class="text-2xl font-bold" style="color: var(--text-primary)">New Record</h1>
			<p class="text-sm mt-1" style="color: var(--text-muted)">
				Create a new entry in the <strong>{collection.name}</strong> collection.
			</p>
		</div>

		<RecordForm
			fields={collection.fields}
			collectionSlug={slug}
			loading={saving}
			onsubmit={handleSubmit}
		/>
	</div>
{/if}
