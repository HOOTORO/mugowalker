import { writable } from 'svelte/store';

function createLog() {
    const { subscribe, set, update } = writable("( ͡°( ͡° ͜ʖ( ͡° ͜ʖ ͡°)ʖ ͡°) ͡°)       ᓚᘏᗢ<br>");

    return {
        subscribe,
        writeLog: (str: string) => update((n) => n + "<br>" + str),
        // decrement: () => update((n) => n - 1),
        reset: () => set("( ͡°( ͡° ͜ʖ( ͡° ͜ʖ ͡°)ʖ ͡°) ͡°)       ᓚᘏᗢ<br>")
    };
}

export let activity = createLog();
