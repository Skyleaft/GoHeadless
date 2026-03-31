<script lang="ts">
	interface Props {
		class?: string;
	}
	let { class: cls = '' }: Props = $props();
</script>

<div class="flex h-5 w-5 items-center justify-center {cls}">
	<span
		class="block h-4 w-4 rounded-full border-2 border-[--brand] border-t-transparent"
		style="animation: spin 0.7s linear infinite"
	></span>
</div>

<style>
	@keyframes spin { to { transform: rotate(360deg); } }
</style>
