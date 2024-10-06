import axios, { AxiosError } from 'axios';
import { z } from 'zod';

const StatsSchema = z.object({
	id: z.number().int().positive(),
	cpu: z.number().min(0).max(100),
	memory: z.number().min(0),
	storage: z.number().min(0),
	network: z.number().min(0),
	uptime: z.number().int().positive(),
	createdAt: z.number().int().positive()
});

export type Stats = z.infer<typeof StatsSchema>;
const ResponseSchema = z.array(StatsSchema);

export async function getStats(): Promise<Stats[]> {
	try {
		const { data } = await axios.get<unknown>('/api/stats');

		const result = ResponseSchema.safeParse(data);

		if (!result.success) {
			console.error('Validation error:', result.error.issues);
			return [];
		}

		return result.data;
	} catch (error) {
		if (error instanceof AxiosError) {
			console.error('Error fetching stats:', error.message);
		} else {
			console.error('Unexpected error:', error);
		}
		return [];
	}
}
