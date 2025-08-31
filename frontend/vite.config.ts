import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');

	return {
		plugins: [sveltekit()],
		server: {
			port: 5173,
			host: 'localhost',
			proxy: {
				'/api': {
					target: env.BACKEND_URL || 'http://localhost:9520',
					changeOrigin: true
				}
			},
			// Enable hot reload in development
			hmr: mode === 'development' ? {
				overlay: true
			} : undefined,
			// Watch for file changes
			watch: mode === 'development' ? {
				usePolling: true,
				interval: 100
			} : undefined
		},
		define: {
			'process.env.BACKEND_URL': JSON.stringify(env.BACKEND_URL || 'http://localhost:9520'),
			'process.env.APP_NAME': JSON.stringify(env.APP_NAME || 'marchive'),
			'process.env.NODE_ENV': JSON.stringify(mode)
		},
		// Optimize for development
		optimizeDeps: {
			include: mode === 'development' ? ['svelte'] : []
		}
	};
});
