import { writable } from 'svelte/store';

let defaultSettings = { str: "(❁´◡`❁)" }

export const settings = writable(defaultSettings);