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

  handleRequestError (error) {
    // Auth Token Error
    if(error.status == 401){
      window.location.href = '#/login'
    }

    return Observable.empty()
  }

  get(url) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) =>
          this.http.get(url, { headers: headers })
          .catch(error => this.handleRequestError(error)));
  }

  post(url, data) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) => this.http.post(url, data, { headers: headers })
          .catch(error => this.handleRequestError(error)));
  }

  put(url, data, contentType: string = 'application/json') {
    return Observable
        .fromPromise(this.buildHeaders(contentType))
        .switchMap((headers) => this.http.put(url, data, { headers: headers })
          .catch(error => this.handleRequestError(error)));
  }

  patch(url, data, contentType: string = 'application/json') {
    return Observable
        .fromPromise(this.buildHeaders(contentType))
        .switchMap((headers) => this.http.patch(url, data, { headers: headers })
          .catch(error => this.handleRequestError(error)));
  }

  delete(url) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) => this.http.delete(url, { headers: headers })
          .catch(error => this.handleRequestError(error)));
  }

}
