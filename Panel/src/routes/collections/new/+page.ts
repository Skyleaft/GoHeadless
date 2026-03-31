import type { PageLoad } from './$types';

export const load: PageLoad = async () => {
    // This is a static "New" page, no data needs to be fetched,
    // but having a +page.ts enables SvelteKit's preloading and 
    // SPA navigation optimizations for this route.
    return {};
};
