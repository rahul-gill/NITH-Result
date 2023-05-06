import { getFilteredAndSortedResults, resultConfig } from '$lib/server/db_ops';
import { SortingType } from '$lib/server/types';
import type { PageServerLoad } from './$types';
import type { Actions } from './$types';

export const load = (async () => {
    const list = getFilteredAndSortedResults("", false, SortingType.cgpi, '', '', 50, 0, 0, 11);

    console.log(list);
    return {
        list,
        branches: resultConfig.available_branches,
        bataches: resultConfig.available_batches,
        sortTypes: sortOptions
    };
}) satisfies PageServerLoad;

export const actions = {
    default: async ({ cookies, request }) => {
        const data = await request.formData();

        console.log(data);
        const search = data.get('search');
        const branch = data.get('branch');
        const batch = data.get('batch');
        const sortBy = data.get('sort');
        const ascendingSearch = data.get('sort_order') === "ASC"

        let sorting: SortingType = sortOptionToDb(sortBy!.toString());


        const list = getFilteredAndSortedResults(
            search!.toString(),
            ascendingSearch,
            sorting,
            branch!.toString() === resultConfig.available_branches[0] ? '' : branch!.toString(),
            batch!.toString() === resultConfig.available_batches[0] ? '' : batch!.toString(),
            50, 0, 0, 11);
        return { list };
    }
} satisfies Actions;

function sortOptionToDb(sortOption: string): SortingType {
    let sorting = SortingType.cgpi;
    switch (sortOption!.toString()) {
        case "Sort by CGPI": sorting = SortingType.cgpi; break;
        case "Sort by Roll Number": sorting = SortingType.rollNumber; break;
        case "Sort by Name": sorting = SortingType.name; break;
    }
    return sorting;
}

const sortOptions = ['Sort by CGPI', 'Sort by Roll Number', 'Sort by Name']