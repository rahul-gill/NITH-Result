<script>
  import { onMount } from "svelte";
  import { fly } from "svelte/transition";

  export let storageKey = "announcementDismissed"; // Unique key for localStorage

  let isVisible = false;

  // Check if the snackbar has been dismissed before
  onMount(() => {
    const dismissed = localStorage.getItem(storageKey) === "true";
    if (!dismissed) {
      isVisible = true;
    }
  });

  // Dismiss the snackbar and save its state
  const dismissSnackbar = () => {
    isVisible = false;
    localStorage.setItem(storageKey, "true"); // Save dismissal in localStorage
  };
</script>

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

  .snackbar a {
    font-weight: 600;
    color: #5bb2f6; /* Equivalent to blue-300 */
    text-decoration: underline;
  }

  .snackbar a:visited {
    color: #1e40af; /* Fallback for blue links */
  }

  .snackbar a:hover {
    color: #1d4ed8; /* Darker blue */
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

  .dark .snackbar a {
    color: #2563eb; /* Equivalent to blue-900 */
  }
</style>

{#if isVisible}
  <div
    class="snackbar"
    transition:fly="{{ y: 50, duration: 300, delay: 1500 }}"
  >
    <span>
      To update result data / take responsibility for it, contact
      <a href="mailto:rgill1@protonmail.com?subject=To update result data or take responsibility for it">
        here
      </a>
    </span>
    <button on:click={dismissSnackbar}>
      Dismiss
    </button>
  </div>
{/if}
