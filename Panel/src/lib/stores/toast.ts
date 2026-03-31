import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface Toast {
	id: string;
	type: ToastType;
	message: string;
	duration?: number;
}

function createToastStore() {
	const { subscribe, update } = writable<Toast[]>([]);

	function add(type: ToastType, message: string, duration = 4000) {
		const id = crypto.randomUUID();
		update((toasts) => [...toasts, { id, type, message, duration }]);
		if (duration > 0) {
			setTimeout(() => remove(id), duration);
		}
	}

	function remove(id: string) {
		update((toasts) => toasts.filter((t) => t.id !== id));
	}

	return {
		subscribe,
		success: (msg: string, dur?: number) => add('success', msg, dur),
		error: (msg: string, dur?: number) => add('error', msg, dur),
		info: (msg: string, dur?: number) => add('info', msg, dur),
		warning: (msg: string, dur?: number) => add('warning', msg, dur),
		remove
	};
}

export const toast = createToastStore();
