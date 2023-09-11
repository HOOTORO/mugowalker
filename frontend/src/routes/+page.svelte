<script lang="ts">
    import { Greet } from "$lib/wailsjs/go/main/App";
    let devIp = "127.0.0.1:5555";
    let result: Promise<string> | null = null;

    function connect() {
        result = Greet(devIp);
    }
</script>

<h1>Device</h1>

<form class="input-box">
    <input bind:value={devIp} />
    <button class="btn rainbow-button" on:click={connect}> CONNECT </button>
</form>
<div>
    <output>
        {#if result}
            {#await result}
                loading
            {:then value}
                {value}
            {/await}
        {/if}
    </output>
</div>

<style lang="scss">
    .input-box {
        display: flex;
        justify-content: space-evenly;
    }
    .rainbow-button {
        background-image: linear-gradient(
            90deg,
            #00d7ff 0%,
            #bada55 20%,
            #ffcf00 40%,
            #fc4f4f 60%,
            #ff00ee 80%,
            #00d7ff 100%
        );
        border-radius: 7px;
        padding: 5px;
        text-transform: uppercase;
        font-size: 1.3rem;
        font-weight: bold;
    }
    .rainbow-button:after {
        content: none;
        background-color: rgba(14, 14, 19, 1);
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 5px;
        padding: 5px;
    }

    .rainbow-button:hover {
        animation: slidebg 3s linear infinite;
    }

    output {
        margin: 2em;
        padding: 1em;
        display: block;
        text-align: center;
        border: 2px solid #00d7ff;
        border-radius: 5px;
    }
</style>
