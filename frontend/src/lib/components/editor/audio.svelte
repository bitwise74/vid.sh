<script lang="ts">
        import { FFmpeg } from '@ffmpeg/ffmpeg'
        import { fetchFile } from '@ffmpeg/util'
        
	const baseURL = "https://cdn.jsdelivr.net/npm/@ffmpeg/core@0.12.10/dist/esm"
        
	let { file }: { file: File } = $props()
        
	let ffmpeg: FFmpeg | null = null
	let appState: 'idle' | 'loading' | 'ready' = $state('idle')
        let isFFmpegLoaded = $state(false)
        
        async function loadFFmpeg() {
                if (!ffmpeg) {
			appState = 'loading'
			ffmpeg = new FFmpeg()

			await ffmpeg.load({
				coreURL: `${baseURL}/ffmpeg-core.js`,
				wasmURL: `${baseURL}/ffmpeg-core.wasm`,
				workerURL: `${baseURL}/ffmpeg-core.worker.js`
			})
			appState = 'ready'
			console.log("FFmpeg loaded")
		}
        }

	async function processAudio() {
                if (!ffmpeg) return

		const url = URL.createObjectURL(file)
		await ffmpeg.writeFile("input.mp4", await fetchFile(url))
		await ffmpeg.exec(["-i", "input.mp4", "-map", "0:a:0", "track.mp3"])
		const data = await ffmpeg.readFile('track.mp3')
		URL.revokeObjectURL(url)

		const exp = URL.createObjectURL(new Blob([(data as any).buffer], {
			type: "audio/mp3"
		}))
		const a = document.createElement('a')
		a.href = exp
		a.download = "track.mp3"
		a.click()
	}
</script>

{#if appState === 'idle'}
<div>
        <p>This tab allows for advanced video editing, but requires you to load a pretty big bundle. Click on the button below to load it.</p>
        <button class="btn btn-outline-danger" onclick={loadFFmpeg}>Load ffmpeg.wasm (~31MB)</button>
</div>
{:else if appState === 'loading'}
<div class="spinner-border" role="status">
        <span class="sr-only">Loading ffmpeg.wasm</span>
</div>
        <p>Loading ffmpeg...</p>
{:else if appState === 'ready'}
        <p>Loaded :)</p>
{/if}

<!-- 
<button on:click={processAudio} disabled={state === 'loading'}>
	{#if state === 'idle'} Extract Audio 🎵 {/if}
	{#if state === 'loading'} Loading FFmpeg... {/if}
	{#if state === 'ready'} Extract Again 🎵 {/if}
</button> -->