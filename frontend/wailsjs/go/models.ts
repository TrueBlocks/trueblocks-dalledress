export namespace app {
	
	export class UserInfoStatus {
	    missingNameEmail: boolean;
	    rpcUnavailable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UserInfoStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.missingNameEmail = source["missingNameEmail"];
	        this.rpcUnavailable = source["rpcUnavailable"];
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

export namespace keys {
	
	export class Accelerator {
	    Key: string;
	    Modifiers: string[];
	
	    static createFrom(source: any = {}) {
	        return new Accelerator(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Key = source["Key"];
	        this.Modifiers = source["Modifiers"];
	    }
	}

}

export namespace menu {
	
	export class Menu {
	    Items: MenuItem[];
	
	    static createFrom(source: any = {}) {
	        return new Menu(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Items = this.convertValues(source["Items"], MenuItem);
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
	export class MenuItem {
	    Label: string;
	    Role: number;
	    Accelerator?: keys.Accelerator;
	    Type: string;
	    Disabled: boolean;
	    Hidden: boolean;
	    Checked: boolean;
	    SubMenu?: Menu;
	
	    static createFrom(source: any = {}) {
	        return new MenuItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Label = source["Label"];
	        this.Role = source["Role"];
	        this.Accelerator = this.convertValues(source["Accelerator"], keys.Accelerator);
	        this.Type = source["Type"];
	        this.Disabled = source["Disabled"];
	        this.Hidden = source["Hidden"];
	        this.Checked = source["Checked"];
	        this.SubMenu = this.convertValues(source["SubMenu"], Menu);
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
	export class CallbackData {
	    MenuItem?: MenuItem;
	
	    static createFrom(source: any = {}) {
	        return new CallbackData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MenuItem = this.convertValues(source["MenuItem"], MenuItem);
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

export namespace msgs {
	
	export enum EventType {
	    STATUS = "statusbar:log",
	    ERROR = "error:message",
	    MANAGER = "manager:change",
	    PROJECTS_UPDATED = "projects:updated",
	    APP_INIT = "app:initialized",
	    APP_READY = "app:ready",
	    VIEW_CHANGE = "app:view-changed",
	    VERSION = "1.0",
	    TAB_CYCLE = "hotkey:tab-cycle",
	    IMAGES_CHANGED = "images:changed",
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

export namespace preferences {
	
	export class Bounds {
	    x: number;
	    y: number;
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new Bounds(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = source["x"];
	        this.y = source["y"];
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}
	export class AppPreferences {
	    version?: string;
	    name?: string;
	    bounds?: Bounds;
	    recentProjects?: string[];
	    lastView?: string;
	    lastTab?: Record<string, string>;
	    lastViewNoWizard?: string;
	    menuCollapsed?: boolean;
	    helpCollapsed?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppPreferences(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.name = source["name"];
	        this.bounds = this.convertValues(source["bounds"], Bounds);
	        this.recentProjects = source["recentProjects"];
	        this.lastView = source["lastView"];
	        this.lastTab = source["lastTab"];
	        this.lastViewNoWizard = source["lastViewNoWizard"];
	        this.menuCollapsed = source["menuCollapsed"];
	        this.helpCollapsed = source["helpCollapsed"];
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
	
	export class Chain {
	    chain: string;
	    chainId: number;
	    remoteExplorer: string;
	    rpcProviders: string[];
	    symbol: string;
	
	    static createFrom(source: any = {}) {
	        return new Chain(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.chain = source["chain"];
	        this.chainId = source["chainId"];
	        this.remoteExplorer = source["remoteExplorer"];
	        this.rpcProviders = source["rpcProviders"];
	        this.symbol = source["symbol"];
	    }
	}
	export class Id {
	    appName: string;
	    orgName: string;
	    github: string;
	    domain: string;
	    twitter: string;
	
	    static createFrom(source: any = {}) {
	        return new Id(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appName = source["appName"];
	        this.orgName = source["orgName"];
	        this.github = source["github"];
	        this.domain = source["domain"];
	        this.twitter = source["twitter"];
	    }
	}
	export class OrgPreferences {
	    version?: string;
	    telemetry?: boolean;
	    theme?: string;
	    language?: string;
	    developerName?: string;
	    logLevel?: string;
	    experimental?: boolean;
	    supportUrl?: string;
	
	    static createFrom(source: any = {}) {
	        return new OrgPreferences(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.telemetry = source["telemetry"];
	        this.theme = source["theme"];
	        this.language = source["language"];
	        this.developerName = source["developerName"];
	        this.logLevel = source["logLevel"];
	        this.experimental = source["experimental"];
	        this.supportUrl = source["supportUrl"];
	    }
	}
	export class UserPreferences {
	    version?: string;
	    theme?: string;
	    language?: string;
	    name?: string;
	    email?: string;
	    chains?: Chain[];
	
	    static createFrom(source: any = {}) {
	        return new UserPreferences(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.theme = source["theme"];
	        this.language = source["language"];
	        this.name = source["name"];
	        this.email = source["email"];
	        this.chains = this.convertValues(source["chains"], Chain);
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

export namespace project {
	
	export class Project {
	    version: string;
	    name: string;
	    last_opened: string;
	    preferences: Record<string, string>;
	    dirty: boolean;
	    data: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.name = source["name"];
	        this.last_opened = source["last_opened"];
	        this.preferences = source["preferences"];
	        this.dirty = source["dirty"];
	        this.data = source["data"];
	    }
	}

}

export namespace utils {
	
	export class Explorer {
	    name: string;
	    url: string;
	    standard: string;
	
	    static createFrom(source: any = {}) {
	        return new Explorer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.url = source["url"];
	        this.standard = source["standard"];
	    }
	}
	export class NativeCurrency {
	    name: string;
	    symbol: string;
	    decimals: number;
	
	    static createFrom(source: any = {}) {
	        return new NativeCurrency(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.symbol = source["symbol"];
	        this.decimals = source["decimals"];
	    }
	}
	export class ChainListItem {
	    name: string;
	    chain: string;
	    icon: string;
	    rpc: string[];
	    faucets: string[];
	    nativeCurrency: NativeCurrency;
	    infoURL: string;
	    shortName: string;
	    chainId: number;
	    networkId: number;
	    explorers: Explorer[];
	
	    static createFrom(source: any = {}) {
	        return new ChainListItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.chain = source["chain"];
	        this.icon = source["icon"];
	        this.rpc = source["rpc"];
	        this.faucets = source["faucets"];
	        this.nativeCurrency = this.convertValues(source["nativeCurrency"], NativeCurrency);
	        this.infoURL = source["infoURL"];
	        this.shortName = source["shortName"];
	        this.chainId = source["chainId"];
	        this.networkId = source["networkId"];
	        this.explorers = this.convertValues(source["explorers"], Explorer);
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
	export class ChainList {
	    chains: ChainListItem[];
	    ChainsMap: Record<number, ChainListItem>;
	
	    static createFrom(source: any = {}) {
	        return new ChainList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.chains = this.convertValues(source["chains"], ChainListItem);
	        this.ChainsMap = this.convertValues(source["ChainsMap"], ChainListItem, true);
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

