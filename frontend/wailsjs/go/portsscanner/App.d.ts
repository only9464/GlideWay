// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {portsscanner} from '../models';
import {context} from '../models';

export function GetScanProgress():Promise<portsscanner.ScanProgress>;

export function GetScanStatus():Promise<string>;

export function ScanPorts(arg1:string,arg2:number,arg3:number,arg4:number):Promise<void>;

export function Startup(arg1:context.Context):Promise<void>;

export function StopScan():Promise<void>;
