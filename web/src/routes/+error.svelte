<script context="module">
	import { page } from '$app/stores';
</script>

<main class={`error-page h-screen flex flex-col justify-center items-center`}>
	<h1 class="text-8xl text-secondary-800 dark:text-surface-50 font-bold mb-2">{$page?.status}</h1>
	{#if $page?.status === 404}
		<p class="text-2xl text-secondary-800 dark:text-surface-50 font-semibold mb-2">
			{$page?.error?.message}
		</p>
	{:else}
		<p class="text-secondary-800 dark:text-surface-50">An unexpected error has occurred.</p>
	{/if}
	<a
		href="/"
		class="bg-primary-600 hover:bg-primary-400 dark:hover:bg-primary-700 text-white px-2 py-2 rounded-lg mb-2"
		>Go back to Home</a
	>
</main>
