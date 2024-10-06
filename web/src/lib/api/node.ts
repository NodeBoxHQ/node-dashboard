import axios, { AxiosError } from 'axios';
import { z } from 'zod';

const LineaSchema = z.object({
	status: z.string(),
	currentHeight: z.number().int().positive(),
	maxHeight: z.number().int().positive()
});

const DuskSchema = z.object({
	status: z.string(),
	version: z.string(),
	currentHeight: z.number().int().positive()
});

export type Linea = z.infer<typeof LineaSchema>;
export type Dusk = z.infer<typeof DuskSchema>;

type NodeSchema = typeof LineaSchema | typeof DuskSchema;
export type NodeData = Linea | Dusk;
export type NodeType = 'pc' | 'linea' | 'dusk';

const schemaMap: Record<NodeType, NodeSchema> = {
	pc: LineaSchema,
	linea: LineaSchema,
	dusk: DuskSchema
};

const defaultDataMap: Record<NodeType, NodeData> = {
	pc: { status: '', currentHeight: 0, maxHeight: 0 },
	linea: { status: '', currentHeight: 0, maxHeight: 0 },
	dusk: { status: '', version: '', currentHeight: 0 }
};

export async function fetchNodeData<T extends NodeData>(nodeType: NodeType): Promise<T> {
	try {
		const { data } = await axios.get<unknown>('/api/node');

		const schema = schemaMap[nodeType];
		const result = schema.safeParse(data);

		if (!result.success) {
			console.error('Validation error:', result.error.issues);
			return defaultDataMap[nodeType] as T;
		}

		return result.data as T;
	} catch (error) {
		if (error instanceof AxiosError) {
			console.error('Error fetching stats:', error.message);
		} else {
			console.error('Unexpected error:', error);
		}
		return defaultDataMap[nodeType] as T;
	}
}

interface BasicNodeData {
	node: string;
	logo: string;
	description: string;
	links: string[];
}

export function getBasicNodeData(node: NodeType): BasicNodeData {
	const nodeData: Record<NodeType, BasicNodeData> = {
		pc: {
			node: 'Linea',
			logo: 'https://bucket.nodebox.cloud/linea-logo.png',
			description:
				'Linea is a network that scales the experience of Ethereum. Its out-of-the-box compatibility with the Ethereum Virtual Machine enables the deployment of already-existing applications, as well as the creation of new ones that would be too costly on Mainnet. It also enables the community to use those dapps, at a fraction of the cost, and at multiples the speed of Mainnet.',
			links: [
				'https://twitter.com/LineaBuild',
				'https://www.youtube.com/@LineaBuild',
				'https://discord.gg/linea',
				'https://linea.build/'
			]
		},
		linea: {
			node: 'Linea',
			logo: 'https://bucket.nodebox.cloud/linea-logo.png',
			description:
				'Linea is a network that scales the experience of Ethereum. Its out-of-the-box compatibility with the Ethereum Virtual Machine enables the deployment of already-existing applications, as well as the creation of new ones that would be too costly on Mainnet. It also enables the community to use those dapps, at a fraction of the cost, and at multiples the speed of Mainnet.',
			links: [
				'https://twitter.com/LineaBuild',
				'https://www.youtube.com/@LineaBuild',
				'https://discord.gg/linea',
				'https://linea.build/'
			]
		},
		dusk: {
			node: 'Dusk',
			logo: 'https://bucket.nodebox.cloud/dusk-logo.png',
			description: 'Dusk Node',
			links: [
				'https://twitter.com/DuskFoundation',
				'https://discord.gg/dusk-official',
				'https://dusk.network/'
			]
		}
	};

	return nodeData[node];
}
