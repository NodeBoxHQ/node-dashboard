import axios, { AxiosError } from 'axios';
import { z } from 'zod';

const HostSchema = z.object({
	id: z.number().int().positive(),
	hostname: z.string(),
	owner: z.string(),
	privateIpv4: z.string(),
	privateIpv6: z.string(),
	ipv4: z.string(),
	ipv6: z.string(),
	node: z.string()
});

export type Host = z.infer<typeof HostSchema>;

export async function getHost(): Promise<Host> {
	try {
		const { data } = await axios.get<unknown>(`/api/host`);
		const result = HostSchema.safeParse(data);

		if (!result.success) {
			console.error('Validation error:', result.error.issues);
			return {} as Host;
		}

		return result.data;
	} catch (error) {
		if (error instanceof AxiosError) {
			console.error('Error fetching stats:', error.message);
		} else {
			console.error('Unexpected error:', error);
		}
		return {} as Host;
	}
}
