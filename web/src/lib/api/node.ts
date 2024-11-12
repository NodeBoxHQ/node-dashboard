import axios, { AxiosError } from 'axios';
import { z } from 'zod';

const LineaSchema = z.object({
	status: z.string(),
	currentHeight: z.number(),
	maxHeight: z.number()
});

const DuskSchema = z.object({
	status: z.string(),
	version: z.string(),
	currentHeight: z.number(),
	stake: z.object({
		stakingAddress: z.string(),
		eligibleStake: z.number(),
		slashes: z.number(),
		hardSlashes: z.number(),
		rewards: z.number()
	}),
});

const JuneoSchema = z.object({
	nodeId: z.string(),
	status: z.string(),
	uptimePercentage: z.number(),
	networkName: z.string()
});

export type Linea = z.infer<typeof LineaSchema>;
export type Dusk = z.infer<typeof DuskSchema>;
export type Juneo = z.infer<typeof JuneoSchema>;

type NodeSchema = typeof LineaSchema | typeof DuskSchema | typeof JuneoSchema;
export type NodeData = Linea | Dusk | Juneo;
export type NodeType = 'pc' | 'linea' | 'dusk' | 'juneo';

const schemaMap: Record<NodeType, NodeSchema> = {
	pc: LineaSchema,
	linea: LineaSchema,
	dusk: DuskSchema,
	juneo: JuneoSchema
};

const defaultDataMap: Record<NodeType, NodeData> = {
	pc: { status: '', currentHeight: 0, maxHeight: 0 },
	linea: { status: '', currentHeight: 0, maxHeight: 0 },
	dusk: { status: '', version: '', currentHeight: 0, stake: { stakingAddress: '', eligibleStake: 0, slashes: 0, hardSlashes: 0, rewards: 0 } },
	juneo: { nodeId: '', status: '', uptimePercentage: 0, networkName: '' }
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
			description:
				'Dusk is a Layer-1 blockchain designed to provide institutional-level and privacy and compliance in order to make it possible for anybody to trade real-world assets from their wallet. Built for regulated and decentralized finance, Dusk aims to evolve the financial landscape by making it possible for regulated assets to be brought on-chain.',
			links: [
				'https://twitter.com/DuskFoundation',
				'https://discord.gg/dusk-official',
				'https://dusk.network/'
			]
		},
		juneo: {
			node: 'Juneo',
			logo: 'https://bucket.nodebox.cloud/juneo-logo.png',
			description:
				'A permissionless protocol, deriving its foundation from the Snowman++ version of the Avalanche DAG consensus.',
			links: [
				'https://twitter.com/JUNEO_official',
				'https://discord.gg/juneosupernet',
				'https://juneosupernet.com/'
			]
		}
	};

	return nodeData[node];
}
