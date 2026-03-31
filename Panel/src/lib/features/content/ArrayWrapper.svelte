<script lang="ts">
	import type { Field } from '$lib/types/collection';
	import Button from '$lib/shared/Button.svelte';

	let { value = $bindable([]), field, children }: { 
		value: any[], 
		field: Field,
		children?: import('svelte').Snippet 
	} = $props();

	function addItem() {
		// Default values based on type
		let defaultValue: any = '';
		if (field.type === 'number' || field.type === 'slider' || field.type === 'rating') defaultValue = 0;
		if (field.type === 'bool' || field.type === 'toggle') defaultValue = false;
		if (field.type === 'file' || field.type === 'image') defaultValue = '';
		
		value = [...value, defaultValue];
	}

	function removeItem(index: number) {
		value = value.filter((_, i) => i !== index);
	}

	function moveItem(index: number, direction: 'up' | 'down') {
		const newIndex = direction === 'up' ? index - 1 : index + 1;
		if (newIndex < 0 || newIndex >= value.length) return;
		
		const newList = [...value];
		[newList[index], newList[newIndex]] = [newList[newIndex], newList[index]];
		value = newList;
	}
</script>

<div class="array-container">
	<div class="flex flex-col gap-3">
		{#each value as item, i}
			<div class="array-item group animate-slide-in">
				<div class="flex-1">
					{@render children?.()}
				</div>
				<div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity pr-2">
					<button 
						type="button" 
						onclick={() => moveItem(i, 'up')}
						disabled={i === 0}
						class="control-btn disabled:opacity-20"
					>↑</button>
					<button 
						type="button" 
						onclick={() => moveItem(i, 'down')}
						disabled={i === value.length - 1}
						class="control-btn disabled:opacity-20"
					>↓</button>
					<button 
						type="button" 
						onclick={() => removeItem(i)}
						class="control-btn text-red-400 hover:bg-red-500/10 hover:text-red-500"
					>✕</button>
				</div>
			</div>
		{:else}
			<div class="py-8 rounded-2xl flex flex-col items-center bg-slate-50/50 dark:bg-slate-900/20 border-2 border-dashed border-[var(--border)] opacity-40">
				<span class="text-3xl mb-2">📋</span>
				<p class="text-xs font-black uppercase tracking-widest">No entries yet</p>
			</div>
		{/each}
	</div>

	<div class="mt-4 flex justify-start">
		<Button variant="ghost" size="sm" onclick={addItem} class="hover:bg-[var(--brand)]/10 text-[var(--brand)]">
			+ Add {field.label || 'Item'}
		</Button>
	</div>
</div>

<style>
	.array-container {
		width: 100%;
	}

	.array-item {
		display: flex;
		align-items: flex-start;
		gap: 12px;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 1.25rem;
		padding: 4px;
		transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: var(--shadow-sm);
	}

	.array-item:hover {
		border-color: var(--brand);
		box-shadow: var(--shadow-md);
		transform: translateY(-1px);
	}

	.control-btn {
		height: 28px;
		width: 28px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 0.75rem;
		font-size: 14px;
		font-weight: 900;
		transition: all 0.2s;
		color: var(--text-muted);
	}

	.control-btn:hover:not(:disabled) {
		background: var(--surface-alt);
		color: var(--text-primary);
	}

	@keyframes slide-in {
		from { transform: translateX(-10px); opacity: 0; }
		to { transform: translateX(0); opacity: 1; }
	}

	.animate-slide-in {
		animation: slide-in 0.2s ease-out forwards;
	}
</style>
