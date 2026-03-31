<script lang="ts">
	import type { FieldType } from '$lib/types/collection';
	import { FIELD_TYPE_GROUPS } from '$lib/types/collection';
	import Badge from '$lib/shared/Badge.svelte';

	interface Props {
		value?: FieldType;
		onchange?: (type: FieldType) => void;
	}

	let { value = $bindable('text' as FieldType), onchange }: Props = $props();

	let open = $state(false);
	let el: HTMLDivElement;

	function select(type: FieldType) {
		value = type;
		onchange?.(type);
		open = false;
	}

	function handleOutside(e: MouseEvent) {
		if (!el?.contains(e.target as Node)) open = false;
	}
</script>

<svelte:window onclick={handleOutside} />

<div class="relative" bind:this={el}>
	<button
		type="button"
		onclick={() => (open = !open)}
		class="flex w-full items-center justify-between gap-2 rounded-lg border px-3 py-2 text-sm transition-all"
		style="
			background: var(--surface);
			border-color: {open ? 'var(--brand)' : 'var(--border)'};
			color: var(--text-primary);
		"
	>
		<Badge type={value} />
		<span style="color: var(--text-muted)" class="text-xs">{open ? '▲' : '▼'}</span>
	</button>

	{#if open}
		<div
			class="absolute left-0 top-full z-50 mt-1 w-72 rounded-xl border p-2 shadow-lg animate-scale-in overflow-y-auto"
			style="
				background: var(--surface);
				border-color: var(--border);
				max-height: 340px;
			"
		>
			{#each FIELD_TYPE_GROUPS as group}
				<div class="mb-2">
					<p class="section-title px-2 py-1">{group.label}</p>
					<div class="grid grid-cols-2 gap-1">
						{#each group.types as item}
							<button
								type="button"
								onclick={() => select(item.value)}
								class="flex items-center gap-2 rounded-lg px-2 py-1.5 text-left text-xs transition-all hover:bg-[--surface-alt]"
								style="
									color: var(--text-primary);
									background: {value === item.value ? 'var(--brand-light)' : 'transparent'};
									font-weight: {value === item.value ? '600' : '400'};
								"
							>
								<span>{item.icon}</span>
								<span>{item.label}</span>
							</button>
						{/each}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	@keyframes scale-in {
		from { opacity: 0; transform: scale(0.96) translateY(-4px); }
		to   { opacity: 1; transform: scale(1) translateY(0); }
	}
	.animate-scale-in { animation: scale-in 0.12s ease forwards; }
</style>
