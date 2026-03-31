<script lang="ts">
	import '../app.css';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { collectionsStore } from '$lib/stores/collections';
	import { getSetupStatus } from '$lib/api/auth';
	import type { Snippet } from 'svelte';
	
	import Sidebar from '$lib/shared/Sidebar.svelte';
	import Topbar from '$lib/shared/Topbar.svelte';
	import Toast from '$lib/shared/Toast.svelte';

	let { children }: { children: Snippet } = $props();

	let sidebarCollapsed = $state(false);

	let isAuthPage = $derived($page.url.pathname === '/login' || $page.url.pathname === '/setup');

	// Load theme from store on mount
	$effect(() => {
		if (browser) {
			const t = localStorage.getItem('theme') ?? 'light';
			document.documentElement.classList.toggle('dark', t === 'dark');
		}
	});

	// Auth Guard & Setup initialization check
	$effect(() => {
		if (!browser) return;

		const checkAuth = async () => {
			// 1. Check if setup is required
			if ($page.url.pathname !== '/setup') {
				try {
					const status = await getSetupStatus();
					if (status.setup_required) {
						goto('/setup');
						return;
					}
				} catch (e) {
					console.error('Failed to check setup status', e);
				}
			}

			// 2. Check if logged in
			if (!isAuthPage && !$auth.isAuthenticated) {
				goto('/login');
			}
			
			// 3. Load collections only if authenticated
			if ($auth.isAuthenticated) {
				collectionsStore.load();
			}
		};

		checkAuth();
	});
</script>

<svelte:head>
	<title>GoHeadless CMS</title>
</svelte:head>

{#if isAuthPage}
	<main class="min-h-screen">
		{@render children()}
	</main>
{:else}
	<div class="app-shell">
		<Sidebar bind:collapsed={sidebarCollapsed} />

		<div class="main-content">
			<Topbar />
			<main class="page-content">
				{@render children()}
			</main>
		</div>
	</div>
{/if}

<Toast />
