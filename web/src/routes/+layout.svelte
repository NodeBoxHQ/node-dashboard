<script lang="ts">
	import '../css/index.css';
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	// import { goto } from '$app/navigation';
	// import { sleep } from '$lib/utils';

	let isLoading: boolean = true;

	onMount(async () => {
		if (browser) {
			isLoading = false;
		}
	});
</script>

{#if isLoading}
	<h1 class="text-red-500">Loading...</h1>
{:else}
	<slot />
{/if}
