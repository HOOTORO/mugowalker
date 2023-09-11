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

