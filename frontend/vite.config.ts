import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'

export default defineConfig({
        build: {
                sourcemap: false,
                minify: 'terser',
                terserOptions: {
                        compress: {
                                drop_console: true,
                                drop_debugger: true
                        }
                }      
        },
        plugins: [sveltekit()],
        resolve: {
                alias: {
                        $styles: new URL('./src/styles', import.meta.url).pathname,
                },
        },
        server: {
                open: true,
		headers: {
		  'Cross-Origin-Opener-Policy': 'same-origin',
		  'Cross-Origin-Embedder-Policy': 'require-corp'
		},
		fs: {
			allow: ['../..']
		}
	  }	
})