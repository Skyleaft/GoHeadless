<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import BaseField from './BaseField.svelte';

	interface Props {
		field: Field;
		data: any;
	}

	let { field, data = $bindable() }: Props = $props();

	function handleCheckChanged(e: Event) {
		data = (e.target as HTMLInputElement).checked;
	}

	function handleMultiSelect(val: any, checked: boolean) {
		const arr = Array.isArray(data) ? [...data] : [];
		if (checked) {
			data = [...arr, val];
		} else {
			data = arr.filter((v) => v !== val);
		}
	}
</script>

<BaseField {field}>
	{#snippet children({ id })}
		{#if field.type === 'select'}
			<select {id} bind:value={data} class="select-input">
				<option value="">Select an option…</option>
				{#each field.options ?? [] as opt}
					<option value={opt.value}>{opt.label}</option>
				{/each}
			</select>

		{:else if field.type === 'radio'}
			<div class="options-group">
				{#each field.options ?? [] as opt}
					<label class="radio-label">
						<input
							type="radio"
							name={field.key}
							value={opt.value}
							checked={data === opt.value}
							onchange={() => (data = opt.value)}
							class="radio-input"
						/>
						<span class="option-text">{opt.label}</span>
					</label>
				{/each}
			</div>

		{:else if field.type === 'checkbox'}
			<label class="checkbox-container">
				<input
					{id}
					type="checkbox"
					checked={!!data}
					onchange={handleCheckChanged}
					class="checkbox-input"
				/>
				<span class="option-text font-medium">{field.label}</span>
			</label>

		{:else if field.type === 'multiselect'}
			<div class="options-group">
				{#each field.options ?? [] as opt}
					<label class="radio-label">
						<input
							type="checkbox"
							checked={Array.isArray(data) && data.includes(opt.value)}
							onchange={(e) => handleMultiSelect(opt.value, (e.target as HTMLInputElement).checked)}
							class="checkbox-input"
						/>
						<span class="option-text">{opt.label}</span>
					</label>
				{/each}
			</div>
		{/if}
	{/snippet}
</BaseField>

<style>
	.select-input {
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

	.select-input:focus {
		outline: none;
		border-color: var(--brand);
	}

	.options-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.radio-label, .checkbox-container {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
	}

	.radio-input {
		accent-color: var(--brand);
	}

	.checkbox-input {
		width: 1rem;
		height: 1rem;
		accent-color: var(--brand);
	}

	.option-text {
		font-size: 0.875rem;
		color: var(--text-secondary);
	}

	.option-text.font-medium {
		font-weight: 500;
		color: var(--text-primary);
	}
</style>
