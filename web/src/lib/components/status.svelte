<script lang="ts">
	import CheckIcon from 'virtual:icons/carbon/checkmark-filled';
	import CrossIcon from 'virtual:icons/carbon/close-filled';
	import SyncIcon from 'virtual:icons/eos-icons/blockchain';
	import QuestionIcon from 'virtual:icons/ph/question';

	export let status: string = 'Syncing';

	let formatted: string = '';
	let Icon: typeof CheckIcon;
	let iconClasses: string;

	if (status === 'Online') {
		formatted = 'Node is Running';
		Icon = CheckIcon;
		iconClasses = 'w-12 h-12 text-green-500';
	} else if (status === 'Offline') {
		formatted = 'Node is Down';
		Icon = CrossIcon;
		iconClasses = 'w-12 h-12 text-red-500';
	} else if (status === 'Syncing') {
		formatted = 'Node is Syncing';
		Icon = SyncIcon;
		iconClasses = 'w-12 h-12 text-yellow-400';
	} else {
		formatted = 'Node is in Unknown State';
		Icon = QuestionIcon;
		iconClasses = 'w-12 h-12 text-gray-500';
	}
</script>

<div class="flex flex-col items-center justify-center gap-4 rounded-2xl bg-white py-4 text-center">
	<svelte:component this={Icon} class={iconClasses} />
	<p class="text-xl font-medium leading-6 text-black">{formatted}</p>
</div>
