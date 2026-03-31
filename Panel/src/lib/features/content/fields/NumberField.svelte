<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import BaseField from './BaseField.svelte';

	interface Props {
		field: Field;
		data: number;
	}

	let { field, data = $bindable() }: Props = $props();

	// Rating hover state
	let hover = $state(0);
</script>

<BaseField {field}>
	{#snippet children({ id })}
		{#if field.type === 'slider'}
			<div class="slider-wrapper">
				<input
					{id}
					type="range"
					bind:value={data}
					min={field.validation?.min ?? 0}
					max={field.validation?.max ?? 100}
					class="slider"
				/>
				<span class="value-display">{data ?? 0}</span>
			</div>
		{:else if field.type === 'rating'}
			<div class="rating-wrapper">
				{#each [1, 2, 3, 4, 5] as star}
					<button
						type="button"
						onclick={() => (data = star)}
						onmouseenter={() => (hover = star)}
						onmouseleave={() => (hover = 0)}
						class="star-btn"
						aria-label="Rate {star} stars out of 5"
						style="color: {star <= (hover || data || 0) ? '#f59e0b' : 'var(--border)'}"
					>
						★
					</button>
				{/each}
			</div>
		{:else}
			<input
				{id}
				type="number"
				bind:value={data}
				placeholder={field.placeholder}
				min={field.validation?.min}
				max={field.validation?.max}
				class="number-input"
			/>
		{/if}
	{/snippet}
</BaseField>

<style>
	.number-input {
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

	.number-input:focus {
		outline: none;
		border-color: var(--brand);
	}

	.slider-wrapper {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.slider {
		flex: 1;
		accent-color: var(--brand);
	}

	.value-display {
		width: 2.5rem;
		text-align: right;
		font-size: 0.875rem;
		font-variant-numeric: tabular-nums;
		color: var(--text-primary);
	}

	.rating-wrapper {
		display: flex;
		gap: 0.25rem;
	}

	.star-btn {
		font-size: 1.5rem;
		transition: color 0.15s, transform 0.1s;
		cursor: pointer;
	}

	.star-btn:hover {
		transform: scale(1.1);
	}
</style>
