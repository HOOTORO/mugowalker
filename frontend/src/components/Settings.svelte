<script lang="ts">
    import { onMount } from "svelte";
    import { CurrentConfiguration } from "../lib/wailsjs/go/backend/Config.js";
    import * as rt from "../lib/wailsjs/runtime/runtime.js";
    import { settings } from "../stores/settings.js";
    function handleSubmit() {
        rt.EventsEmit("config:settings", $settings);
    }
    onMount(async () => {
        $settings = await CurrentConfiguration();
    });
</script>

<div class="component">
    <div>
        <h2>ImageMagick Args</h2>
        <form on:submit|preventDefault={handleSubmit}>
            <span
                ><input
                    bind:value={$settings.imagick.colorspace}
                />colorspace</span
            >
            <span><input bind:value={$settings.imagick.alpha} />alpha</span>
            <span
                ><input
                    bind:value={$settings.imagick.threshold}
                />threshold</span
            >
            <span><input bind:value={$settings.imagick.edge} />edge</span>
            <span><input bind:value={$settings.imagick.negate} />negate</span>
            <span
                ><input
                    bind:value={$settings.imagick.blackthreshold}
                />blackthreshold</span
            >

            <button type="submit"> Save </button>
        </form>
    </div>
    <div>
        <h2>Tesseract Args</h2>
        <form on:submit|preventDefault={handleSubmit}>
            <span><input bind:value={$settings.tesseract.psm} />psm</span>
            <span><input bind:value={$settings.tesseract.args} />args</span>

            <button type="submit"> Save </button>
        </form>
    </div>
</div>

<style lang="scss">
    form {
        display: block;
        input {
            margin: 5px;
        }
        button {
            background-color: chartreuse;
            border-radius: 10px;
            width: 90%;
        }
    }
</style>
