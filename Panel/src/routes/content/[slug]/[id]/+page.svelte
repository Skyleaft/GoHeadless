<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { collectionsApi } from '$lib/api/collections';
	import { contentApi } from '$lib/api/content';
	import { toast } from '$lib/stores/toast';
	import RecordForm from '$lib/features/content/RecordForm.svelte';
	import Button from '$lib/shared/Button.svelte';
	import Spinner from '$lib/shared/Spinner.svelte';
	import type { Collection, ContentRecord } from '$lib/types/collection';

	let slug = $derived(page.params.slug as string);
	let id = $derived(page.params.id as string);

	let collection = $state<Collection | null>(null);
	let record = $state<ContentRecord | null>(null);
	let loading = $state(true);
	let saving = $state(false);

	$effect(() => {
		loading = true;
		Promise.all([
			collectionsApi.get(slug),
			contentApi.get(slug, id)
		]).then(([coll, rec]) => {
			collection = coll;
			record = rec;
			loading = false;
		}).catch((err) => {
			toast.error(err.message ?? 'Record not found');
			loading = false;
		});
	});

	async function handleSubmit(data: ContentRecord) {
		saving = true;
		try {
			await contentApi.update(slug, id, data);
			toast.success('Record updated!');
			goto(`/content/${slug}`);
		} catch (err: any) {
			toast.error(err.message ?? 'Failed to update record');
		} finally {
			saving = false;
		}
	}
</script>

<svelte:head>
	<title>Edit Record — {collection?.name ?? slug} — GoHeadless CMS</title>
</svelte:head>

{#if loading}
	<div class="flex justify-center py-24"><Spinner /></div>
{:else if collection && record}
	<div class="flex flex-col gap-6 animate-fade-in max-w-2xl">
		<div>
			<div class="flex items-center gap-2 text-sm mb-2">
				<a href="/content/{slug}" class="transition hover:text-[--brand]" style="color: var(--text-muted)">
					← {collection.name}
				</a>
			</div>
			<div class="flex items-start justify-between">
				<div>
					<h1 class="text-2xl font-bold" style="color: var(--text-primary)">Edit Record</h1>
					<code class="text-xs font-mono mt-1" style="color: var(--text-muted)">{id}</code>
				</div>
			</div>
		</div>

		<RecordForm
			fields={collection.fields}
			initial={record}
			collectionSlug={slug}
			loading={saving}
			onsubmit={handleSubmit}
		/>
	</div>
{/if}
