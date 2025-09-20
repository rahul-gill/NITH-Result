<script lang="ts">
	import { json } from '@sveltejs/kit';
	import type { PageData } from './$types';
	import { onMount } from 'svelte';

	import SEOPage from '../../components/SEOPage.svelte';
	export let data: PageData;

	let student = data.student;

	let isMobile = false;

	function handleResize(event: Event) {
		isMobile = window.innerWidth < 768;
	}

	onMount(() => {
		window.addEventListener('resize', handleResize);
		return () => {
			window.removeEventListener('resize', handleResize);
		};
	});
</script>


<SEOPage title="NITH Results for {data.student.roll_number} {data.student.name}" 
	description="NIT Hamirpur semester results for {data.student.roll_number} {data.student.name} "
	canonical="https://result-nith-rg.vercel.app/{data.student.roll_number}">
<div class="container">
	<div class="student-details">
		Name: {student.name}<br />
		Roll Number: {student.roll_number}<br />
		Batch: {student.batch}<br />
		Branch: {student.branch}<br />
		CGPI: {student.cgpi}
	</div>
	{#each student.semester_results as semester}
		<h1># Sem {semester.semester_number}</h1>
		<h3>CGPI: {semester.cgpi}, SGPI: {semester.sgpi}</h3>

		<table>
			<thead>
				<tr>
					<th>Subject name</th>
					<th>Subject code</th>
					<th>Points</th>
					<th>Grade</th>
					<th>SubGP</th>
				</tr>
			</thead>
			<tbody>
				{#each semester.subject_results as subject}
					<tr>
						<td data-label="Subject name">{subject.subject_name}</td>
						<td data-label="Subject code">{subject.subject_code}</td>
						<td data-label="Points">{subject.sub_point}</td>
						<td data-label="Grade">{subject.grade}</td>
						<td data-label="SubGP">{subject.sub_gp}</td>
					</tr>
				{/each}
			</tbody>
		</table>
		<!-- {/if} -->
	{/each}
</div>
</SEOPage>

<style>
	.container {
		text-align: center;
	}
	.student-details {
		border-radius: 10px;
		color: var(--on-container);
		margin: 20px;
		padding: 20px;
		text-align: left;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
		background-color: var(--container);
		transition: all 0.2s ease-in-out;
		transform: scale(1);
	}
	.student-details:hover {
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
		transform: scale(1.05);
		cursor: pointer;
	}
	table {
		border-collapse: collapse;
		width: 100%;
	}

	th,
	td {
		text-align: left;
		padding: 8px;
		border-bottom: 1px solid var(--on-background);
	}

	th {
		background-color: var(--on-container);
		color: var(--container);
	}

	tr:hover {
		background-color: var(--container);
	}

	@media screen and (max-width: 600px) {
		table thead {
			border: none;
			clip: rect(0 0 0 0);
			height: 1px;
			margin: -1px;
			overflow: hidden;
			padding: 0;
			position: absolute;
			width: 1px;
		}

		table tr {
			border-bottom: 3px solid var(--on-background);
			display: block;
			margin-bottom: 0.625em;
		}

		table td {
			border-bottom: 1px solid var(--on-background);
			display: block;
			font-size: 0.8em;
			text-align: right;
		}

		table td::before {
			content: attr(data-label);
			float: left;
			font-weight: bold;
			text-transform: uppercase;
		}

		table td:last-child {
			border-bottom: 0;
		}
	}
</style>
