import { Injectable } from '@angular/core';

declare const webpackGlobalVars: any;

@Injectable()
export class ConfigService {
  CLA_API_URL:string = webpackGlobalVars.claApiUrl;
    // Add more configuration variables as needed
}
