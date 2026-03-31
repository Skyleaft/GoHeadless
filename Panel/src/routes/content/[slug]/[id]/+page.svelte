<script lang="ts">
	import { goto } from '$app/navigation';
	import { contentApi } from '$lib/api/content';
	import { toast } from '$lib/stores/toast';
	import RecordForm from '$lib/features/content/RecordForm.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let slug = $derived(data.slug);
	let id = $derived(data.id);
	let collection = $derived(data.collection);
	let record = $derived(data.record);
	let saving = $state(false);

	async function handleSubmit(formData: any) {
		saving = true;
		try {
			await contentApi.update(slug, id, formData);
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

{#if collection && record}
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
