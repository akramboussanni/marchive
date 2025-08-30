import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {

	const env = loadEnv(mode, process.cwd(), '');

	return {
		plugins: [sveltekit()],
		server: {
			proxy: {
				'/api': {
					target: env.BACKEND_URL || 'http://localhost:9520',
					changeOrigin: true
				}
			}
		},
		define: {
			'process.env.BACKEND_URL': JSON.stringify(env.BACKEND_URL || 'http://localhost:9520'),
			'process.env.APP_NAME': JSON.stringify(env.APP_NAME || 'Marchived')
		}
	};
});
