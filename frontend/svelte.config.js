import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/kit/vite';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter({
			pages: 'build',
			assets: 'build',
			fallback: 'index.html',
			precompress: false,
			strict: false
		}),
		prerender: {
			handleHttpError: ({ path, referrer, message }) => {
				// Ignore missing favicon and other non-critical assets
				if (path.includes('favicon') || path.includes('.ico')) {
					return;
				}
				// Log other errors but don't fail the build
				console.warn(`Prerender warning for ${path}: ${message}`);
			}
		}
	}
};

export default config;
