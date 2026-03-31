<script lang="ts">
	let { value = $bindable([]), placeholder = "Add item...", onchange }: { 
		value: string[], 
		placeholder?: string,
		onchange?: (val: string[]) => void 
	} = $props();

	let inputValue = $state('');

	function addTag() {
		const trimmed = inputValue.trim();
		if (trimmed && !value.includes(trimmed)) {
			value = [...value, trimmed];
			onchange?.(value);
		}
		inputValue = '';
	}

	function removeTag(index: number) {
		value = value.filter((_, i) => i !== index);
		onchange?.(value);
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' || e.key === ',') {
			e.preventDefault();
			addTag();
		} else if (e.key === 'Backspace' && !inputValue && value.length > 0) {
			removeTag(value.length - 1);
		}
	}
</script>

<div class="tag-input-container group">
	<div class="flex flex-wrap gap-2 p-2 min-h-[3rem]">
		{#each value as tag, i}
			<div class="tag animate-scale-in">
				<span class="text-xs font-black tracking-tight">{tag}</span>
				<button 
					type="button" 
					onclick={() => removeTag(i)}
					class="ml-1 hover:text-red-500 transition-colors"
				>✕</button>
			</div>
		{/each}
		<input
			type="text"
			bind:value={inputValue}
			onkeydown={handleKeydown}
			onblur={addTag}
			{placeholder}
			class="flex-1 bg-transparent border-none outline-none text-sm font-bold min-w-[120px] px-2"
			style="color: var(--text-primary)"
		/>
	</div>
</div>

<style>
	.tag-input-container {
		background: var(--surface-alt);
		border: 2px solid transparent;
		border-radius: 1rem;
		transition: all 0.2s;
		box-shadow: inset 0 2px 4px rgba(0,0,0,0.05);
	}

	.tag-input-container:focus-within {
		border-color: var(--brand);
		background: var(--surface);
		box-shadow: 0 0 0 4px rgba(var(--brand-rgb), 0.1);
	}

	.tag {
		background: var(--brand);
		color: white;
		padding: 4px 10px;
		border-radius: 0.75rem;
		display: flex;
		align-items: center;
		gap: 4px;
		box-shadow: 0 4px 6px -1px rgba(var(--brand-rgb), 0.2);
	}

	@keyframes scale-in {
		from { transform: scale(0.9); opacity: 0; }
		to { transform: scale(1); opacity: 1; }
	}

	.animate-scale-in {
		animation: scale-in 0.2s ease-out forwards;
	}
</style>
