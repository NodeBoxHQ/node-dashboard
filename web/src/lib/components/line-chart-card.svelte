<script lang="ts">
	import { Line } from 'svelte-chartjs';
	import type { Stats } from '../api/stats';
	import { bytesToGiB, convertStampToDateTime } from '$lib/utils';
	import {
		Chart as ChartJS,
		Title,
		Tooltip,
		LineElement,
		LinearScale,
		PointElement,
		CategoryScale,
		Filler,
		type TooltipItem
	} from 'chart.js';
	import { onMount, afterUpdate, onDestroy } from 'svelte';

	import CPUIcon from 'virtual:icons/ph/cpu';
	import RamIcon from 'virtual:icons/ri/ram-line';
	import StorageIcon from 'virtual:icons/mdi/storage';
	import NetworkIcon from 'virtual:icons/mdi/network';

	ChartJS.register(Title, Tooltip, LineElement, LinearScale, PointElement, CategoryScale, Filler);

	export let data: Stats[] = [];
	export let type: keyof Stats = 'cpu';

	let chart: ChartJS<'line', number[]>;
	let prevDataLength = 0;
	let chartContainer: HTMLDivElement;

	function getLabelType(key: string): string {
		switch (key) {
			case 'cpu':
				return 'CPU Usage';
			case 'memory':
				return 'Memory Usage';
			case 'storage':
				return 'Disk Usage';
			case 'network':
				return 'Network Usage';
			default:
				return 'Unknown';
		}
	}

	const chartData = {
		labels: [] as string[],
		datasets: [
			{
				label: getLabelType(type),
				data: [] as number[],
				borderColor: '#57F287',
				backgroundColor: (context: { chart: { ctx: CanvasRenderingContext2D } }) => {
					const ctx = context.chart.ctx;
					const gradient = ctx.createLinearGradient(0, 0, 0, 66);
					gradient.addColorStop(0, 'rgba(255, 255, 255, 0.2)');
					gradient.addColorStop(1, 'rgba(255, 255, 255, 0)');
					return gradient;
				},
				borderWidth: 2,
				tension: 0.4,
				fill: true,
				pointRadius: 0
			}
		]
	};

	const options = {
		responsive: true,
		maintainAspectRatio: false,
		scales: {
			y: {
				display: false
			},
			x: {
				display: false
			}
		},
		plugins: {
			legend: {
				display: false
			},
			tooltip: {
				enabled: true,
				mode: 'index' as const,
				intersect: false,
				callbacks: {
					label: function (tooltipItem: TooltipItem<'line'>) {
						const value = (tooltipItem.dataset.data[tooltipItem.dataIndex] as number) ?? 0;
						const set = tooltipItem.dataset.label;

						if (set === 'CPU Usage' || set === 'Memory Usage' || set === 'Disk Usage') {
							return `${set}: ${value.toFixed(2)}%`;
						} else {
							return `${set}: ${bytesToGiB(value)} GiB`;
						}
					}
				}
			}
		},
		interaction: {
			intersect: false,
			mode: 'index' as const
		}
	};

	onMount(() => {
		if (data.length > 0) {
			updateChartData();
		}
		window.addEventListener('resize', handleResize);
		handleResize(); // Initial resize
	});

	onDestroy(() => {
		window.removeEventListener('resize', handleResize);
	});

	afterUpdate(() => {
		if (data.length > prevDataLength) {
			updateChartData();
			prevDataLength = data.length;
		}
	});

	function updateChartData() {
		if (!chart) return;

		const newData = data.slice(prevDataLength);
		newData.forEach((item) => {
			chartData.labels.push(convertStampToDateTime(item.createdAt));
			chartData.datasets[0].data.push(item[type]);
		});

		const maxDataPoints = 20;
		if (chartData.labels.length > maxDataPoints) {
			chartData.labels = chartData.labels.slice(-maxDataPoints);
			chartData.datasets[0].data = chartData.datasets[0].data.slice(-maxDataPoints);
		}

		chart.update('none');
	}

	function handleResize() {
		if (chartContainer && chart) {
			const containerWidth = chartContainer.clientWidth;
			chart.resize(containerWidth, 50); // Fixed height of 50px
		}
	}

	let Icon: typeof CPUIcon;
	let text: string = '';

	$: {
		if (type === 'cpu') {
			Icon = CPUIcon;
			text = data.length > 0 ? `${data[data.length - 1].cpu}%` : '0%';
		} else if (type === 'memory') {
			Icon = RamIcon;
			text = data.length > 0 ? `${data[data.length - 1].memory}%` : '0%';
		} else if (type === 'storage') {
			Icon = StorageIcon;
			text = data.length > 0 ? `${data[data.length - 1].storage}%` : '0%';
		} else if (type === 'network') {
			Icon = NetworkIcon;
			text = data.length > 0 ? `${bytesToGiB(data[data.length - 1].network)} GiB` : '0 GiB';
		}
	}
</script>

<div class="rounded-2xl border border-dark-600 bg-dark-800 p-3">
	<div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 md:gap-6">
		<div class="flex items-center gap-4">
			<svelte:component this={Icon} class="w-10 h-10" />
			<div>
				<h3 class="text-xl font-semibold leading-9 sm:text-2xl">
					{getLabelType(type).replace('Usage', '')}
				</h3>
				<p
					class="mt-2 flex h-6 w-fit items-center justify-center rounded-full bg-[#ECFDF3] px-2.5 text-sm font-medium text-[#027A48]"
				>
					{text}
				</p>
			</div>
		</div>
		<div class="w-full md:w-[200px] flex-shrink-0" bind:this={chartContainer}>
			<div class="h-[50px]">
				<Line data={chartData} {options} bind:chart />
			</div>
		</div>
	</div>
</div>
