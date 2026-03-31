<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import BaseField from './BaseField.svelte';

	interface Props {
		field: Field;
		data: string;
	}

	let { field, data = $bindable() }: Props = $props();
</script>

<BaseField {field}>
	{#snippet children({ id })}
		<textarea
			{id}
			bind:value={data}
			placeholder={field.placeholder}
			rows="4"
			class="textarea-input"
		></textarea>
	{/snippet}
</BaseField>

<style>
	.textarea-input {
		padding: 0.5rem 0.75rem; /* px-3 py-2 */
		font-size: 0.875rem; /* text-sm */
		width: 100%;
		min-height: 5rem;
		resize: vertical;
		border-radius: 0.5rem; /* rounded-lg */
		border: 1px solid var(--border);
		background: var(--surface);
		color: var(--text-primary);
		transition: border-color 0.2s, box-shadow 0.2s;
	}

	.textarea-input:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px rgba(var(--brand-rgb), 0.1);
	}
</style>
