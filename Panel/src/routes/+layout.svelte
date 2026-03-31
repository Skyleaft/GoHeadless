<script lang="ts">
	import '../app.css';
	import { onMount, type Snippet } from 'svelte';
	import { app } from '$lib/app.svelte';
	import Sidebar from '$lib/shared/Sidebar.svelte';
	import Topbar from '$lib/shared/Topbar.svelte';
	import Toast from '$lib/shared/Toast.svelte';
	import Spinner from '$lib/shared/Spinner.svelte';

	let { children }: { children: Snippet } = $props();
	let sidebarCollapsed = $state(false);

	onMount(() => {
		app.initialize();
	});
</script>

<svelte:head>
	<title>GoHeadless CMS</title>
</svelte:head>

{#if !app.isInitialLoaded && !app.isAuthPage}
	<div class="inset-0 fixed z-[9999] flex items-center justify-center bg-[--bg]">
		<Spinner />
	</div>
{:else if app.isAuthPage}
	<main class="min-h-screen">
		{@render children()}
	</main>
{:else}
	<div class="app-shell bg-[--bg]">
		<Sidebar bind:collapsed={sidebarCollapsed} />

		<div
			class="main-content transition-all duration-300"
			style="margin-left: {sidebarCollapsed
				? 'calc(80px + 2rem)'
				: 'calc(var(--sidebar-width) + 2rem)'}"
		>
			<Topbar />
			<main class="page-content animate-fade-in">
				{@render children()}
			</main>
		</div>
	</div>
{/if}

<Toast />

<style>
	.app-shell {
		display: flex;
		height: 100vh;
		overflow: hidden;
	}

	.main-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		min-width: 0;
		overflow: hidden;
	}

	.page-content {
		flex: 1;
		overflow-y: auto;
		padding: 2rem;
		scroll-behavior: smooth;
	}
</style>
