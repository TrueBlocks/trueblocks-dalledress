export namespace app {
	
	export class GenerationProgress {
	    active: boolean;
	    series: string;
	    seed: string;
	    phase: string;
	    percent: number;
	    etaSeconds: number;
	    phasePercent: number;
	    phaseETASeconds: number;
	    phaseIndex: number;
	    phaseCount: number;
	    done: boolean;
	    cacheHit: boolean;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new GenerationProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.active = source["active"];
	        this.series = source["series"];
	        this.seed = source["seed"];
	        this.phase = source["phase"];
	        this.percent = source["percent"];
	        this.etaSeconds = source["etaSeconds"];
	        this.phasePercent = source["phasePercent"];
	        this.phaseETASeconds = source["phaseETASeconds"];
	        this.phaseIndex = source["phaseIndex"];
	        this.phaseCount = source["phaseCount"];
	        this.done = source["done"];
	        this.cacheHit = source["cacheHit"];
	        this.error = source["error"];
	    }
	}
	export class RuntimeInfo {
	    dataDir: string;
	    databaseVersion: string;
	    archiveHash: string;
	
	    static createFrom(source: any = {}) {
	        return new RuntimeInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dataDir = source["dataDir"];
	        this.databaseVersion = source["databaseVersion"];
	        this.archiveHash = source["archiveHash"];
	    }
	}

}

export namespace dalle {
	
	export class ArtifactSet {
	    generated?: string;
	    annotated?: string;
	
	    static createFrom(source: any = {}) {
	        return new ArtifactSet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.generated = source["generated"];
	        this.annotated = source["annotated"];
	    }
	}
	export class DatabaseRecordsResult {
	    name: string;
	    version: string;
	    columns?: string[];
	    records: storage.DatabaseRecord[];
	
	    static createFrom(source: any = {}) {
	        return new DatabaseRecordsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.columns = source["columns"];
	        this.records = this.convertValues(source["records"], storage.DatabaseRecord);
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
	export class ExportImageOptions {
	    Dir: string;
	    IncludePrompt: boolean;
	    IncludeData: boolean;
	    IncludeTitle: boolean;
	    IncludeTerse: boolean;
	    IncludeEnhanced: boolean;
	    IncludeTechnical: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ExportImageOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Dir = source["Dir"];
	        this.IncludePrompt = source["IncludePrompt"];
	        this.IncludeData = source["IncludeData"];
	        this.IncludeTitle = source["IncludeTitle"];
	        this.IncludeTerse = source["IncludeTerse"];
	        this.IncludeEnhanced = source["IncludeEnhanced"];
	        this.IncludeTechnical = source["IncludeTechnical"];
	    }
	}
	export class ExportImageResult {
	    dir: string;
	    files: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new ExportImageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dir = source["dir"];
	        this.files = source["files"];
	    }
	}
	export class GenerateRequest {
	    input: string;
	    seed?: string;
	    series?: string;
	    recipe?: string;
	    backstyle?: string;
	    enhance?: boolean;
	    image?: boolean;
	    annotate?: boolean;
	    force?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GenerateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input = source["input"];
	        this.seed = source["seed"];
	        this.series = source["series"];
	        this.recipe = source["recipe"];
	        this.backstyle = source["backstyle"];
	        this.enhance = source["enhance"];
	        this.image = source["image"];
	        this.annotate = source["annotate"];
	        this.force = source["force"];
	    }
	}
	export class MetadataStatus {
	    completed: boolean;
	    cacheHit: boolean;
	
	    static createFrom(source: any = {}) {
	        return new MetadataStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.completed = source["completed"];
	        this.cacheHit = source["cacheHit"];
	    }
	}
	export class StageStatus {
	    status: string;
	    cacheHit?: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new StageStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.cacheHit = source["cacheHit"];
	        this.error = source["error"];
	    }
	}
	export class PipelineStages {
	    selected: StageStatus;
	    prompted: StageStatus;
	    enhanced: StageStatus;
	    generated: StageStatus;
	    annotated: StageStatus;
	
	    static createFrom(source: any = {}) {
	        return new PipelineStages(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.selected = this.convertValues(source["selected"], StageStatus);
	        this.prompted = this.convertValues(source["prompted"], StageStatus);
	        this.enhanced = this.convertValues(source["enhanced"], StageStatus);
	        this.generated = this.convertValues(source["generated"], StageStatus);
	        this.annotated = this.convertValues(source["annotated"], StageStatus);
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
	export class PromptSet {
	    prompt?: string;
	    dataPrompt?: string;
	    titlePrompt?: string;
	    tersePrompt?: string;
	    enhancedPrompt?: string;
	
	    static createFrom(source: any = {}) {
	        return new PromptSet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.prompt = source["prompt"];
	        this.dataPrompt = source["dataPrompt"];
	        this.titlePrompt = source["titlePrompt"];
	        this.tersePrompt = source["tersePrompt"];
	        this.enhancedPrompt = source["enhancedPrompt"];
	    }
	}
	export class SelectedRecord {
	    attribute: string;
	    database: string;
	    rowIndex: number;
	    record: string;
	
	    static createFrom(source: any = {}) {
	        return new SelectedRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.attribute = source["attribute"];
	        this.database = source["database"];
	        this.rowIndex = source["rowIndex"];
	        this.record = source["record"];
	    }
	}
	export class MetadataDatabase {
	    version: string;
	    archiveHash: string;
	
	    static createFrom(source: any = {}) {
	        return new MetadataDatabase(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.archiveHash = source["archiveHash"];
	    }
	}
	export class MetadataRecipe {
	    name: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new MetadataRecipe(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	    }
	}
	export class MetadataSeries {
	    name: string;
	    hash: string;
	    source: string;
	
	    static createFrom(source: any = {}) {
	        return new MetadataSeries(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.hash = source["hash"];
	        this.source = source["source"];
	    }
	}
	export class ImageMetadata {
	    metadataVersion: string;
	    imageId: string;
	    input: string;
	    seed: string;
	    series: MetadataSeries;
	    recipe: MetadataRecipe;
	    database: MetadataDatabase;
	    selectedRecords: SelectedRecord[];
	    prompts: PromptSet;
	    artifacts: ArtifactSet;
	    stages: PipelineStages;
	    status: MetadataStatus;
	
	    static createFrom(source: any = {}) {
	        return new ImageMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.metadataVersion = source["metadataVersion"];
	        this.imageId = source["imageId"];
	        this.input = source["input"];
	        this.seed = source["seed"];
	        this.series = this.convertValues(source["series"], MetadataSeries);
	        this.recipe = this.convertValues(source["recipe"], MetadataRecipe);
	        this.database = this.convertValues(source["database"], MetadataDatabase);
	        this.selectedRecords = this.convertValues(source["selectedRecords"], SelectedRecord);
	        this.prompts = this.convertValues(source["prompts"], PromptSet);
	        this.artifacts = this.convertValues(source["artifacts"], ArtifactSet);
	        this.stages = this.convertValues(source["stages"], PipelineStages);
	        this.status = this.convertValues(source["status"], MetadataStatus);
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
	export class GenerateResult {
	    input: string;
	    seed: string;
	    series: string;
	    recipe: string;
	    metadataPath?: string;
	    generatedPath?: string;
	    annotatedPath?: string;
	    metadata: ImageMetadata;
	
	    static createFrom(source: any = {}) {
	        return new GenerateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input = source["input"];
	        this.seed = source["seed"];
	        this.series = source["series"];
	        this.recipe = source["recipe"];
	        this.metadataPath = source["metadataPath"];
	        this.generatedPath = source["generatedPath"];
	        this.annotatedPath = source["annotatedPath"];
	        this.metadata = this.convertValues(source["metadata"], ImageMetadata);
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
	
	export class ImageMetadataRecord {
	    path: string;
	    metadata: ImageMetadata;
	    archived: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ImageMetadataRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.metadata = this.convertValues(source["metadata"], ImageMetadata);
	        this.archived = source["archived"];
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
	    purpose?: string;
	    deleted?: boolean;
	    adverbs: string[];
	    adjectives: string[];
	    nouns: string[];
	    emotions: string[];
	    occupations: string[];
	    actions: string[];
	    artstyles: string[];
	    litstyles: string[];
	    colors: string[];
	    viewpoints: string[];
	    gazes: string[];
	    backstyles: string[];
	    compositions: string[];
	    colorLimit?: string;
	    modifiedAt?: string;
	    version?: string;
	    source?: string;
	
	    static createFrom(source: any = {}) {
	        return new Series(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.last = source["last"];
	        this.suffix = source["suffix"];
	        this.purpose = source["purpose"];
	        this.deleted = source["deleted"];
	        this.adverbs = source["adverbs"];
	        this.adjectives = source["adjectives"];
	        this.nouns = source["nouns"];
	        this.emotions = source["emotions"];
	        this.occupations = source["occupations"];
	        this.actions = source["actions"];
	        this.artstyles = source["artstyles"];
	        this.litstyles = source["litstyles"];
	        this.colors = source["colors"];
	        this.viewpoints = source["viewpoints"];
	        this.gazes = source["gazes"];
	        this.backstyles = source["backstyles"];
	        this.compositions = source["compositions"];
	        this.colorLimit = source["colorLimit"];
	        this.modifiedAt = source["modifiedAt"];
	        this.version = source["version"];
	        this.source = source["source"];
	    }
	}

}

export namespace storage {
	
	export class DatabaseFileManifest {
	    name: string;
	    path: string;
	    hash: string;
	    rows: number;
	    columns: string[];
	
	    static createFrom(source: any = {}) {
	        return new DatabaseFileManifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.hash = source["hash"];
	        this.rows = source["rows"];
	        this.columns = source["columns"];
	    }
	}
	export class DatabaseArchiveManifest {
	    version: string;
	    archiveHash: string;
	    files: DatabaseFileManifest[];
	
	    static createFrom(source: any = {}) {
	        return new DatabaseArchiveManifest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.archiveHash = source["archiveHash"];
	        this.files = this.convertValues(source["files"], DatabaseFileManifest);
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
	
	export class DatabaseRecord {
	    key: string;
	    values: string[];
	
	    static createFrom(source: any = {}) {
	        return new DatabaseRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.values = source["values"];
	    }
	}

}

