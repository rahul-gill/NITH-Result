<script>
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';

	export let storageKey = 'attendance_tracker_ad'; // Unique key for localStorage

	let isVisible = false;

	// Check if the snackbar has been dismissed before
	onMount(() => {
		const dismissed = localStorage.getItem(storageKey) === 'true';
		if (!dismissed) {
			isVisible = true;
		}
	});

	// Dismiss the snackbar and save its state
	const dismissSnackbar = () => {
		isVisible = false;
		localStorage.setItem(storageKey, 'true'); // Save dismissal in localStorage
	};
</script>

{#if isVisible}
	<div
		class="snackbar"
		in:fly={{ y: 50, duration: 300, delay: 1500 }}
		out:fly={{ y: 50, duration: 300 }}
	>
		<a
			style="display: flex; align-items: center;"
			target="_blank" rel="noopener noreferrer"
			href="https://play.google.com/store/apps/details?id=com.github.rahul_gill.attendance"
		>
			<img src="https://play.google.com/favicon.ico" style="padding:0.2em" alt="play_store_icon" />
			<span>Attendance and bunk tracker app(Android) </span>
		</a>

		<button on:click={dismissSnackbar} id="xse" class="dark"> Dismiss </button>
	</div>
{/if}

<style>
	.snackbar {
		position: fixed;
		bottom: 1rem;
		left: 1rem;
		right: 1rem;
		max-width: 28rem;
		margin: 0 auto;
		padding: 1rem;
		border-radius: 0.5rem;
		box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 1rem;
		background-color: #065f46; /* Equivalent to green-950 */
		color: #ffffff; /* White text */
		transition: all 0.3s ease-in-out;
	}


	a {
		color: #ffffff;
		text-decoration: none;
    font-weight: 600;
		transition: all 0.3s ease-in-out;
	}
	a:link {
		color: #dd9999;
		text-decoration: none;
	}
	a:active {
		color: #dd9999;
		text-decoration: none;
	}
	a:visited {
		color: #ffffff;
		text-decoration: none;
	}
	a:hover {
		color: #dd9999;
		text-decoration: none;
	}


	.snackbar button {
		font-size: 0.875rem; /* Text-sm */
		font-weight: 600; /* Font-semibold */
		background: none;
		border: none;
		color: inherit;
		cursor: pointer;
	}

	.snackbar button:hover {
		color: #e5e7eb; /* Equivalent to gray-200 */
	}

	@media (min-width: 768px) {
		.snackbar {
			max-width: 36rem; /* Equivalent to md:max-w-md */
		}
	}

	.dark .snackbar {
		background-color: #d1fae5; /* Equivalent to green-200 */
		color: #000000; /* Dark text */
	}


	#xse {
		background: #000000;
		transition: all 0.3s ease-in-out;
	}

	#xse:hover {
		background: #111111;
	}
</style>
