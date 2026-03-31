import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export type ThemeType = 'light' | 'dark' | 'system';

function getInitialTheme(): ThemeType {
	if (!browser) return 'system';
	const stored = localStorage.getItem('theme');
	if (stored === 'light' || stored === 'dark' || stored === 'system') return stored;
	return 'system';
}

export const theme = writable<ThemeType>(getInitialTheme());

function applyTheme(val: ThemeType) {
	if (!browser) return;
	const isDark =
		val === 'dark' ||
		(val === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches);

	const html = document.documentElement;
	if (isDark) {
		html.classList.add('dark');
	} else {
		html.classList.remove('dark');
	}
}

// Subscribe to store changes to apply theme and persist
theme.subscribe((val) => {
	if (!browser) return;
	localStorage.setItem('theme', val);
	applyTheme(val);
});

// React to system preference changes if in system mode
if (browser) {
	window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
		theme.update((t) => {
			applyTheme(t);
			return t;
		});
	});
}

// Cycle through the available themes
export function toggleTheme() {
	theme.update((t) => {
		if (t === 'light') return 'dark';
		if (t === 'dark') return 'system';
		return 'light';
	});
}
