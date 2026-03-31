<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import BaseField from './BaseField.svelte';

	interface Props {
		field: Field;
		data: string;
	}

	let { field, data = $bindable() }: Props = $props();

	let inputType = $derived(
		field.type === 'datepicker' ? 'date' :
		field.type === 'timepicker' ? 'time' :
		field.type === 'datetimepicker' ? 'datetime-local' : 'date'
	);
</script>

<BaseField {field}>
	{#snippet children({ id })}
		<div class="datetime-input-container">
			<input
				{id}
				type={inputType}
				bind:value={data}
				class="datetime-input"
			/>
		</div>
	{/snippet}
</BaseField>

<style>
	.datetime-input {
		height: 2.25rem;
		border-radius: 0.5rem;
		padding: 0 0.75rem;
		font-size: 0.875rem;
		width: 100%;
		border: 1px solid var(--border);
		background: var(--surface);
		color: var(--text-primary);
		transition: border-color 0.2s;
	}

	.datetime-input:focus {
		outline: none;
		border-color: var(--brand);
	}

	/* Specific style for calendar icons in some browsers */
	.datetime-input::-webkit-calendar-picker-indicator {
		filter: var(--calendar-icon-filter, none);
		cursor: pointer;
	}
</style>
