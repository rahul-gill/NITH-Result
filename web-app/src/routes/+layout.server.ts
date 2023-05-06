import { resultConfig } from '$lib/server/db_ops';
import type { LayoutServerLoad } from './$types';

export const load = (async () => {
    return {
        date: resultConfig.last_update_date
    };
}) satisfies LayoutServerLoad;