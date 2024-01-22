export namespace main {
	
	export class Dalledress {
	    ens: string;
	    addr: string;
	    seed: string;
	    adverb: string;
	    adverbNum: number;
	    adjective: string;
	    adjectiveNum: number;
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
	
	    static createFrom(source: any = {}) {
	        return new Dalledress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ens = source["ens"];
	        this.addr = source["addr"];
	        this.seed = source["seed"];
	        this.adverb = source["adverb"];
	        this.adverbNum = source["adverbNum"];
	        this.adjective = source["adjective"];
	        this.adjectiveNum = source["adjectiveNum"];
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
	    }
	}

}

