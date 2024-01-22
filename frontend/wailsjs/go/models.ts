export namespace main {
	
	export class Dalledress {
	    ens: string;
	    addr: string;
	    seed: string;
	    adverb: string;
	    adverbSeed: string;
	    adverbNum: number;
	    adjective: string;
	    adjectiveSeed: string;
	    adjectiveNum: number;
	    noun: string;
	    nounSeed: string;
	    nounNum: number;
	    style: string;
	    styleSeed: string;
	    styleNum: number;
	    color1: string;
	    color1Seed: string;
	    color1Num: number;
	    color2: string;
	    color2Seed: string;
	    color2Num: number;
	    color3: string;
	    color3Seed: string;
	    color3Num: number;
	    variant1: string;
	    variant1Seed: string;
	    variant1Num: number;
	    variant2: string;
	    variant2Seed: string;
	    variant2Num: number;
	    variant3: string;
	    variant3Seed: string;
	    variant3Num: number;
	    style2: string;
	    style2Seed: string;
	    style2Num: number;
	    background: string;
	    backgroundSeed: string;
	    backgroundNum: number;
	
	    static createFrom(source: any = {}) {
	        return new Dalledress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ens = source["ens"];
	        this.addr = source["addr"];
	        this.seed = source["seed"];
	        this.adverb = source["adverb"];
	        this.adverbSeed = source["adverbSeed"];
	        this.adverbNum = source["adverbNum"];
	        this.adjective = source["adjective"];
	        this.adjectiveSeed = source["adjectiveSeed"];
	        this.adjectiveNum = source["adjectiveNum"];
	        this.noun = source["noun"];
	        this.nounSeed = source["nounSeed"];
	        this.nounNum = source["nounNum"];
	        this.style = source["style"];
	        this.styleSeed = source["styleSeed"];
	        this.styleNum = source["styleNum"];
	        this.color1 = source["color1"];
	        this.color1Seed = source["color1Seed"];
	        this.color1Num = source["color1Num"];
	        this.color2 = source["color2"];
	        this.color2Seed = source["color2Seed"];
	        this.color2Num = source["color2Num"];
	        this.color3 = source["color3"];
	        this.color3Seed = source["color3Seed"];
	        this.color3Num = source["color3Num"];
	        this.variant1 = source["variant1"];
	        this.variant1Seed = source["variant1Seed"];
	        this.variant1Num = source["variant1Num"];
	        this.variant2 = source["variant2"];
	        this.variant2Seed = source["variant2Seed"];
	        this.variant2Num = source["variant2Num"];
	        this.variant3 = source["variant3"];
	        this.variant3Seed = source["variant3Seed"];
	        this.variant3Num = source["variant3Num"];
	        this.style2 = source["style2"];
	        this.style2Seed = source["style2Seed"];
	        this.style2Num = source["style2Num"];
	        this.background = source["background"];
	        this.backgroundSeed = source["backgroundSeed"];
	        this.backgroundNum = source["backgroundNum"];
	    }
	}

}

