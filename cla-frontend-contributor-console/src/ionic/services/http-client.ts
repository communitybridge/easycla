// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Injectable} from '@angular/core';
import {Headers, Http} from '@angular/http';
import {KeycloakService} from './keycloak/keycloak.service';
import {Observable} from 'rxjs/Rx';
import {AuthService} from './auth.service';

@Injectable()
export class HttpClient {

  constructor(
    public http: Http,
    private keycloak: KeycloakService,
    private authService: AuthService) {
  }

  private buildAuthHeaders(contentType: string = 'application/json') {
    let headers = new Headers({
      'Accept': 'application/json',
      'Content-Type': contentType,
    });

    if (this.authService.isAuthenticated()) {
      return this.authService.getIdToken().then(
        (token) => {
          if (token) {
            headers.append('Authorization', 'Bearer ' + token);
            return headers;
          }
        }
      );
    } else {
      return Promise.resolve(headers);
    }

  }

  setHttp(http: Http) {
    this.http = http; // allow alternate http library
  }

  buildHeaders(contentType: string = 'application/json') {
    let headers = new Headers({
      'Accept': 'application/json',
      'Content-Type': contentType
    });
    return Promise.resolve(headers);
  }


  get(url) {
    return Observable
      .fromPromise(this.buildAuthHeaders())
      .switchMap((headers) => this.http.get(url, {headers: headers}));
  }

  getWithCreds(url) {
    return Observable
      .fromPromise(this.buildAuthHeaders())
      .switchMap((headers) => this.http.get(url, {headers: headers, withCredentials: true}));
  }

  post(url, data) {
    return Observable
      .fromPromise(this.buildAuthHeaders())
      .switchMap((headers) => this.http.post(url, data, {headers: headers}));
  }

  postWithCreds(url, data) {
    return Observable
      .fromPromise(this.buildAuthHeaders())
      .switchMap((headers) => this.http.post(url, data, {headers: headers, withCredentials: true}));
  }

  put(url, data, contentType: string = 'application/json') {
    return Observable
      .fromPromise(this.buildAuthHeaders(contentType))
      .switchMap((headers) => this.http.put(url, data, {headers: headers}));
  }

  patch(url, data, contentType: string = 'application/json') {
    return Observable
      .fromPromise(this.buildAuthHeaders(contentType))
      .switchMap((headers) => this.http.patch(url, data, {headers: headers}));
  }

  delete(url) {
    return Observable
      .fromPromise(this.buildAuthHeaders())
      .switchMap((headers) => this.http.delete(url, {headers: headers}));
  }

  deleteWithBody(url, body) {
    return Observable
      .fromPromise(this.buildAuthHeaders())
      .switchMap((headers) => this.http.delete(url, {body: body, headers: headers}));
  }

  deleteWithCredsAndBody(url, body) {
    return Observable
      .fromPromise(this.buildAuthHeaders())
      .switchMap((headers) => this.http.delete(url, {body: body, headers: headers, withCredentials: true}));
  }
  securedGet(url) {
    return Observable
      .fromPromise(this.buildAuthHeaders())
      .switchMap((headers) => this.http.get(url, {headers: headers}));
  }

  securedPost(url, data) {
    return Observable
      .fromPromise(this.buildAuthHeaders())
      .switchMap((headers) => this.http.post(url, data, {headers: headers}));
  }

  securedPut(url, data, contentType: string = 'application/json') {
    return Observable
      .fromPromise(this.buildAuthHeaders(contentType))
      .switchMap((headers) => this.http.put(url, data, {headers: headers}));
  }

  securedPatch(url, data, contentType: string = 'application/json') {
    return Observable
      .fromPromise(this.buildAuthHeaders(contentType))
      .switchMap((headers) => this.http.patch(url, data, {headers: headers}));
  }

  securedDelete(url) {
    return Observable
      .fromPromise(this.buildAuthHeaders())
      .switchMap((headers) => this.http.delete(url, {headers: headers}));
  }

}
