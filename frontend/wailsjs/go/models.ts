export namespace app {
	
	export class NameEx {
	    address: base.Address;
	    decimals: number;
	    deleted?: boolean;
	    isContract?: boolean;
	    isCustom?: boolean;
	    isErc20?: boolean;
	    isErc721?: boolean;
	    isPrefund?: boolean;
	    name: string;
	    source: string;
	    symbol: string;
	    tags: string;
	    // Go type: base
	    prefund?: any;
	    type: names.Parts;
	
	    static createFrom(source: any = {}) {
	        return new NameEx(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = this.convertValues(source["address"], base.Address);
	        this.decimals = source["decimals"];
	        this.deleted = source["deleted"];
	        this.isContract = source["isContract"];
	        this.isCustom = source["isCustom"];
	        this.isErc20 = source["isErc20"];
	        this.isErc721 = source["isErc721"];
	        this.isPrefund = source["isPrefund"];
	        this.name = source["name"];
	        this.source = source["source"];
	        this.symbol = source["symbol"];
	        this.tags = source["tags"];
	        this.prefund = this.convertValues(source["prefund"], null);
	        this.type = source["type"];
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
	export class TransactionEx {
	    blockNumber: number;
	    transactionIndex: number;
	    timestamp: number;
	    date: string;
	    from: base.Address;
	    fromName: string;
	    to: base.Address;
	    toName: string;
	    // Go type: base
	    wei: any;
	    ether: string;
	    function: string;
	    hasToken: boolean;
	    isError: boolean;
	    logCount: number;
	
	    static createFrom(source: any = {}) {
	        return new TransactionEx(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.blockNumber = source["blockNumber"];
	        this.transactionIndex = source["transactionIndex"];
	        this.timestamp = source["timestamp"];
	        this.date = source["date"];
	        this.from = this.convertValues(source["from"], base.Address);
	        this.fromName = source["fromName"];
	        this.to = this.convertValues(source["to"], base.Address);
	        this.toName = source["toName"];
	        this.wei = this.convertValues(source["wei"], null);
	        this.ether = source["ether"];
	        this.function = source["function"];
	        this.hasToken = source["hasToken"];
	        this.isError = source["isError"];
	        this.logCount = source["logCount"];
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

export namespace base {
	
	export class Address {
	    address: number[];
	
	    static createFrom(source: any = {}) {
	        return new Address(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	    }
	}

}

export namespace config {
	
	export class Session {
	    x: number;
	    y: number;
	    width: number;
	    height: number;
	    title: string;
	    lastRoute: string;
	    lastTab: string;
	    lastAddress: string;
	    lastSeries: string;
	
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
	        this.lastRoute = source["lastRoute"];
	        this.lastTab = source["lastTab"];
	        this.lastAddress = source["lastAddress"];
	        this.lastSeries = source["lastSeries"];
	    }
	}

}

export namespace dalle {
	
	export class Attribute {
	    database: string;
	    name: string;
	    bytes: string;
	    number: number;
	    factor: number;
	    count: number;
	    selector: number;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new Attribute(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.database = source["database"];
	        this.name = source["name"];
	        this.bytes = source["bytes"];
	        this.number = source["number"];
	        this.factor = source["factor"];
	        this.count = source["count"];
	        this.selector = source["selector"];
	        this.value = source["value"];
	    }
	}
	export class DalleDress {
	    original: string;
	    fileName: string;
	    seed: string;
	    prompt?: string;
	    dataPrompt?: string;
	    titlePrompt?: string;
	    tersePrompt?: string;
	    enhancedPrompt?: string;
	    attributes: Attribute[];
	
	    static createFrom(source: any = {}) {
	        return new DalleDress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.original = source["original"];
	        this.fileName = source["fileName"];
	        this.seed = source["seed"];
	        this.prompt = source["prompt"];
	        this.dataPrompt = source["dataPrompt"];
	        this.titlePrompt = source["titlePrompt"];
	        this.tersePrompt = source["tersePrompt"];
	        this.enhancedPrompt = source["enhancedPrompt"];
	        this.attributes = this.convertValues(source["attributes"], Attribute);
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
	export class Series {
	    last?: number;
	    suffix: string;
	    adverbs: string[];
	    adjectives: string[];
	    nouns: string[];
	    emotions: string[];
	    occupations: string[];
	    actions: string[];
	    artstyles: string[];
	    litstyles: string[];
	    colors: string[];
	    orientations: string[];
	    gazes: string[];
	    backstyles: string[];
	
	    static createFrom(source: any = {}) {
	        return new Series(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.last = source["last"];
	        this.suffix = source["suffix"];
	        this.adverbs = source["adverbs"];
	        this.adjectives = source["adjectives"];
	        this.nouns = source["nouns"];
	        this.emotions = source["emotions"];
	        this.occupations = source["occupations"];
	        this.actions = source["actions"];
	        this.artstyles = source["artstyles"];
	        this.litstyles = source["litstyles"];
	        this.colors = source["colors"];
	        this.orientations = source["orientations"];
	        this.gazes = source["gazes"];
	        this.backstyles = source["backstyles"];
	    }
	}

}

export namespace names {
	
	export enum Parts {
	    REGULAR = 2,
	    CUSTOM = 4,
	    PREFUND = 8,
	    BADDRESS = 16,
	}

}

export namespace output {
	
	export class RenderCtx {
	
	
	    static createFrom(source: any = {}) {
	        return new RenderCtx(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace servers {
	
	export enum Type {
	    FILESERVER = 0,
	    SCRAPER = 1,
	    MONITOR = 2,
	    API = 3,
	    IPFS = 4,
	}
	export enum State {
	    STOPPED = 0,
	    RUNNING = 1,
	    PAUSED = 2,
	}

}

