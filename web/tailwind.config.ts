import type { Config } from 'tailwindcss';
import { fontFamily } from 'tailwindcss/defaultTheme';

const config = {
	darkMode: 'class',
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			fontFamily: {
				sans: ['Inter', ...fontFamily.sans]
			},
			colors: {
				primary: {
					50: '#fceee7',
					100: '#f4cab4',
					200: '#efb08f',
					300: '#e88c5c',
					400: '#e3753d',
					500: '#dc530c',
					600: '#c84c0b',
					700: '#9c3b09',
					800: '#792e07',
					900: '#5c2305'
				},

				dark: {
					50: '#383838',
					100: '#353535',
					200: '#333333',
					300: '#2d2d2d',
					400: '#2c2c2c',
					500: '#272727',
					600: '#242424',
					700: '#222222',
					800: '#1d1d1d',
					900: '#121212'
				},

				violet: {
					50: '#f3ebfd',
					100: '#d9c1f9',
					200: '#c6a3f6',
					300: '#ac7af2',
					400: '#9c60f0',
					500: '#8338ec',
					600: '#7733d7',
					700: '#5d28a8',
					800: '#481f82',
					900: '#371863'
				},

				red: {
					50: '#fcf1f0',
					100: '#f5d3d1',
					200: '#f0bebb',
					300: '#e9a19c',
					400: '#e58e89',
					500: '#de726b',
					600: '#ca6861',
					700: '#9e514c',
					800: '#7a3f3b',
					900: '#5d302d'
				},

				blue: {
					50: '#eceffd',
					100: '#c5cefa',
					200: '#a9b6f7',
					300: '#8195f4',
					400: '#6981f1',
					500: '#4361ee',
					600: '#3d58d9',
					700: '#3045a9',
					800: '#253583',
					900: '#1c2964'
				},

				foundation: {
					50: '#eceffd',
					100: '#c5cefa',
					200: '#a9b6f7',
					300: '#8195f4',
					400: '#6981f1',
					500: '#4361ee',
					600: '#3d58d9',
					700: '#3045a9',
					800: '#253583',
					900: '#1c2964'
				}
			}
		}
	},
	plugins: []
};

export default config;
