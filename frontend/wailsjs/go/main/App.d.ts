// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function GetDirsearchProgress():Promise<main.DirsearchProgress>;

export function GetDirsearchStatus():Promise<string>;

export function GetScanProgress():Promise<main.ScanProgress>;

export function GetScanStatus():Promise<string>;

export function OpenFileDialog():Promise<string>;

export function ScanPorts(arg1:string,arg2:number,arg3:number,arg4:number):Promise<void>;

export function StartDirsearch(arg1:string,arg2:string,arg3:number):Promise<void>;

export function StopDirsearch():Promise<void>;

export function StopScan():Promise<void>;
