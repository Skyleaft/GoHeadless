<script lang="ts">
	import { onMount } from 'svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		open?: boolean;
		title?: string;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		onclose?: () => void;
		children: Snippet;
		footer?: Snippet;
	}

	let {
		open = $bindable(false),
		title,
		size = 'md',
		onclose,
		children,
		footer
	}: Props = $props();

	let dialog: HTMLDialogElement;

	const sizes = {
		sm: 'max-w-sm',
		md: 'max-w-lg',
		lg: 'max-w-2xl',
		xl: 'max-w-4xl'
	};

	$effect(() => {
		if (!dialog) return;
		if (open) {
			dialog.showModal();
		} else {
			dialog.close();
		}
	});

	function handleClose() {
		open = false;
		onclose?.();
	}

	function handleBackdrop(e: MouseEvent) {
		if (e.target === dialog) handleClose();
	}

	onMount(() => {
		dialog.addEventListener('cancel', handleClose);
		return () => dialog.removeEventListener('cancel', handleClose);
	});
</script>

<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<dialog
	bind:this={dialog}
	onclick={handleBackdrop}
	class="modal-dialog m-auto rounded-xl p-0 shadow-2xl border w-[calc(100%-2rem)] {sizes[size]} animate-scale-in"
	style="
		background: var(--surface);
		border-color: var(--border);
		color: var(--text-primary);
	"
	class:hidden={!open}
>
	<!-- Header -->
	{#if title}
		<div
			class="flex items-center justify-between px-6 py-4 border-b"
			style="border-color: var(--border)"
		>
			<h3 class="text-base font-semibold" style="color: var(--text-primary)">{title}</h3>
			<button
				onclick={handleClose}
				class="flex items-center justify-center w-8 h-8 rounded-lg text-lg transition-colors hover:bg-[--surface-alt]"
				style="color: var(--text-muted)"
				aria-label="Close modal"
			>
				✕
			</button>
		</div>
	{/if}

	<!-- Body -->
	<div class="p-6">
		{@render children()}
	</div>

	<!-- Footer -->
	{#if footer}
		<div
			class="flex items-center justify-end gap-3 px-6 py-4 border-t"
			style="border-color: var(--border); background: var(--surface-alt)"
		>
			{@render footer()}
		</div>
	{/if}
</dialog>

<style>
	.modal-dialog {
		position: fixed;
		inset: 0;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal-dialog::backdrop {
		background: rgb(0 0 0 / 0.5);
		backdrop-filter: blur(2px);
	}

	@keyframes scale-in {
		from {
			opacity: 0;
			transform: scale(0.95) translateY(-8px);
		}
		to {
			opacity: 1;
			transform: scale(1) translateY(0);
		}
	}

	.animate-scale-in {
		animation: scale-in 0.18s cubic-bezier(0.16, 1, 0.3, 1) forwards;
	}
</style>
