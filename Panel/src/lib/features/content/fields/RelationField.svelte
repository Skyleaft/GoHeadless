<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import BaseField from './BaseField.svelte';

	interface Props {
		field: Field;
		data: any;
	}

	let { field, data = $bindable() }: Props = $props();

	let options = $state<{ label: string; value: string }[]>([]);
	let loading = $state(false);
	let searched = $state(false);

	async function loadRelation() {
		if (searched || !field.relation) return;
		loading = true;
		try {
			const res = await fetch(`/api/v1/content/${field.relation.collection}`);
			const records: Record<string, any>[] = await res.json();
			const displayField = field.relation.field;
			options = records.map((r) => ({
				label: String(r[displayField] ?? r._id ?? ''),
				value: String(r._id ?? r.id ?? '')
			}));
		} catch {
			// ignore
		} finally {
			loading = false;
			searched = true;
		}
	}
</script>

<BaseField {field}>
	{#snippet children({ id })}
		<select {id} bind:value={data} onfocus={loadRelation} class="relation-select">
			<option value="">
				{loading ? 'Loading…' : 'Select related record…'}
			</option>
			{#each options as opt}
				<option value={opt.value}>{opt.label}</option>
			{/each}
		</select>
	{/snippet}
</BaseField>

<style>
	.relation-select {
		height: 2.25rem;
		border-radius: 0.5rem;
		padding: 0 0.75rem;
		font-size: 0.875rem;
		width: 100%;
		border: 1px solid var(--border);
		background: var(--surface);
		color: var(--text-primary);
		transition: border-color 0.2s;
		cursor: pointer;
	}

	.relation-select:focus {
		outline: none;
		border-color: var(--brand);
	}
</style>
