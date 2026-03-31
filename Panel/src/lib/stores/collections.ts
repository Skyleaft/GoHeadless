import { writable } from 'svelte/store';
import type { Collection } from '$lib/types/collection';
import { collectionsApi } from '$lib/api/collections';

function createCollectionsStore() {
	const { subscribe, set, update } = writable<Collection[]>([]);
	let loaded = false;

	async function load(force = false) {
		if (loaded && !force) return;
		const data = await collectionsApi.list();
		set(data ?? []);
		loaded = true;
	}

	function addLocal(c: Collection) {
		update((all) => [...all, c]);
	}

	function removeLocal(slug: string) {
		update((all) => all.filter((c) => c.slug !== slug));
	}

	function reset() {
		set([]);
		loaded = false;
	}

	return { subscribe, load, addLocal, removeLocal, reset };
}

export const collectionsStore = createCollectionsStore();
