<script lang="ts">
	import { getBasicNodeData, type NodeType } from '$lib/api/node';
	import SocialButton from './social-button.svelte';

	export let node: string;
	const nodeInfo = getBasicNodeData(node as NodeType);
</script>

<h2 class="mb-3 text-xl font-medium leading-6 hidden" style="font-size: 2em; margin-bottom: 1em;">
	Node Dashboard
</h2>

<h4 class="mb-3 text-xl font-medium leading-6">Project</h4>
<p class="mb-8 pb-1.5 text-base font-medium">
	<img src={nodeInfo.logo} alt="{nodeInfo.node} Logo" width="150" class="object-contain" />
</p>

<h4 class="mb-3 text-xl font-medium leading-6">Description</h4>
<p class="mb-8 pb-1.5 text-base font-medium text-justify">
	{nodeInfo.description}
</p>

<h4 class="mb-3 text-xl font-medium leading-6">Useful Links</h4>
<ul class="mb-[62px] flex gap-3">
	{#each nodeInfo.links as link}
		<SocialButton url={link} />
	{/each}
</ul>
