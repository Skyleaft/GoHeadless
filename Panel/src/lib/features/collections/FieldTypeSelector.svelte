<script lang="ts">
	import type { FieldType } from '$lib/types/collection';
	import { FIELD_TYPE_GROUPS } from '$lib/types/collection';
	import Badge from '$lib/shared/Badge.svelte';
	import { onMount, tick } from 'svelte';

	interface Props {
		value?: FieldType;
		onchange?: (type: FieldType) => void;
	}

	let { value = $bindable('text' as FieldType), onchange }: Props = $props();

	let open = $state(false);
	let el: HTMLDivElement;
	let dropdownEl: HTMLDivElement;
	let dropdownPortal: HTMLDivElement;

	function select(type: FieldType) {
		value = type;
		onchange?.(type);
		open = false;
	}

	function handleOutside(e: MouseEvent) {
		if (!el?.contains(e.target as Node)) open = false;
	}

	function positionDropdown() {
		if (!dropdownEl || !el) return;
		const rect = el.getBoundingClientRect();
		dropdownEl.style.position = 'fixed';
		dropdownEl.style.left = `${rect.left}px`;
		dropdownEl.style.top = `${rect.bottom + 4}px`;
		dropdownEl.style.zIndex = '9999';
	}

	$effect(() => {
		if (open) {
			tick().then(() => {
				if (dropdownEl) {
					if (!dropdownPortal) {
						dropdownPortal = document.createElement('div');
						dropdownPortal.id = 'field-type-dropdown-portal';
						document.body.appendChild(dropdownPortal);
					}
					dropdownPortal.appendChild(dropdownEl);
					positionDropdown();
				}
			});
		} else if (dropdownPortal && dropdownEl) {
			if (dropdownEl.parentNode === dropdownPortal) {
				dropdownPortal.appendChild(dropdownEl);
			}
		}
	});

	onMount(() => {
		window.addEventListener('scroll', positionDropdown, true);
		window.addEventListener('resize', positionDropdown);
		return () => {
			window.removeEventListener('scroll', positionDropdown, true);
			window.removeEventListener('resize', positionDropdown);
			if (dropdownPortal && dropdownPortal.parentNode) {
				dropdownPortal.remove();
			}
		};
	});
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
</div>

{#if open}
	<div
		bind:this={dropdownEl}
		class="w-72 rounded-xl border p-2 shadow-lg animate-scale-in overflow-y-auto"
		style="
			background: var(--surface);
			border-color: var(--border);
			max-height: 340px;
		"
	>
		{#each FIELD_TYPE_GROUPS as group (group.label)}
			<div class="mb-2">
				<p class="section-title px-2 py-1">{group.label}</p>
				<div class="grid grid-cols-2 gap-1">
					{#each group.types as item (item.value)}
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

<style>
	@keyframes scale-in {
		from { opacity: 0; transform: scale(0.96) translateY(-4px); }
		to   { opacity: 1; transform: scale(1) translateY(0); }
	}
	.animate-scale-in { animation: scale-in 0.12s ease forwards; }
</style>