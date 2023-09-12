<script lang="ts">
    import { device } from "../stores/device";
    import { Greet } from "../lib/wailsjs/go/main/App";

    let result: Promise<string> | null = null;

    function connect() {
        result = Greet($device.dest);
    }
</script>

<div class="compo">
    <input bind:value={$device.dest} />
    <button class="btn rainbow-button" on:click={connect}> CONNECT </button>
    <output>
        {#if result}
            {#await result}
                loading
            {:then value}
                <span style="color: green;text-align:right">Conected!</span>
                <!-- {value} -->
            {:catch error}
                <span style="color: red;text-align:right">
                    ERROR: {error.message}
                </span>
            {/await}
        {:else}
            <span style="color: red;text-align:right">Disconected</span>
        {/if}
    </output>
</div>
