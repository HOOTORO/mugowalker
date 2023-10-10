<script lang="ts">
    import { afterUpdate, onMount } from "svelte";
    import Device from "./components/Device.svelte";
    import Settings from "./components/Settings.svelte";
    import * as rt from "./lib/wailsjs/runtime/runtime.js";
    import { activity } from "./stores/activity.js";
    import { message } from "./stores/message.js";

    let loo: string;

    activity.subscribe((activity) => {
        loo = activity;
    });

    onMount(async () => {
        rt.EventsOn("message", (msg) => {
            rt.LogPrint(msg);
            $message = msg;
            activity.writeLog(msg);
        });

        rt.EventsOn("init", (msg) => {
            rt.LogPrint("INIT msg!" + msg);
            $message = msg;
            activity.writeLog(msg);
        });
    });

    afterUpdate(async () => {
        const ou = document.querySelector(".lof");
        if (ou && ou.innerHTML.length > 980) {
            ou.scrollBy(0, 100);
        }
    });
</script>

<div class="main">
    <div class="userinput">
        <Device />
        <h1 id="⚙️">Settings</h1>
        <Settings />
    </div>
    <div class="output">
        <h1 id="🌝">
            <span>ACTIVITY LOG</span>
            <button on:click={activity.reset}>CLEAN</button>
        </h1>

        <p class="lof">{@html loo}</p>
    </div>
</div>

<style lang="scss">
    .main {
        display: flex;
        padding: 1.5rem;
        justify-content: space-between;
    }

    .userinput {
        max-width: 45%;
    }
    .output {
        display: inline-flex;
        flex-direction: column;
        justify-items: center;
        border: 1px outset #8400ff;
        border-radius: 5px;
        width: 50%;
        max-height: 760px;
    }
    .lof {
        padding-left: 5px;
        padding-right: 5px;
        font-size: 10px;
        text-align: left;
        opacity: 0.9;
        background-color: transparent;
        text-wrap: nowrap;
        overflow-y: auto;
    }
</style>
