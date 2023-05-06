import { getStudent } from '$lib/server/db_ops';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';


export const load = (async ({ params }) => {
    const { rollno } = params;
    const student = getStudent(rollno);
    if (student === null) {
        throw error(404, "Roll number not found");
    }
    return {
        student
    };
}) satisfies PageServerLoad;