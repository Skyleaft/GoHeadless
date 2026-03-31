import { collectionsApi } from '$lib/api/collections';
import { contentApi } from '$lib/api/content';
import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ params }) => {
	try {
		const [collection, record] = await Promise.all([
			collectionsApi.get(params.slug),
			contentApi.get(params.slug, params.id)
		]);

		return {
			collection,
			record,
			slug: params.slug,
			id: params.id
		};
	} catch (err: any) {
		throw error(404, err.message || 'Record not found');
	}
};
