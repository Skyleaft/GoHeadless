<script lang="ts">
	import type { Snippet } from 'svelte';
	import Modal from './Modal.svelte';
	import Button from './Button.svelte';

	interface Props {
		open?: boolean;
		title?: string;
		message?: string;
		confirmLabel?: string;
		cancelLabel?: string;
		loading?: boolean;
		onconfirm?: () => void;
		oncancel?: () => void;
		children?: Snippet;
	}

	let {
		open = $bindable(false),
		title = 'Are you sure?',
		message = 'This action cannot be undone.',
		confirmLabel = 'Confirm',
		cancelLabel = 'Cancel',
		loading = false,
		onconfirm,
		oncancel
	}: Props = $props();

	function cancel() {
		open = false;
		oncancel?.();
	}
	function confirm() {
		onconfirm?.();
	}
</script>

<Modal bind:open size="sm" {title} onclose={cancel}>
	<p class="text-sm" style="color: var(--text-secondary)">{message}</p>

	{#snippet footer()}
		<Button variant="ghost" onclick={cancel}>{cancelLabel}</Button>
		<Button variant="danger" {loading} onclick={confirm}>{confirmLabel}</Button>
	{/snippet}
</Modal>
