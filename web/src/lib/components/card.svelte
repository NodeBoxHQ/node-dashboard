<script lang="ts">
	export let title: string;
	export let value: string | number;

	let formattedValue = '';

	$: {
		if (typeof value === 'string' && value.length >= 16) {
			formattedValue = value.slice(0, 16) + '...';
		} else {
			formattedValue = value.toString();
		}
	}
</script>

{#if value !== ''}
	<div
			class="min-w-48 flex-1 space-y-5 rounded-2xl border-dark-600 bg-dark-800 p-3 transition-all hover:shadow-sm hover:ring-4 hover:ring-dark-600 md:min-w-64 md:max-w-[325px]"
	>
		<p class="text-sm font-medium leading-6">{title}</p>
		<p class="text-base font-bold leading-6 break-words" title={String(value)}>
			{formattedValue}
		</p>
	</div>
{/if}
