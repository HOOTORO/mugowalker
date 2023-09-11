<script>
    import * as rt from "./lib/wailsjs/runtime/runtime.js";

    function quit() {
        rt.Quit();
    }
</script>

<nav>
    <a href="/">📲</a>
    <a href="/account">🤖</a>
    <a href="/settings">⚙️</a>
    <a href="/imaginer">🖼️</a>
</nav>

<div class="control" id="winbtn">
    <button class="btn" on:click={quit}>❌</button>
</div>

<style lang="scss">
    /* //////////////// */
    /* // Nav layout / */
    /* ////////////// */

    /* Horizontal */
    nav {
        display: flex;
        flex-direction: row;
        padding: 0;
        margin: 0;
        justify-content: center;
    }

    nav a {
        display: block;
        font-size: 2rem;
        padding: 2px;
        // margin: 1px;
        text-decoration: none;
        margin-top: -7px;
    }

    nav a:hover:not(.active),
    nav a.active {
        background-color: #3aeba1;
        background-image: linear-gradient(
            210deg,
            #ff0387 0%,
            #ff4949 49%,
            #ff8e04 100%
        );
        color: white;
        border-bottom-left-radius: 5px;
        border-bottom-right-radius: 5px;
    }

    nav a:hover {
        animation: slidebg 3s linear infinite;
    }
</style>
