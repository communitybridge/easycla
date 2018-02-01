import { Injectable } from '@angular/core';
import { Http } from '@angular/http';

import 'rxjs/Rx';

@Injectable()
export class S3Service {
  http: any;

  constructor(public http: Http) {
    this.http = http;
  }

  public setHttp(http: any) {
    this.http = http; // allow configuration for alternate http library
  }

  uploadToS3(url, file, contentType) {
    return this.http.put(url, file, contentType)
      .map((res) => res);
  }

}
