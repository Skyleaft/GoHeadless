<script lang="ts">
	import { toast } from '$lib/stores/toast';
	import type { Toast } from '$lib/stores/toast';

	const icons: Record<string, string> = {
		success: '✓',
		error: '✕',
		warning: '⚠',
		info: 'ℹ'
	};

	const styles: Record<string, string> = {
		success: 'border-l-4 border-green-500 bg-green-50 text-green-800 dark:bg-green-900/20 dark:text-green-300',
		error: 'border-l-4 border-red-500 bg-red-50 text-red-800 dark:bg-red-900/20 dark:text-red-300',
		warning: 'border-l-4 border-amber-500 bg-amber-50 text-amber-800 dark:bg-amber-900/20 dark:text-amber-300',
		info: 'border-l-4 border-blue-500 bg-blue-50 text-blue-800 dark:bg-blue-900/20 dark:text-blue-300'
	};

	let toasts = $derived([...$toast]);
</script>

<div class="toast-container">
	{#each toasts as t (t.id)}
		<div
			class="pointer-events-auto flex min-w-72 max-w-sm items-start gap-3 rounded-xl p-4 shadow-lg {styles[t.type]} animate-fade-in"
			role="alert"
		>
			<span class="flex-shrink-0 text-base font-bold">{icons[t.type]}</span>
			<p class="flex-1 text-sm font-medium">{t.message}</p>
			<button
				onclick={() => toast.remove(t.id)}
				class="flex-shrink-0 text-base opacity-60 hover:opacity-100 transition-opacity"
				aria-label="Dismiss"
			>
				✕
			</button>
		</div>
	{/each}
</div>

<style>
	@keyframes fade-in {
		from { opacity: 0; transform: translateX(16px); }
		to   { opacity: 1; transform: translateX(0); }
	}
	.animate-fade-in { animation: fade-in 0.2s ease forwards; }
</style>
