export namespace main {
	
	export class DirsearchProgress {
	    current: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new DirsearchProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current = source["current"];
	        this.total = source["total"];
	    }
	}
	export class ScanProgress {
	    current_port: number;
	    total_ports: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.current_port = source["current_port"];
	        this.total_ports = source["total_ports"];
	        this.status = source["status"];
	    }
	}

}

