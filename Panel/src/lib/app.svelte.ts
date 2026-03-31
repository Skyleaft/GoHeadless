import { browser } from '$app/environment';
import { page } from '$app/state';
import { goto } from '$app/navigation';
import { auth } from '$lib/stores/auth';
import { collectionsStore } from '$lib/stores/collections';
import { getSetupStatus } from '$lib/api/auth';
import { get } from 'svelte/store';

export class AppController {
	isInitialLoaded = $state(false);
	isChecking = $state(false);

	isAuthPage = $derived(
		page.url.pathname === '/login' ||
		page.url.pathname === '/setup'
	);

	currentPath = $derived(page.url.pathname);

	constructor() { }

	async initialize() {
		if (this.isChecking || !browser) return;
		this.isChecking = true;

		// Reactive theme management (Now safely in an effect context during init)
		$effect(() => {
			const theme = localStorage.getItem('theme') ?? 'light';
			document.documentElement.classList.toggle('dark', theme === 'dark');
		});

		try {
			const currentPath = page.url.pathname;
			const authenticated = get(auth).isAuthenticated;

			// 1. Initial System Check (Setup Required?)
			if (currentPath !== '/setup') {
				const status = await getSetupStatus();
				if (status.setup_required) {
					await goto('/setup', { replaceState: true });
					return;
				}
			}

			// 2. Auth Guard
			if (!authenticated && !this.isAuthPage) {
				await goto('/login', { replaceState: true });
				return;
			}

			// 3. User normalization (Prevent logged-in users from login/setup)
			if (authenticated && this.isAuthPage) {
				await goto('/', { replaceState: true });
				return;
			}

			// 4. Critical Data Bootstrap
			if (authenticated) {
				await collectionsStore.load();
			}

			this.isInitialLoaded = true;
		} catch (e) {
			console.error('[AppController] Initialization failed:', e);
		} finally {
			this.isChecking = false;
		}
	}
}

export const app = new AppController();
