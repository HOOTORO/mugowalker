<script lang="ts">
    import type * as mods from "../lib/wailsjs/go/models";

    import { onMount } from "svelte";
    import { CurrentPilot } from "../lib/wailsjs/go/backend/Config.js";
    import * as rt from "../lib/wailsjs/runtime/runtime.js";
    import { activity } from "../stores/activity";
    import { device } from "../stores/device";
    import Butt from "./Butt.svelte";

    let d: mods.pilot.Pilot;
    device.subscribe((device) => {
        d = device;
    });
    const tasks = [
        "Daily",
        "Push Campain",
        "Kings Tower",
        "World Tree",
        "Forsaken Necropolis",
        "Towers of Light",
        "Brutal Citadel",
        "Celestial Sanctum",
        "Infernal Fortress",
    ];
    const devStatus = {
        online: `<span style="color: green;text-align:right"><strong>CONNECTED</strong></span>`,
        offline: `<span style="color: red;text-align:right"><strong>DISCONNECTED</strong></span>`,
        con: `CONNECT`,
        dis: `DISCONNECT`,
        loading: `CONNECTING`,
    };
    let buttonName = devStatus.con;
    let isAvailable = devStatus.offline;
    onMount(async () => {
        $device = await CurrentPilot();
        if ($device.online) {
            isAvailable = devStatus.online;
            buttonName = devStatus.dis;
        }

        rt.EventsOn("devstate", (msg) => {
            rt.LogPrint("Connection result -> " + msg);
            $device.game = msg;
            activity.writeLog(msg);
            if (msg === "success") {
                $device.online = true;
                buttonName = devStatus.dis;
                isAvailable = devStatus.online;
            } else {
                $device.online = false;
                buttonName = devStatus.con;
                isAvailable = devStatus.offline;
            }
        });
    });
    function connect() {
        buttonName = devStatus.loading;
        activity.writeLog(`Connecting to ${$device.dev}`);
        rt.EventsEmit("connection", $device.dev);
    }

    function runTask(task: string) {
        rt.EventsEmit("task", task);
    }
</script>

<div class="component">
    <h1>Device <span><sup>{@html isAvailable}</sup></span></h1>
    <input bind:value={$device.dev} />
    <Butt name={buttonName} on:click={connect} ro={!$device.online} />
    <h3>Pilot namaeva</h3>
    <input bind:value={$device.account} />

    <div>
        <h2>Tasks</h2>
        <div class="btn-container">
            {#each tasks as task}
                <button
                    on:click={() => runTask(task)}
                    disabled={!$device.online}>{task}</button
                >
            {/each}
        </div>
    </div>
</div>

<style lang="scss">
    .component {
        display: block;
        padding: 1em;
    }
    .btn-container {
        display: flex;
        flex-wrap: wrap;
        align-content: space-around;
        justify-content: center;
        align-items: stretch;
        flex-direction: row;

        button {
            border-radius: 10px;
            background-color: coral;
            margin: 4px;
        }
    }
</style>
