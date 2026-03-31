<script lang="ts">
	import '../app.css';
	import Sidebar from '$lib/shared/Sidebar.svelte';
	import Topbar from '$lib/shared/Topbar.svelte';
	import Toast from '$lib/shared/Toast.svelte';
	import { collectionsStore } from '$lib/stores/collections';
	import { theme } from '$lib/stores/theme';
	import { browser } from '$app/environment';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	let sidebarCollapsed = $state(false);

	// Load theme from store on mount
	$effect(() => {
		if (browser) {
			const t = localStorage.getItem('theme') ?? 'light';
			document.documentElement.classList.toggle('dark', t === 'dark');
		}
	});

	// Load collections into store
	$effect(() => {
		collectionsStore.load();
	});
</script>

<svelte:head>
	<title>GoHeadless CMS</title>
</svelte:head>

<div class="app-shell">
	<Sidebar bind:collapsed={sidebarCollapsed} />

	<div class="main-content">
		<Topbar />
		<main class="page-content">
			{@render children()}
		</main>
	</div>
</div>

<Toast />
