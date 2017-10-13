import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { KeycloakService } from './keycloak/keycloak.service';
import { Observable } from 'rxjs/Rx';

@Injectable()
export class HttpClient {

  constructor(public http: Http, private keycloak: KeycloakService) {}

  buildHeaders(){
    let headers = new Headers({
      'Accept': 'application/json',
      'Content-Type': 'application/json; charset=utf-8'
    })

    return this.keycloak.getToken().then(
      (token) => {
        if(token){
          headers.append('Authorization', 'Bearer ' + token);
          return headers;
        }
      }
    );

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

  put(url, data) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) => this.http.put(url, data, { headers: headers }));
  }

  delete(url) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) => this.http.delete(url, { headers: headers }));
  }

}
