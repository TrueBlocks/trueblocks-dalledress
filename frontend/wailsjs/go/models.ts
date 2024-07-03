export namespace config {
	
	export class Session {
	    x: number;
	    y: number;
	    width: number;
	    height: number;
	    title: string;
	
	    static createFrom(source: any = {}) {
	        return new Session(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.title = source["title"];
	    }
	}

}

