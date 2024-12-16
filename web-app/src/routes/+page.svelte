<script lang="ts">
	import { enhance } from '$app/forms';
	import type { PageData } from './$types';
	import SEOPage from '../components/SEOPage.svelte';
	import { initD, setInitDTrue, search, branch, batch, sort, ascDesc, resultList } from '../stores/query_params';
	import IntersectionObserver from '../components/IntersectionObserver.svelte';

	export let data: PageData;

	const branches = data.branches;
	const batches = data.bataches;
	const sorts = data.sortTypes;
	if(!initD){
		search.set('');
		branch.set(branches[0]);
		batch.set(batches[0]);
		sort.set(sorts[0]);
		resultList.set(data.list)
		setInitDTrue();
	}

	let timeoutId = setTimeout(() => {}, 300);
	function handleInput() {
		clearTimeout(timeoutId);
		timeoutId = setTimeout(handleSubmit, 300);
	}

	let currentPage = -1
	let submitButtonRef: HTMLButtonElement;
	function handleSubmit() {
		currentPage = -1;
		submitButtonRef.click();
	}
	async function onLoadMore(){
		if(currentPage == -1){
			currentPage = 1;
		} else {
			currentPage += 1
		}
		const formData = new FormData();
					
        formData.append('search', $search);
        formData.append('branch', $branch);
        formData.append('batch', $batch);
        formData.append('sort', $sort);
        formData.append('sort_order', $ascDesc);
        formData.append('load_page', currentPage);
		const r = await fetch('/', {
			method: 'POST',
			body: formData
		});
		r.body
		const body = await r.json()
		const resultListX = (body).list ?? [];
		resultList.set($resultList.concat(...resultListX))
	}

	let isErrorToastVisible = false;
	function showErrorToast() {
		isErrorToastVisible = true;
		setTimeout(() => {
			isErrorToastVisible = false;
		}, 1000);
	}
</script>

<svelte:head>
    <title>NITH results</title> 
</svelte:head>

<SEOPage title="NITH Results" description="NIT Hamirpur semester results" canonical="result-nith-rg.vercel.app/">

<form
	id="root_form"
	on:submit|preventDefault={handleSubmit}
	method="POST"
	use:enhance={({ form, data, action, cancel, submitter }) => {
		return async ({ result, update }) => {
			if (result.type === 'success' || result.type === 'failure') {
				if (result.type == 'failure') {
					showErrorToast();
				} else {
					resultList.set(result.data?.list ?? []);
				}
			}
		};
	}}
>
	<input
		type="text"
		id="search"
		name="search"
		placeholder="Search by roll number or name"
		bind:value={$search}
		on:input={handleInput}
	/>

	<select id="branch" name="branch" bind:value={$branch} on:change={handleSubmit}>
		{#each branches as branchOption}
			<option value={branchOption}>{branchOption}</option>
		{/each}
	</select>

	<select id="batch" name="batch" bind:value={$batch} on:change={handleSubmit}>
		{#each batches as batchOption}
			<option value={batchOption}>{batchOption}</option>
		{/each}
	</select>

	<select id="sort" name="sort" bind:value={$sort} on:change={handleSubmit}>
		{#each sorts as sortOption}
			<option value={sortOption}>{sortOption}</option>
		{/each}
	</select>

	<select id="sort_order" name="sort_order" bind:value={$ascDesc} on:change={handleSubmit}>
		<option value="DESC">Descending</option>
		<option value="ASC">Ascending</option>
	</select>

	<button bind:this={submitButtonRef} type="submit" class="submit-button">Submit</button>
</form>
{#if $resultList.length == 0}
	<div class="no-results-message">No students found</div>
{:else}
	<div class="result-card-list">
		{#each $resultList as resultItem, index}
			<a class="result-card" href="/{resultItem.roll_number}">
				<div class="list-index">#{index + 1}</div>
				{resultItem.name}<br />
				RollNumber: {resultItem.roll_number} <br />
				Branch: {resultItem.branch} <br />
				Batch: {resultItem.batch} <br />
				CGPI: {resultItem.cgpi}
				<div class="ranks">
					Class #{resultItem.class_rank}<br />
					Branch #{resultItem.branch_rank}<br />
					Year #{resultItem.year_rank}<br />
				</div>
			</a>
		{/each}
	</div>
	<IntersectionObserver let:intersecting top={200} onIntersectingTrue={onLoadMore} />
{/if}


</SEOPage>

<style>
	form {
		display: flex;
		flex-wrap: wrap;
	}
	input {
		flex: 1;
	}
	input,
	select {
		border-radius: 10px;
		padding: 10px;
		margin-left: 8px;
		margin-right: 8px;
		margin-top: 8px;
		margin-bottom: 8px;
		border: 2px solid var(--outline);
		background-color: var(--container);
		color: var(--on-container);
		transition: all 0.2s ease-in-out;
		
	}
	input:focus,
	select:focus {
		outline: none;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
		transform: scale(1.02);
		border-color: var(--primary);
	}
	.submit-button {
		display: none;
	}

	.no-results-message {
		text-align: center;
		margin-top: 30px;
		width: calc(100vw - 16px);
		max-width: 1200px;
	}
	.result-card-list {
		display: flex;
		flex-wrap: wrap;
		margin-top: 8px;
		justify-content: space-evenly;
		width: calc(100vw - 16px);
		max-width: 1200px;
	}

	.result-card {
		position: relative;
		padding: 8px;
		margin-bottom: 20px;
		border-radius: 10px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
		background-color: var(--container);
		transition: all 0.2s ease-in-out;
		transform: scale(1);

		width: 95%;
	}

	.ranks {
		position: absolute;
		top: 36px;
		right: 8px;
		text-align: right;
		color: var(--primary);
		border-radius: 50%;
	}

	.result-card:hover {
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
		transform: scale(1.05);
		cursor: pointer;
	}

	.list-index {
		position: relative;
		top: 0;
		left: 0;
		padding: 3px 6px;
		font-size: 1.2em;
		font-weight: bold;
		color: var(--primary);
		border-radius: 50%;
	}

	a {
		color: inherit;
		text-decoration: none;
	}

	@media (max-width: 440px) {
		.result-card-list {
			font-size: 0.8em;
		}
	}

	@media (max-width: 350px) {
		.result-card-list {
			font-size: 0.6em;
		}
	}

	@media (min-width: 950px) {
		.result-card {
			width: calc(45% - 10px);
		}
	}

	@media (min-width: 992px) {
		.result-card {
			width: calc(40% - 10px);
		}
	}

</style>
