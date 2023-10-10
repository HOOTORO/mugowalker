import { writable } from 'svelte/store';

let defaultDevice = {
    dev: "127.0.0.1:5555",
    account: "simpoleman",
    game: "Disconnected",
    online: false
}

export const device = writable(defaultDevice);