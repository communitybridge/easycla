/*
 * Copyright 2017 Red Hat, Inc. and/or its affiliates
 * and other contributors as indicated by the @author tags.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/// <reference path="keycloak.d.ts"/>

import { Injectable } from '@angular/core';

var Keycloak = require('./keycloak'); // load keycloak.js locally
type KeycloakClient = KeycloakModule.KeycloakClient;

@Injectable()
export class KeycloakService {
  static keycloakAuth: KeycloakClient = Keycloak('assets/keycloak.json');

  static init(options?: any): Promise<any> {
    return new Promise((resolve, reject) => {
      KeycloakService.keycloakAuth
        .init(options)
        .success(() => {
          resolve();
        })
        .error((errorData: any) => {
          reject(errorData);
        });
    });
  }

  authenticated = (): boolean => {
    return KeycloakService.keycloakAuth.authenticated;
  };

  login() {
    KeycloakService.keycloakAuth.login();
  }

  logout = () => {
    return KeycloakService.keycloakAuth.logout();
  };

  account() {
    KeycloakService.keycloakAuth.accountManagement();
  }

  createLogoutUrl = () => {
    return KeycloakService.keycloakAuth.createLogoutUrl();
  };

  profile(): Promise<any> {
    return new Promise((resolve, reject) => {
      KeycloakService.keycloakAuth
        .loadUserProfile()
        .success((profile: any) => {
          resolve(profile);
        })
        .error((errorData: any) => {
          reject(errorData);
        });
    });
  }

  getToken(): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      if (KeycloakService.keycloakAuth.token) {
        KeycloakService.keycloakAuth
          .updateToken(5)
          .success(() => {
            resolve(<string>KeycloakService.keycloakAuth.token);
          })
          .error(() => {
            this.login();
            return reject('Failed to refresh token');
          });
      } else {
        this.login();
        return reject('Not logged in');
      }
    });
  }

  getTokenParsed(): Promise<any> {
    return new Promise((resolve, reject) => {
      if (KeycloakService.keycloakAuth.tokenParsed) {
        KeycloakService.keycloakAuth
          .updateToken(5)
          .success(() => {
            resolve(KeycloakService.keycloakAuth.tokenParsed);
          })
          .error(() => {
            this.login();
            return reject('Failed to refresh token');
          });
      } else {
        this.login();
        return reject('Not logged in');
      }
    });
  }
}
