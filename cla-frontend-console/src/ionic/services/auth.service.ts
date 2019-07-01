// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from "@angular/core";
import * as auth0 from "auth0-js";
import * as jwt_decode from "jwt-decode";
import { getAuthURLFromWindow } from "./auth.utils";
import { EnvConfig } from "./cla.env.utils";
(window as any).global = window;

@Injectable()

export class AuthService {

  auth0 = new auth0.WebAuth({
    clientID: EnvConfig['auth0-clientId'],
    domain: EnvConfig['auth0-domain'],
    responseType: "token id_token",
    redirectUri: getAuthURLFromWindow(),
    scope: 'openid email profile',
  });


  public login(): void {
    this.auth0.authorize();
  }

  /* parseHash method to parse a URL hash fragment when the user is redirected back to your application
   * in order to extract the result of an Auth0 authentication response
   */
  public handleAuthentication(): void {
    this.auth0.parseHash((err, authResult) => {
      if (authResult && authResult.accessToken && authResult.idToken) {
        this.setSession(authResult);
      } else if (err) {
        console.log(err);
        alert(
          `Authentication Error: ${
            err.error
          }. Check the console for further details.`
        );
      }
    });
  }

  private setSession(authResult): void {
    // Set the time that the access token will expire at
    const expiresAt = JSON.stringify(
      authResult.expiresIn * 1000 + new Date().getTime()
    );
    localStorage.setItem("access_token", authResult.accessToken);
    localStorage.setItem("id_token", authResult.idToken);
    localStorage.setItem("expires_at", expiresAt);
  }

  public logout(): void {
    // Remove tokens and expiry time from localStorage
    localStorage.removeItem("access_token");
    localStorage.removeItem("id_token");
    localStorage.removeItem("expires_at");
  }

  public isAuthenticated(): boolean {
    // Check whether the current time is past the
    // access token's expiry time
    const expiresAt = JSON.parse(localStorage.getItem("expires_at") || "{}");
    return new Date().getTime() < expiresAt;
  }

  public getIdToken(): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      if (this.isAuthenticated() && localStorage.getItem("id_token")) {
        resolve(localStorage.getItem("id_token"));
      } else {
        return reject("Id token not found. Please login.");
      }
    });
  }

  public parseIdToken(token: string): Promise<any> {
    return new Promise((resolve, reject) => {
      try {
        resolve(jwt_decode(token));
      } catch (error) {
        return reject(error);
      }
    });
  }

  public getUserInfo(): Promise<any> {
    return new Promise<any>((resolve, reject) => {
      var accessToken = localStorage.getItem('access_token');
      if (!accessToken) {
        reject('Access Token must exist to fetch profile');
      }
      this.auth0.client.userInfo(accessToken, function(err, profile) {
        if(profile) {
          return resolve(profile);
        }
      })
    });
  } 


  
}
