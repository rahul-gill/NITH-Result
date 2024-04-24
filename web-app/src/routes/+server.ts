
import { getFilteredAndSortedResults, resultConfig } from '$lib/server/db_ops';
import { SortingType } from '$lib/server/types';
import { json } from '@sveltejs/kit';
export async function POST({ request, cookies }) {
	const data = await request.formData();

        const search = data.get('search');
        const branch = data.get('branch');
        const batch = data.get('batch');
        const sortBy = data.get('sort');
        const ascendingSearch = data.get('sort_order') === "ASC";
        const loadPage = data.get('load_page');

        let sorting: SortingType = sortOptionToDb(sortBy!.toString());


        console.info("loadPagxZe", loadPage)
        const list = getFilteredAndSortedResults(
            search!.toString(),
            ascendingSearch,
            sorting,
            branch!.toString() === resultConfig.available_branches[0] ? '' : branch!.toString(),
            batch!.toString() === resultConfig.available_batches[0] ? '' : batch!.toString(),
            50,(loadPage == -1) ? 0 : loadPage  , 0, 11);
        return json({ list });
}



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