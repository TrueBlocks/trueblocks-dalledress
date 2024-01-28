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
	export class Dalledress {
	    ens: string;
	    addr: string;
	    seed: string;
	    adverb: Attribute;
	    adjective: string;
	    adjectiveNum: number;
	    emotion1: string;
	    emotion1Num: number;
	    emotion2: string;
	    emotion2Num: number;
	    literary: string;
	    literaryNum: number;
	    noun: string;
	    nounNum: number;
	    style: string;
	    styleNum: number;
	    color1: string;
	    color1Num: number;
	    color2: string;
	    color2Num: number;
	    color3: string;
	    color3Num: number;
	    variant1: string;
	    variant1Num: number;
	    variant2: string;
	    variant2Num: number;
	    variant3: string;
	    variant3Num: number;
	    style2: string;
	    style2Num: number;
	    background: string;
	    backgroundNum: number;
	    orientation: string;
	    orientationNum: number;
	    prompt: string;
	    promptError: any;
	    data: string;
	    dataError: any;
	    terse: string;
	    terseError: any;
	
	    static createFrom(source: any = {}) {
	        return new Dalledress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ens = source["ens"];
	        this.addr = source["addr"];
	        this.seed = source["seed"];
	        this.adverb = this.convertValues(source["adverb"], Attribute);
	        this.adjective = source["adjective"];
	        this.adjectiveNum = source["adjectiveNum"];
	        this.emotion1 = source["emotion1"];
	        this.emotion1Num = source["emotion1Num"];
	        this.emotion2 = source["emotion2"];
	        this.emotion2Num = source["emotion2Num"];
	        this.literary = source["literary"];
	        this.literaryNum = source["literaryNum"];
	        this.noun = source["noun"];
	        this.nounNum = source["nounNum"];
	        this.style = source["style"];
	        this.styleNum = source["styleNum"];
	        this.color1 = source["color1"];
	        this.color1Num = source["color1Num"];
	        this.color2 = source["color2"];
	        this.color2Num = source["color2Num"];
	        this.color3 = source["color3"];
	        this.color3Num = source["color3Num"];
	        this.variant1 = source["variant1"];
	        this.variant1Num = source["variant1Num"];
	        this.variant2 = source["variant2"];
	        this.variant2Num = source["variant2Num"];
	        this.variant3 = source["variant3"];
	        this.variant3Num = source["variant3Num"];
	        this.style2 = source["style2"];
	        this.style2Num = source["style2Num"];
	        this.background = source["background"];
	        this.backgroundNum = source["backgroundNum"];
	        this.orientation = source["orientation"];
	        this.orientationNum = source["orientationNum"];
	        this.prompt = source["prompt"];
	        this.promptError = source["promptError"];
	        this.data = source["data"];
	        this.dataError = source["dataError"];
	        this.terse = source["terse"];
	        this.terseError = source["terseError"];
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

