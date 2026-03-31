<script lang="ts">
	import { page } from '$app/state';
	import { collectionsApi } from '$lib/api/collections';
	import { contentApi } from '$lib/api/content';
	import { toast } from '$lib/stores/toast';
	import RecordTable from '$lib/features/content/RecordTable.svelte';
	import Button from '$lib/shared/Button.svelte';
	import Spinner from '$lib/shared/Spinner.svelte';
	import type { Collection, ContentRecord } from '$lib/types/collection';

	let slug = $derived(page.params.slug as string);
	let collection = $state<Collection | null>(null);
	let records = $state<ContentRecord[]>([]);
	let loading = $state(true);

	$effect(() => {
		loading = true;
		Promise.all([
			collectionsApi.get(slug),
			contentApi.list(slug)
		]).then(([coll, recs]) => {
			collection = coll;
			records = recs ?? [];
			loading = false;
		}).catch((err) => {
			toast.error(err.message ?? 'Failed to load data');
			loading = false;
		});
	});

	async function handleDelete(id: string) {
		try {
			await contentApi.delete(slug, id);
			records = records.filter((r) => (r._id ?? r.id) !== id);
			toast.success('Record deleted');
		} catch (err: any) {
			toast.error(err.message ?? 'Failed to delete record');
		}
	}
</script>

<svelte:head>
	<title>{collection?.name ?? slug} Content — GoHeadless CMS</title>
</svelte:head>

{#if loading}
	<div class="flex justify-center py-24"><Spinner /></div>
{:else if collection}
	<div class="flex flex-col gap-6 animate-fade-in">
		<!-- Header -->
		<div class="flex items-start justify-between gap-4">
			<div class="flex flex-col gap-1">
				<div class="flex items-center gap-3">
					<h1 class="text-2xl font-bold" style="color: var(--text-primary)">{collection.name}</h1>
					<span
						class="rounded-full px-2.5 py-0.5 text-xs font-medium"
						style="background: var(--surface-alt); color: var(--text-muted)"
					>{records.length} records</span>
				</div>
				<div class="flex items-center gap-3">
					<code class="text-sm font-mono" style="color: var(--text-muted)">/{collection.slug}</code>
					<a
						href="/collections/{slug}"
						class="text-xs transition hover:underline"
						style="color: var(--text-muted)"
					>View schema →</a>
				</div>
			</div>
			<Button variant="primary" onclick={() => (window.location.href = `/content/${slug}/new`)}>
				+ New Record
			</Button>
		</div>

		<RecordTable
			fields={collection.fields}
			{records}
			collectionSlug={slug}
			ondelete={handleDelete}
		/>
	</div>
{/if}
