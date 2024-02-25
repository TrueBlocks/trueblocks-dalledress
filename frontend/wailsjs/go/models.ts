export namespace main {
	
	export class Attribute {
	    seed: string;
	    num: number;
	    val: string;
	
	    static createFrom(source: any = {}) {
	        return new Attribute(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.seed = source["seed"];
	        this.num = source["num"];
	        this.val = source["val"];
	    }
	}
	export class Value {
	    val: string;
	    error: any;
	
	    static createFrom(source: any = {}) {
	        return new Value(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.val = source["val"];
	        this.error = source["error"];
	    }
	}
	export class Dalledress {
	    ens: string;
	    addr: string;
	    seed: string;
	    adverb: Attribute;
	    adjective: Attribute;
	    emotionShort: Attribute;
	    emotion: Attribute;
	    literary: Attribute;
	    noun: Attribute;
	    style: Attribute;
	    shortStyle: string;
	    style2: Attribute;
	    color1: Attribute;
	    color2: Attribute;
	    color3: Attribute;
	    variant1: Attribute;
	    variant2: Attribute;
	    variant3: Attribute;
	    background: Attribute;
	    orientation: Attribute;
	    prompt: Value;
	    data: Value;
	    terse: Value;
	
	    static createFrom(source: any = {}) {
	        return new Dalledress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ens = source["ens"];
	        this.addr = source["addr"];
	        this.seed = source["seed"];
	        this.adverb = this.convertValues(source["adverb"], Attribute);
	        this.adjective = this.convertValues(source["adjective"], Attribute);
	        this.emotionShort = this.convertValues(source["emotionShort"], Attribute);
	        this.emotion = this.convertValues(source["emotion"], Attribute);
	        this.literary = this.convertValues(source["literary"], Attribute);
	        this.noun = this.convertValues(source["noun"], Attribute);
	        this.style = this.convertValues(source["style"], Attribute);
	        this.shortStyle = source["shortStyle"];
	        this.style2 = this.convertValues(source["style2"], Attribute);
	        this.color1 = this.convertValues(source["color1"], Attribute);
	        this.color2 = this.convertValues(source["color2"], Attribute);
	        this.color3 = this.convertValues(source["color3"], Attribute);
	        this.variant1 = this.convertValues(source["variant1"], Attribute);
	        this.variant2 = this.convertValues(source["variant2"], Attribute);
	        this.variant3 = this.convertValues(source["variant3"], Attribute);
	        this.background = this.convertValues(source["background"], Attribute);
	        this.orientation = this.convertValues(source["orientation"], Attribute);
	        this.prompt = this.convertValues(source["prompt"], Value);
	        this.data = this.convertValues(source["data"], Value);
	        this.terse = this.convertValues(source["terse"], Value);
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
	export class Settings {
	    title: string;
	    width: number;
	    height: number;
	    x: number;
	    y: number;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.x = source["x"];
	        this.y = source["y"];
	    }
	}

}

