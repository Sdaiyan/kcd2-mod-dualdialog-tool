export namespace main {
	
	export class CategoryConfig {
	    id: string;
	    sourceFile: string;
	    outputFile: string;
	    separator: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CategoryConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.sourceFile = source["sourceFile"];
	        this.outputFile = source["outputFile"];
	        this.separator = source["separator"];
	        this.enabled = source["enabled"];
	    }
	}

}

