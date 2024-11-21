<script lang="ts">
	import { getStats, type Stats } from '$lib/api/stats';
	import { getHost, type Host } from '$lib/api/host';
	import { onMount } from 'svelte';
	import ListIcon from 'virtual:icons/fluent/list-rtl-16-filled';
	import ClockIcon from 'virtual:icons/mdi/clock-time-four-outline';

	import SocialButton from '$lib/components/social-button.svelte';
	import LineChart from '$lib/components/line-chart-card.svelte';
	import Card from '$lib/components/card.svelte';
	import Status from '$lib/components/status.svelte';
	import NodeInfo from '$lib/components/node-info.svelte';

	import { writable } from 'svelte/store';
	import { uptimeToRelative, firstLetterUppercase } from '$lib/utils';
	import { fetchNodeData, type NodeData, type NodeType } from '$lib/api/node';

	const statsStore = writable<Stats[]>([]);
	const hostStore = writable<Host>({
		id: 0,
		hostname: '',
		owner: '',
		privateIpv4: '',
		privateIpv6: '',
		ipv4: '',
		ipv6: '',
		node: '',
		version: ''
	});

	let latestUptime = 0;
	let nodeData: NodeData | undefined;
	let hasNodeData = false;

	$: {
		const latestStats =
			$statsStore.length > 0
				? $statsStore[$statsStore.length - 1]
				: { cpu: 0, memory: 0, storage: 0, network: 0, uptime: 0 };
		latestUptime = latestStats.uptime;

		const latestHost = $hostStore || {
			id: 0,
			hostname: '',
			owner: '',
			privateIpv4: '',
			privateIpv6: '',
			ipv4: '',
			ipv6: '',
			node: ''
		};

		hostStore.set(latestHost);
	}

	onMount(async () => {
		const [fetchedStats, fetchedHost] = await Promise.all([getStats(), getHost()]);

		statsStore.set(fetchedStats);
		hostStore.set(fetchedHost);

		nodeData = await fetchNodeData($hostStore.node as NodeType);
		hasNodeData = true;
	});

	async function updateStats() {
		const fetchedStats = await getStats();
		const fetchedHost = await getHost();
		statsStore.set(fetchedStats);
		hostStore.set(fetchedHost);
	}

	onMount(() => {
		const interval = setInterval(updateStats, 5000);
		return () => clearInterval(interval);
	});

	$: isLinea = $hostStore.node.toLowerCase() === 'linea';
	$: isDusk = $hostStore.node.toLowerCase() === 'dusk';
	$: isJuneo = $hostStore.node.toLowerCase() === 'juneo';
	$: isHyperliquid = $hostStore.node.toLowerCase() === 'hyperliquid' || $hostStore.node === 'pc';
</script>

<svelte:head>
	<title>NodeBox - {firstLetterUppercase($hostStore.node)} Dashboard</title>
</svelte:head>

<div class="bg-black text-white">
	<div class="overflow-hidden px-6 pb-20 pt-8 sm:px-10 md:pb-[106px] lg:px-20 lg:pt-11">
		<header class="mx-auto flex max-w-7xl items-center justify-center gap-6">
			<a href="index.html" class="block w-fit">
				<img src="/img/nodebox-logo.png" alt="Nodebox Logo" width="235" class="object-contain" />
			</a>
			<button
				class="btn-toggle flex size-[62px] items-center justify-center rounded-2xl border border-dark-600 bg-dark-800 hidden"
			>
				<ListIcon class="w-12 h-12" />
			</button>
		</header>

		<main class="mx-auto max-w-7xl">
			<div class="pt-14 md:pt-[36px]">
				{#if hasNodeData}
					<NodeInfo node={$hostStore.node} />
				{/if}

				<div class="relative z-0">
					{#if hasNodeData && nodeData}
						<Status status={nodeData.status} />
					{/if}

					<div
						class="absolute -top-[68px] left-[217px] -z-10 hidden h-[218px] w-[343px] transform-gpu rounded-[100%] bg-[#D35678] opacity-60 blur-[450px] md:block"
					></div>
					<div
						class="absolute -top-[41px] right-[217px] -z-10 hidden h-[163px] w-[257px] transform-gpu rounded-[100%] bg-[#FFA900] opacity-60 blur-[450px] md:block"
					></div>
				</div>

				<div class="relative z-10 mt-5 flex flex-wrap gap-5">
					<Card title="Dashboard Version" value={$hostStore.version} />
					<Card title="Node" value={$hostStore.node} />
					<Card title="Owner" value={$hostStore.owner} />
					<Card title="Public IPv4" value={$hostStore.ipv4} />
					<Card title="Public IPv6" value={$hostStore.ipv6} />
					<Card title="Private IPv4" value={$hostStore.privateIpv4} />
					<Card title="Private IPv6" value={$hostStore.privateIpv6} />

					{#if hasNodeData && nodeData}
						{#if isLinea && 'currentHeight' in nodeData && 'maxHeight' in nodeData}
							{#if nodeData.status === 'Online' || nodeData.status === 'Syncing'}
								<Card title="Current Height" value={nodeData.currentHeight.toString()} />
								<Card title="Maximum Height" value={nodeData.maxHeight.toString()} />
							{/if}
						{:else if isDusk && 'currentHeight' in nodeData && 'version' in nodeData && 'stake' in nodeData}
							{#if nodeData.status === 'Online'}
								<Card title="Dusk Version" value={nodeData.version} />
								<Card title="Height" value={nodeData.currentHeight.toString()} />
								<Card title="Staking Address" value={nodeData.stake.stakingAddress} />
								<Card title="Eligible Stake" value="{nodeData.stake.eligibleStake}" />
								<Card title="Accumulated Rewards" value="{nodeData.stake.rewards}" />
							{/if}
						{:else if isJuneo && 'uptimePercentage' in nodeData && 'networkName' in nodeData && 'nodeId' in nodeData}
							<Card title="Network Name" value={`${nodeData.networkName}`} />
							<Card title="Node ID" value={nodeData.nodeId} />
							<Card title="Uptime Percentage" value={`${nodeData.uptimePercentage.toFixed(2)}%`} />
						{:else if isHyperliquid && 'applyDuration' in nodeData}
							{#if nodeData.status === 'Online'}
								<Card title="Hyperliquid Version" value={nodeData.version} />
								<Card title="Current Height" value={nodeData.currentHeight.toString()} />
								<Card title="Block Time" value={nodeData.blockTime} />
								<Card title="Apply Duration" value={nodeData.applyDuration.toString()} />
							{/if}
						{/if}
					{/if}

					<div
						class="absolute -top-2.5 left-1/2 -z-10 hidden h-[353px] w-[555px] -translate-x-1/2 transform-gpu rounded-[100%] bg-[#B121C7] opacity-60 blur-[300px] md:block"
					></div>
				</div>

				<div class="relative z-10 mt-5 grid grid-cols-1 gap-5 sm:grid-cols-2">
					<LineChart data={$statsStore} type="cpu" />
					<LineChart data={$statsStore} type="memory" />
					<LineChart data={$statsStore} type="storage" />
					<LineChart data={$statsStore} type="network" />

					<div class="rounded-2xl border border-dark-600 bg-dark-800 p-6">
						<div class="flex items-center justify-between gap-6">
							<div class="flex items-center gap-4">
								<ClockIcon class="w-10 h-10" />
								<div>
									<h3 class="text-xl font-semibold leading-9 sm:text-2xl">Uptime</h3>
								</div>
							</div>
							<p>
								{uptimeToRelative(latestUptime)}
							</p>
						</div>
					</div>
				</div>
			</div>
		</main>

		<footer class="mx-auto mt-14 max-w-7xl md:mt-[72px]">
			<a href="index.html" class="block w-fit" style="margin-bottom: 1.5em">
				<img src="/img/nodebox-logo.png" alt="Nodebox" width="235" class="object-contain" />
			</a>

			<div class="flex flex-col items-start justify-between gap-8 md:flex-row md:gap-10">
				<div class="w-full max-w-[720px]">
					<div class="mb-8">
						<p class="text-base font-medium">
							Nodes management made easy. Simplify your entry into the blockchain space with our
							simple Node-as-a-Service solution, enabling you to join the ongoing revolution in
							blockchain technology.
						</p>
						<p style="margin-top: 0.75em">Nodebox &copy; 2024</p>
					</div>
				</div>
				<div class="w-[168px]">
					<p class="mb-3 text-xl font-medium leading-6">Follow us</p>
					<ul class="flex gap-3">
						<SocialButton url="https://x.com/nodebox_cloud" />
						<SocialButton url="https://discord.gg/VgTbHXRn6k" />
						<SocialButton url="https://telegram.me/nodebox_cloud" />
					</ul>
				</div>
			</div>
		</footer>
	</div>
</div>
