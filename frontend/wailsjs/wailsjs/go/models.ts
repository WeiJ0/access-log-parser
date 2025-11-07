export namespace app {
	
	export class ParseFileRequest {
	    filePath: string;
	
	    static createFrom(source: any = {}) {
	        return new ParseFileRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	    }
	}
	export class ParseFileResponse {
	    success: boolean;
	    logFile?: models.LogFile;
	    errorMessage: string;
	    errorSamples: parser.ParseError[];
	
	    static createFrom(source: any = {}) {
	        return new ParseFileResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.logFile = this.convertValues(source["logFile"], models.LogFile);
	        this.errorMessage = source["errorMessage"];
	        this.errorSamples = this.convertValues(source["errorSamples"], parser.ParseError);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class SelectFileResponse {
	    success: boolean;
	    filePath: string;
	    errorMessage: string;
	
	    static createFrom(source: any = {}) {
	        return new SelectFileResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.filePath = source["filePath"];
	        this.errorMessage = source["errorMessage"];
	    }
	}
	export class ValidateFormatRequest {
	    filePath: string;
	
	    static createFrom(source: any = {}) {
	        return new ValidateFormatRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	    }
	}
	export class ValidateFormatResponse {
	    success: boolean;
	    valid: boolean;
	    errorMessage: string;
	
	    static createFrom(source: any = {}) {
	        return new ValidateFormatResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.valid = source["valid"];
	        this.errorMessage = source["errorMessage"];
	    }
	}

}

export namespace models {
	
	export class LogEntry {
	    ip: string;
	    // Go type: time
	    timestamp: any;
	    method: string;
	    url: string;
	    protocol: string;
	    statusCode: number;
	    responseBytes: number;
	    referer: string;
	    userAgent: string;
	    user?: string;
	    requestTime?: number;
	    lineNumber: number;
	    rawLine: string;
	    parseError?: string;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.method = source["method"];
	        this.url = source["url"];
	        this.protocol = source["protocol"];
	        this.statusCode = source["statusCode"];
	        this.responseBytes = source["responseBytes"];
	        this.referer = source["referer"];
	        this.userAgent = source["userAgent"];
	        this.user = source["user"];
	        this.requestTime = source["requestTime"];
	        this.lineNumber = source["lineNumber"];
	        this.rawLine = source["rawLine"];
	        this.parseError = source["parseError"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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
	export class LogFile {
	    path: string;
	    name: string;
	    size: number;
	    // Go type: time
	    loadedAt: any;
	    totalLines: number;
	    parsedLines: number;
	    errorLines: number;
	    entries: LogEntry[];
	    statistics: any;
	    parseTime: number;
	    statTime: number;
	    memoryUsed: number;
	
	    static createFrom(source: any = {}) {
	        return new LogFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.loadedAt = this.convertValues(source["loadedAt"], null);
	        this.totalLines = source["totalLines"];
	        this.parsedLines = source["parsedLines"];
	        this.errorLines = source["errorLines"];
	        this.entries = this.convertValues(source["entries"], LogEntry);
	        this.statistics = source["statistics"];
	        this.parseTime = source["parseTime"];
	        this.statTime = source["statTime"];
	        this.memoryUsed = source["memoryUsed"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

export namespace parser {
	
	export class ParseError {
	    lineNumber: number;
	    line: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new ParseError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lineNumber = source["lineNumber"];
	        this.line = source["line"];
	        this.error = source["error"];
	    }
	}

}

