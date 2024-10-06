import { DateTime } from 'luxon';

export async function sleep(ms: number) {
	return new Promise((resolve) => setTimeout(resolve, ms));
}

export function convertStampToDateTime(stamp: number): string {
	return DateTime.fromMillis(stamp * 1000).toFormat('dd-MM-yy HH:mm:ss');
}

export function bytesToGiB(bytes: number): number {
	return parseFloat((bytes / 1024 / 1024 / 1024).toFixed(2));
}

export function uptimeToRelative(seconds: number): string {
	const date = DateTime.local().minus({ seconds });
	return date.toRelative();
}
