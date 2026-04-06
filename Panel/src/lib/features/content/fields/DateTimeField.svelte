<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import BaseField from './BaseField.svelte';

	interface Props {
		field: Field;
		data: any;
	}

	let { field, data = $bindable() }: Props = $props();

	let inputType = $derived(
		field.type === 'datepicker' ? 'date' :
		field.type === 'timepicker' ? 'time' :
		field.type === 'datetimepicker' ? 'datetime-local' : 'date'
	);

	let _start = $state('');
	let _end = $state('');

	$effect(() => {
		if (field.type === 'daterange') {
			if (data && typeof data === 'object') {
				_start = data.start || '';
				_end = data.end || '';
			} else if (typeof data === 'string') {
                // optional fallback to avoid errors if incorrectly stored
                _start = '';
                _end = '';
            }
		}
	});

	function updateDateRange() {
		data = { start: _start, end: _end };
	}
</script>

<BaseField {field}>
	{#snippet children({ id })}
		{#if field.type === 'daterange'}
			<div class="daterange-input-container">
				<input
					{id}
					type="date"
					bind:value={_start}
					oninput={updateDateRange}
					class="datetime-input flex-1"
				/>
				<span class="daterange-separator">to</span>
				<input
					type="date"
					bind:value={_end}
					oninput={updateDateRange}
					class="datetime-input flex-1"
				/>
			</div>
		{:else}
			<div class="datetime-input-container">
				<input
					{id}
					type={inputType}
					bind:value={data}
					class="datetime-input"
				/>
			</div>
		{/if}
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

	.daterange-input-container {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		width: 100%;
	}

	.daterange-separator {
		font-size: 0.875rem;
		color: var(--text-muted);
		font-weight: 500;
	}

	/* Specific style for calendar icons in some browsers */
	.datetime-input::-webkit-calendar-picker-indicator {
		filter: var(--calendar-icon-filter, none);
		cursor: pointer;
	}
</style>
