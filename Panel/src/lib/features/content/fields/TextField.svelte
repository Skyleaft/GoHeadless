<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import BaseField from './BaseField.svelte';

	interface Props {
		field: Field;
		data: any;
	}

	let { field, data = $bindable() }: Props = $props();

	let inputType = $derived(
		field.type === 'email'
			? 'email'
			: field.type === 'password'
				? 'password'
				: field.type === 'url'
					? 'url'
					: field.type === 'phone'
						? 'tel'
						: 'text'
	);
</script>

<BaseField {field}>
	{#snippet children({ id })}
		<input
			{id}
			type={inputType}
			bind:value={data}
			placeholder={field.placeholder}
			required={field.required}
			class="text-input"
		/>
	{/snippet}
</BaseField>

<style>
	.text-input {
		height: 2.25rem; /* h-9 */
		border-radius: 0.5rem; /* rounded-lg */
		padding: 0 0.75rem; /* px-3 */
		font-size: 0.875rem; /* text-sm */
		width: 100%;
		border: 1px solid var(--border);
		background: var(--surface);
		color: var(--text-primary);
		transition:
			border-color 0.2s,
			box-shadow 0.2s;
	}

	.text-input:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px rgba(var(--brand-rgb), 0.1);
	}
</style>
