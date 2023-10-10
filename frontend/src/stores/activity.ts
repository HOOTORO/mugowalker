import { writable } from 'svelte/store';

function createLog() {
    const head = "<h3 style='color:orange'>( ͡°( ͡° ͜ʖ( ͡° ͜ʖ ͡°)ʖ ͡°) ͡°)     ᓚᘏᗢ </h3>"
    const { subscribe, set, update } = writable(head);

    return {
        subscribe,
        writeLog: (str: string) => update((n) => n + "<br>" + str),
        reset: () => set(head)
    };
}

export let activity = createLog();
