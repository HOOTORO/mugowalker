import { writable } from 'svelte/store';

let defaultDevice = { device: "127.0.0.1:5555" }

export const device = writable(defaultDevice);