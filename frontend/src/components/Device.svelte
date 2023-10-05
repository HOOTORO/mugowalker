<script lang="ts">
    import { afterUpdate } from "svelte";
    import { AdbConnect } from "../lib/wailsjs/go/backend/Config";
    import * as rt from "../lib/wailsjs/runtime/runtime.js";
    import { activity } from "../stores/activity";
    import { device } from "../stores/device";
    import Butt from "./Butt.svelte";

    const devStatus = {
        online: `<span style="color: green;text-align:right"><strong>CONNECTED</strong></span>`,
        offline: `<span style="color: red;text-align:right"><strong>DISCONNECTED</strong></span>`,
        con: `CONNECT`,
        dis: `DISCONNECT`,
        loading: `CONNECTING`,
    };
    let buttonName = devStatus.con;
    let btnActive = true;
    let isAvailable = devStatus.offline;

    // let result: Promise<string | boolean> | null = null;

    function connect() {
        buttonName = devStatus.loading;
        btnActive = false;
        activity.writeLog(`Connecting to ${$device.dest}`);
        rt.EventsEmit("task", $device.dest);
        AdbConnect($device.dest)
            .then((x) => {
                if (x) {
                    isAvailable = devStatus.online;
                    buttonName = devStatus.dis;
                    activity.writeLog("connection successful");
                } else {
                    btnActive = true;
                    buttonName = devStatus.con;
                    isAvailable = devStatus.offline;
                    activity.writeLog("connection failed");
                }
            })
            .catch((e) => {
                btnActive = true;
                buttonName = devStatus.con;
                isAvailable = devStatus.offline;
                activity.writeLog(devStatus.offline);
            });
    }
    afterUpdate(async () => {
        rt.EventsEmit("config", $device.dest);
    }); // console.log(result);
    // activity += result.
    function updCfg() {
        rt.EventsEmit("config", $device.dest);
    }
</script>

<div class="component">
    <h1 id="📲">Device <span><sup>{@html isAvailable}</sup></span></h1>
    <input bind:value={$device.dest} on:change={updCfg} />
    <Butt name={buttonName} on:click={connect} ro={btnActive} />
</div>

<style lang="scss">
    .component {
        display: block;
        padding: 1em;

        h1:before {
            content: attr(id);
            text-shadow: none;
        }
    }
</style>
