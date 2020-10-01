// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';
import * as auth0 from 'auth0-js';
import * as jwt_decode from 'jwt-decode';
import { getAuthURLFromWindow } from 'src/app/config/auth-utils';
import { AppSettings } from 'src/app/config/app-settings';
import { EnvConfig } from 'src/app/config/cla-env-utils';

(window as any).global = window;

@Injectable()
export class AuthService {
  auth0;

  constructor(
  ) { }

  public login(type): void {
    this.auth0 = new auth0.WebAuth({
      clientID: EnvConfig.default['auth0-clientId'],
      domain: EnvConfig.default['auth0-domain'],
      responseType: 'token id_token',
      redirectUri: getAuthURLFromWindow(type),
      scope: 'openid email profile'
    });
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
      }
    });
  }

  public isAuthenticated(): boolean {
    // Check whether the current time is past the
    // access token's expiry time
    const expiresAt = JSON.parse(this.getItem(AppSettings.EXPIRES_AT));
    if (expiresAt) {
      return new Date().getTime() < expiresAt;
    }
    return false;
  }

  public getIdToken(): string {
    const tokenId = JSON.parse(this.getItem(AppSettings.ID_TOKEN));
    return tokenId;
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
      const accessToken = localStorage.getItem('access_token');
      if (!accessToken) {
        reject('Access Token must exist to fetch profile');
      }
      this.auth0.client.userInfo(accessToken, function (err, profile) {
        if (profile) {
          return resolve(profile);
        }
      });
    });
  }

  private setSession(authResult): void {
    // Set the time that the access token will expire at
    const expiresAt = JSON.stringify(authResult.expiresIn * 1000 + new Date().getTime());
    this.setItem(AppSettings.ACCESS_TOKEN, authResult.accessToken);
    this.setItem(AppSettings.ID_TOKEN, authResult.idToken);
    this.setItem(AppSettings.EXPIRES_AT, expiresAt);
    this.setItem(AppSettings.USER_ID, authResult.idTokenPayload.nickname);
    this.setItem(AppSettings.USER_EMAIL, authResult.idTokenPayload.email);
    this.setItem(AppSettings.USER_NAME, authResult.idTokenPayload.name);
  }

  hasTokenValid() {
    const tokenId = this.getIdToken();
    if (tokenId !== undefined && tokenId !== null && tokenId.length > 0) {
      return true;
    }
    return false;
  }

  getItem<T>(key: string): T {
    const result = localStorage.getItem(key);
    let resultJson = null;
    if (result != null) {
      resultJson = result;
    }
    return resultJson;
  }

  setItem<T>(key: string, value: T) {
    localStorage.setItem(key, JSON.stringify(value));
  }

  removeItem<T>(key: string) {
    localStorage.removeItem(key);
  }
}
