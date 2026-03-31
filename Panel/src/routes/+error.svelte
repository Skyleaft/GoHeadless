<script lang="ts">
	import { page } from '$app/state';
</script>

<svelte:head>
	<title>Error — GoHeadless CMS</title>
</svelte:head>

<div class="flex flex-col items-center justify-center gap-6 py-24 text-center">
	<div
		class="flex h-24 w-24 items-center justify-center rounded-3xl text-5xl"
		style="background: var(--surface-alt)"
	>
		⚠
	</div>
	<div class="flex flex-col gap-2">
		<h1 class="text-3xl font-bold" style="color: var(--text-primary)">
			{page.status ?? 'Error'}
		</h1>
		<p class="text-base" style="color: var(--text-secondary)">
			{page.error?.message ?? 'Something went wrong.'}
		</p>
	</div>
	<a
		href="/"
		class="inline-flex items-center gap-2 rounded-xl px-5 py-2.5 text-sm font-medium text-white"
		style="background: var(--brand)"
	>← Back to Dashboard</a>
</div>
