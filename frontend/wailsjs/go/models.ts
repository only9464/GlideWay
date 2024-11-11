export namespace gitdorker {
	
	export class GithubResult {
	    Status: boolean;
	    Total: number;
	    Items: string[];
	    Link: string;
	
	    static createFrom(source: any = {}) {
	        return new GithubResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Status = source["Status"];
	        this.Total = source["Total"];
	        this.Items = source["Items"];
	        this.Link = source["Link"];
	    }
	}

}

export namespace portsscanner {
	
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

