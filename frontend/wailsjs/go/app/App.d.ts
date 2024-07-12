// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {types} from '../models';
import {config} from '../models';
import {dalle} from '../models';

export function GenerateEnhanced(arg1:string):Promise<string>;

export function GenerateImage(arg1:string):Promise<string>;

export function GetData(arg1:string):Promise<string>;

export function GetEnhanced(arg1:string):Promise<string>;

export function GetExistingAddrs():Promise<Array<string>>;

export function GetFilelist(arg1:string):Promise<Array<string>>;

export function GetFilename(arg1:string):Promise<string>;

export function GetJson(arg1:string):Promise<string>;

export function GetLast(arg1:string):Promise<string>;

export function GetNames(arg1:number,arg2:number):Promise<Array<types.Name>>;

export function GetNamesCnt():Promise<number>;

export function GetPrompt(arg1:string):Promise<string>;

export function GetSeries(arg1:string):Promise<string>;

export function GetSession():Promise<config.Session>;

export function GetTerse(arg1:string):Promise<string>;

export function GetTitle(arg1:string):Promise<string>;

export function HandleLines():Promise<void>;

export function LoadSeries():Promise<dalle.Series>;

export function MakeDalleDress(arg1:string):Promise<dalle.DalleDress>;

export function ReloadDatabases():Promise<void>;

export function Save(arg1:string):Promise<boolean>;

export function SetLast(arg1:string,arg2:string):Promise<void>;

export function String():Promise<string>;
