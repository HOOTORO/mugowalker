<script lang="ts">
    import * as rt from "./lib/wailsjs/runtime/runtime.js";
    import Account from "./components/Account.svelte";
    import Device from "./components/Device.svelte";
    import Settings from "./components/Settings.svelte";
    // import { main } from "$lib/wailsjs/go/models.js";
    // import { account } from "./stores/account.js";
    // import { settings } from "./stores/settings.js";
    // import { device } from "./stores/device.js";

    function quit() {
        rt.Quit();
    }
</script>

<div class="control" id="winbtn">
    <button class="btn" on:click={quit}>❌</button>
</div>

<div class="main">
    <h1><span>📲</span>Device</h1>
    <Device />
    <h1><span>🤖</span>Account</h1>
    <Account />
    <h1><span>⚙️</span>Settings</h1>
    <Settings />
    <h1><span>❤️</span>Imaginer</h1>
</div>

<style lang="scss">
    .main {
        display: block;
    }

    :global(.compo) {
        display: flex;
        padding: 1em;
        border: 0.1px solid #00d7ff;
        border-radius: 5px;
        justify-content: space-between;
    }
</style>
