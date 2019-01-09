import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { KeycloakService } from './keycloak/keycloak.service';
import { Observable } from 'rxjs/Rx';
import { AuthService } from './auth.service';

@Injectable()
export class HttpClient {

  constructor(
    public http: Http,
    private keycloak: KeycloakService,
    private authService: AuthService) {}

  buildAuthHeaders(contentType: string = 'application/json') {
    let headers = new Headers({
      'Accept': 'application/json',
      'Content-Type': contentType
    });

    if (this.authService.isAuthenticated()) {
      return this.authService.getIdToken().then(
        (token) => {
          if (token){
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
        .fromPromise(this.buildHeaders())
        .switchMap((headers) => this.http.get(url, { headers: headers }));
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


  securedGet(url) {
    return Observable
        .fromPromise(this.buildAuthHeaders())
        .switchMap((headers) => this.http.get(url, { headers: headers }));
  }

  securedPost(url, data) {
    return Observable
        .fromPromise(this.buildAuthHeaders())
        .switchMap((headers) => this.http.post(url, data, { headers: headers }));
  }

  securedPut(url, data, contentType: string = 'application/json') {
    return Observable
        .fromPromise(this.buildAuthHeaders(contentType))
        .switchMap((headers) => this.http.put(url, data, { headers: headers }));
  }

  securedPatch(url, data, contentType: string = 'application/json') {
    return Observable
        .fromPromise(this.buildAuthHeaders(contentType))
        .switchMap((headers) => this.http.patch(url, data, { headers: headers }));
  }

  securedDelete(url) {
    return Observable
        .fromPromise(this.buildAuthHeaders())
        .switchMap((headers) => this.http.delete(url, { headers: headers }));
  }

}
