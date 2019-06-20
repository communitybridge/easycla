// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Injectable } from '@angular/core';
import { Http } from '@angular/http';

import 'rxjs/Rx';

@Injectable()
export class S3Service {

  constructor(public http: Http) {

  }

  public setHttp(http: any) {
    this.http = http; // allow configuration for alternate http library
  }

  uploadToS3(url, file, contentType) {
    return this.http.put(url, file, contentType)
      .map((res) => res);
  }

}
