export namespace main {
	
	export class OptionsType {
	    caseInsensitive?: boolean;
	    wholeWord?: boolean;
	    wholeLine?: boolean;
	    filenameOnly?: boolean;
	    filesWoMatches?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OptionsType(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.caseInsensitive = source["caseInsensitive"];
	        this.wholeWord = source["wholeWord"];
	        this.wholeLine = source["wholeLine"];
	        this.filenameOnly = source["filenameOnly"];
	        this.filesWoMatches = source["filesWoMatches"];
	    }
	}

}

export namespace settings {
	
	export class Imagick {
	    colorspace: string;
	    alpha: string;
	    threshold: string;
	    edge: string;
	    negate: boolean;
	    "gaussian-blur": string;
	    "auto-threshold": string;
	
	    static createFrom(source: any = {}) {
	        return new Imagick(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.colorspace = source["colorspace"];
	        this.alpha = source["alpha"];
	        this.threshold = source["threshold"];
	        this.edge = source["edge"];
	        this.negate = source["negate"];
	        this["gaussian-blur"] = source["gaussian-blur"];
	        this["auto-threshold"] = source["auto-threshold"];
	    }
	}
	export class Pilot {
	    dev: string;
	    account: string;
	    game: string;
	    online: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Pilot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dev = source["dev"];
	        this.account = source["account"];
	        this.game = source["game"];
	        this.online = source["online"];
	    }
	}
	export class Tesseract {
	    psm: number;
	    args: string[];
	
	    static createFrom(source: any = {}) {
	        return new Tesseract(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.psm = source["psm"];
	        this.args = source["args"];
	    }
	}
	export class Settings {
	    drawstep: boolean;
	    logfile: string;
	    loglevel: string;
	    imagick: Imagick;
	    tesseract: Tesseract;
	    // Go type: struct { Instance string "json:\"instance\""; Cmd string "json:\"cmd\""; Package string "json:\"package\"" }
	    bluestacks: any;
	    ignoredwords: string[];
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.drawstep = source["drawstep"];
	        this.logfile = source["logfile"];
	        this.loglevel = source["loglevel"];
	        this.imagick = this.convertValues(source["imagick"], Imagick);
	        this.tesseract = this.convertValues(source["tesseract"], Tesseract);
	        this.bluestacks = this.convertValues(source["bluestacks"], Object);
	        this.ignoredwords = source["ignoredwords"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

