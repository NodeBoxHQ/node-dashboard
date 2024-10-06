<script lang="ts">
	import TwitterIcon from 'virtual:icons/prime/twitter';
	import DiscordIcon from 'virtual:icons/mdi/discord';
	import TelegramIcon from 'virtual:icons/mdi/telegram';
	import YoutubeIcon from 'virtual:icons/mdi/youtube';
	import GlobeIcon from 'virtual:icons/mdi/globe';

	export let url: string = '';

	let Icon: typeof TwitterIcon;

	if (url.includes('twitter') || url.includes('x.com')) {
		Icon = TwitterIcon;
	} else if (url.includes('discord')) {
		Icon = DiscordIcon;
	} else if (url.includes('telegram')) {
		Icon = TelegramIcon;
	} else if (url.includes('youtube')) {
		Icon = YoutubeIcon;
	} else {
		Icon = GlobeIcon;
	}
</script>

<li>
	<a
		href={url}
		class="flex size-12 items-center justify-center rounded-xl border border-dark-600 bg-dark-800 transition-all hover:bg-dark-600"
		target="_blank"
	>
		<svelte:component this={Icon} class="w-8 h-8" />
	</a>
</li>
