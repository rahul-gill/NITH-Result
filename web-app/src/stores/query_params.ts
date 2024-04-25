import type { StudentResultCompact } from '$lib/server/types';
import { writable, type Writable } from 'svelte/store';

const sortOptions = ['Sort by CGPI', 'Sort by Roll Number', 'Sort by Name'];
export let initD = false;
export let setInitDTrue = () => { initD = true };
export let search: Writable<string | null> = writable(null);
export let branch: Writable<string | null> = writable(null);
export let batch: Writable<string | null> = writable(null);
export let sort: Writable<string | null> = writable(null);
export let ascDesc: Writable<string | null> = writable("DESC");

export let resultList: Writable<StudentResultCompact[] > = writable([]);