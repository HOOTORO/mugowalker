<script lang="ts">
    import { device } from "../stores/device";
    import { AdbConnect } from "../lib/wailsjs/go/main/App";
    import Butt from "./Butt.svelte";
    import { activity } from "../stores/activity";

    const devStatus = {
        online: `<span style="color: green;text-align:right"><strong>CONNECTED</strong></span>`,
        offline: `<span style="color: red;text-align:right"><strong>DISCONNECTED</strong></span>`,
        con: `CONNECT`,
        dis: `DISCONNECT`,
        loading: `CONNECTING`,
    };
    let btname = devStatus.con;
    let btnActifve = true;
    let isAvailiable = devStatus.offline;

    // let result: Promise<string | boolean> | null = null;

    function connect() {
        btname = devStatus.loading;
        btnActifve = false;
        activity.writeLog(`Connecting to ${$device.dest}`);
        AdbConnect($device.dest)
            .then((x) => {
                console.log(x);
                if (x) {
                    isAvailiable = devStatus.online;
                    btname = devStatus.dis;
                    activity.writeLog("connection succesful");
                } else {
                    btnActifve = true;
                    btname = devStatus.con;
                    isAvailiable = devStatus.offline;
                    activity.writeLog("connection failed");
                }
            })
            .catch((e) => {
                btnActifve = true;
                btname = devStatus.con;
                isAvailiable = devStatus.offline;
                console.log(e);
                activity.writeLog(devStatus.offline);
            });
        // console.log(result);
        // activity += result.
    }
</script>

<div class="component">
    <h1 id="📲">Device <span><sup>{@html isAvailiable}</sup></span></h1>
    <input bind:value={$device.dest} />
    <Butt name={btname} on:click={connect} ro={btnActifve} />
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
