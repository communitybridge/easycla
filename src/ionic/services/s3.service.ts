import { Injectable } from '@angular/core';
import { Http } from '@angular/http';

import 'rxjs/Rx';

@Injectable()
export class S3Service {
  http: any;

  constructor(public http: Http) {

  }

  public setHttp(http: Http) {
    this.http = http; // allow configuration for alternate http library
  }

  uploadToS3(url, file, contentType) {
    return this.http.put(url, file, contentType)
      .map((res) => res);
  }

}
