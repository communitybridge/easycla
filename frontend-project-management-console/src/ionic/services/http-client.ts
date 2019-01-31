import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Observable } from 'rxjs/Rx';
import { HttpErrorResponse, HttpResponse } from '@angular/common/http';
import { catchError, retry } from 'rxjs/Operators';


@Injectable()
export class HttpClient {

  constructor(
    private http: Http
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

  private handleError(error: HttpErrorResponse, errorHandlerObject: object) {
    if (error.error instanceof ErrorEvent) {
        /* 
            this would be on the browser... I would need an instance of this in
            order to think of a situation to use this...
        */
    } else if(error.status) {
        /*
            more then likely this will come up when it receieves an error from an 
            actual server side call...
        */
        if(
            errorHandlerObject[error.status]
        ){
            // if you have the error code defined... call it...
            errorHandlerObject[error.status](error);
        }
        /* 
            IF the errorObject passed in has a default instance you can 
            at least gracefully catch the error 
        */ 
        errorHandlerObject["default"] && errorHandlerObject["default"](error);
    }
  };

  get(url, errorHandlerObject) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) => this.http.get(url, { headers: headers })
            .pipe(
                catchError((res) => {
                    this.handleError(res, errorHandlerObject);
                    return Observable.empty();
                }
            ))
        );
  }

  post(url, data, errorHandlerObject) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) => this.http.post(url, data, { headers: headers })
            .pipe(
                catchError((res) => {
                    this.handleError(res, errorHandlerObject);
                    return Observable.empty();
                }
            ))
        );
  }

  put(url, data, errorHandlerObject, contentType: string = 'application/json') {
    return Observable
        .fromPromise(this.buildHeaders(contentType))
        .switchMap((headers) => this.http.put(url, data, { headers: headers })
            .pipe(
                catchError((res) => {
                    this.handleError(res, errorHandlerObject);
                    return Observable.empty();
                }
            ))
        );
  }

  patch(url, data, errorHandlerObject, contentType: string = 'application/json') {
    return Observable
        .fromPromise(this.buildHeaders(contentType))
        .switchMap((headers) => this.http.patch(url, data, { headers: headers })
            .pipe(
                catchError((res) => {
                    this.handleError(res, errorHandlerObject);
                    return Observable.empty();
                }
            ))
        );
  }

  delete(url, errorHandlerObject) {
    return Observable
        .fromPromise(this.buildHeaders())
        .switchMap((headers) => this.http.delete(url, { headers: headers })
            .pipe(
                catchError((res) => {
                    this.handleError(res, errorHandlerObject);
                    return Observable.empty();
                }
            ))
        );
  }

}
