import { writable } from 'svelte/store';

let defaultSettings = {
    drawstep: false,
    logfile: "application.log",
    loglevel: "FATAL",
    imagick: {
        colorspace: "Gray",
        alpha: "off",
        threshold: "80%",
        edge: "2",
        negate: true,
        blackthreshold: "90%",
    },
    tesseract: {
        psm: 3,
        args: ["quiet"]
    },
    bluestacks: {
        instance: "Rvc64",
        cmd: "launchApp",
        package: "com.lilithgames.hgame.gp.id"
    },
    ignoredwords: ["Go", "Up", "In", "Tap"]
}

export const settings = writable(defaultSettings);