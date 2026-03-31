<script lang="ts">
	import { goto } from '$app/navigation';
	import { contentApi } from '$lib/api/content';
	import { toast } from '$lib/stores/toast';
	import RecordForm from '$lib/features/content/RecordForm.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let slug = $derived(data.slug);
	let collection = $derived(data.collection);
	let saving = $state(false);

	async function handleSubmit(formData: Record<string, unknown>) {
		saving = true;
		try {
			await contentApi.create(slug, formData);
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

{#if collection}
	<div class="gap-6 animate-fade-in max-w-2xl flex flex-col">
		<div>
			<div class="gap-2 text-sm mb-2 flex items-center">
				<a
					href="/content/{slug}"
					class="transition hover:text-[--brand]"
					style="color: var(--text-muted)"
				>
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
