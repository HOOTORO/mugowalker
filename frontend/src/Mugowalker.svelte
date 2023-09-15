<script lang="ts">
    import { onMount, afterUpdate, tick } from "svelte";
    import * as rt from "./lib/wailsjs/runtime/runtime.js";
    import Device from "./components/Device.svelte";
    // import { main } from "$lib/wailsjs/go/models.js";
    import { activity } from "./stores/activity.js";
    import { message } from "./stores/message.js";

    let loo: string;

    activity.subscribe((activity) => {
        loo = activity;
    });

    onMount(async () => {
        rt.EventsOn("message", (msg) => {
            rt.LogPrint("recivew msg!" + msg);
            $message = msg;
            activity.writeLog(msg);
        });
    });

    function quit() {
        rt.Quit();
    }
</script>

<div class="control" id="winbtn">
    <button class="btn" on:click={quit}>❌</button>
</div>

<div class="main">
    <Device />
    <!-- <h1 id="🤖">Account</h1> -->
    <!-- <Account /> -->
    <!--  <span>
        <h1 id="⚙️">Settings</h1>
        <Settings />
    </span>
     <span> <h1 id="❤️">Imaginer</h1></span> -->

    <h1 id="🌝">ACTIVITY</h1>
    <p class="lof">{@html loo}</p>
</div>

<style lang="scss">
    .main {
        display: block;
    }
    .lof {
        position: absolute;
        padding-left: 10px;
        font-size: small;
        text-align: left;
        border: 2px outset #8400ff;
        opacity: 0.9;
        border-radius: 5px;
        background-color: transparent;
        text-wrap: nowrap;
    }
</style>
