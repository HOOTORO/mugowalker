<script lang="ts">
    import { settings } from "../stores/settings.js";
</script>

<div class="compo">
    <input bind:value={$settings.str} />
</div>
