<script lang="ts">
    import { device } from "../stores/device.js";
</script>

<div id="device">
    <span>{$device}</span>
</div>

<style>
    #device {
        display: flex;
        flex-direction: column;
        margin: 0px;
        padding: 10px;
    }
</style>
