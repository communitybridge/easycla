// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Injectable } from '@angular/core';
import * as auth0 from 'auth0-js';
import * as jwt_decode from 'jwt-decode';
import { getAuthURLFromWindow } from './auth.utils';
import { EnvConfig } from './cla.env.utils';
import { App } from 'ionic-angular';
import Timer = NodeJS.Timer;

(window as any).global = window;

@Injectable()
export class AuthService {
  sessionTimer: Timer = null;

  auth0 = new auth0.WebAuth({
    clientID: EnvConfig['auth0-clientId'],
    domain: EnvConfig['auth0-domain'],
    responseType: 'token id_token',
    redirectUri: getAuthURLFromWindow()
  });

  constructor(private app: App) {}

  // constructor(public router: Router) {} Right now haven't figure out how ionic does routing
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
        this.handleSessionTimeout();
      } else if (err) {
        this.warn(err);
        alert(`Authentication Error: ${err.error}. Check the console for further details.`);
      }
    });
  }

  private setSession(authResult): void {
    // Set the time that the access token will expire at
    const expiresAt = JSON.stringify(authResult.expiresIn * 1000 + new Date().getTime());
    // console.log(authResult);
    localStorage.setItem('access_token', authResult.accessToken);
    localStorage.setItem('id_token', authResult.idToken);
    localStorage.setItem('expires_at', expiresAt);
    localStorage.setItem('userid', authResult.idTokenPayload.nickname);
    localStorage.setItem('user_email', authResult.idTokenPayload.email);
    localStorage.setItem('user_name', authResult.idTokenPayload.name);
  }

  private handleSessionTimeout(): void {
    const FIVE_MINUTES_MS = 60000 * 5;
    this.sessionTimer = setInterval(() => {
      //this.debug('checking session...');
      if (!this.isAuthenticated()) {
        // Make a note
        const expiresAt = JSON.parse(localStorage.getItem('expires_at') || '{}');
        this.info('Session expired at ' + new Date(expiresAt).toISOString() + ' - logging out.');

        // Cancel the timer
        if (this.sessionTimer != null) {
          this.clearSessionTimer();
        }

        // Logout the user and redirect to the login page
        //this.debug('logging user out...');
        this.logout();
        //this.debug('redirecting to login page...');
        this.app
          .getRootNav()
          .setRoot('LoginPage')
          .catch(error => this.warn(error));
      }
    }, FIVE_MINUTES_MS);
  }

  private clearSessionTimer(): void {
    clearInterval(this.sessionTimer);
  }

  public logout(): void {
    // Remove tokens and expiry time from localStorage
    localStorage.removeItem('access_token');
    localStorage.removeItem('id_token');
    localStorage.removeItem('expires_at');
    localStorage.removeItem('userid');
    localStorage.removeItem('user_email');
    localStorage.removeItem('user_name');

    const localServicesMode = (process.env.USE_LOCAL_SERVICES || 'false').toLowerCase() === 'true';
    // property is: cla-corp-console-link-{STAGE} -> https://corporate.dev.lfcla.com
    // but ask for the key without the prefix or the suffix/stage
    let redirectUri = EnvConfig['corp-console-link'];
    if (localServicesMode) {
      redirectUri = 'http://localhost:8101';
    }
    console.log('Redirecting to: ' + redirectUri + '/#/login');
    // This library seems to want the options to be `redirectUri`, but the API docs say to use: `returnTo`
    // `returnTo` works
    this.auth0.logout({
      returnTo: redirectUri + '/#/login'
    });
  }

  public isAuthenticated(): boolean {
    // Check whether the current time is past the
    // access token's expiry time
    const expiresAt = JSON.parse(localStorage.getItem('expires_at') || '{}');
    //return new Date().getTime() < expiresAt;
    return new Date().getTime() < expiresAt && localStorage.getItem('userid') != null;
  }

  public getIdToken(): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      if (this.isAuthenticated() && localStorage.getItem('id_token')) {
        resolve(localStorage.getItem('id_token'));
      } else {
        return reject('Id token not found. Please login.');
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

  /**
   * Helper function to show a formatted console log output with the time.
   * @param message the message to print out
   */
  private info(message: string): void {
    console.log('[' + new Date().toISOString() + '][INFO] ' + message);
  }

  /**
   * Helper function to show a formatted console log output with the time.
   * @param message the message to print out
   */
  private debug(message: string): void {
    console.log('[' + new Date().toISOString() + '][DEBUG] ' + message);
  }

  /**
   * Helper function to show a formatted console log output with the time.
   * @param message the message to print out
   */
  private warn(message: string): void {
    console.log('[' + new Date().toISOString() + '][WARN] ' + message);
  }
}
