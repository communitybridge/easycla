// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Observable } from 'rxjs/Rx';

@Injectable()
export class HttpClient {

  constructor(
    private http: Http,
  ) {

  }

  public setHttp(http: Http) {
    this.http = http; // allow alternate http library
  }

  private buildHeaders(contentType: string = 'application/json') {
    let headers = new Headers({
      'Accept': 'application/json',
      'Content-Type': contentType
    });
    return Promise.resolve(headers);
  }

  get(url) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) =>
          this.http.get(url, { headers: headers }));
  }

  post(url, data) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) => this.http.post(url, data, { headers: headers }));
  }

  put(url, data, contentType: string = 'application/json') {
    return Observable
        .fromPromise(this.buildHeaders(contentType))
        .switchMap((headers) => this.http.put(url, data, { headers: headers }));
  }

  patch(url, data, contentType: string = 'application/json') {
    return Observable
        .fromPromise(this.buildHeaders(contentType))
        .switchMap((headers) => this.http.patch(url, data, { headers: headers }));
  }

  delete(url) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) => this.http.delete(url, { headers: headers }));
  }

}
