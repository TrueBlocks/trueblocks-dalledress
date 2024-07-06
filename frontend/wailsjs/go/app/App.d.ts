// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {config} from '../models';
import {dalle} from '../models';

export function GenerateEnhanced(arg1:string):Promise<string>;

export function GenerateImage(arg1:string):Promise<string>;

export function GetData(arg1:string):Promise<string>;

export function GetEnhanced(arg1:string):Promise<string>;

export function GetFileList():Promise<Array<string>>;

export function GetImage(arg1:string):Promise<string>;

export function GetJson(arg1:string):Promise<string>;

export function GetLastAddress():Promise<string>;

export function GetLastRoute():Promise<string>;

export function GetLastSeries():Promise<string>;

export function GetLastTab():Promise<string>;

export function GetNames(arg1:number,arg2:number):Promise<Array<string>>;

export function GetPrompt(arg1:string):Promise<string>;

export function GetSeries(arg1:string):Promise<string>;

export function GetSession():Promise<config.Session>;

export function GetTerse(arg1:string):Promise<string>;

export function GetTitle(arg1:string):Promise<string>;

export function HandleLines():Promise<void>;

export function LoadSeries():Promise<dalle.Series>;

export function MaxNames():Promise<number>;

export function Refresh(arg1:string):Promise<string>;

export function SetLastAddress(arg1:string):Promise<void>;

export function SetLastRoute(arg1:string):Promise<void>;

export function SetLastSeries(arg1:string):Promise<void>;

export function SetLastTab(arg1:string):Promise<void>;

export function String():Promise<string>;
