// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component, ViewChild } from '@angular/core';
import { Nav, Platform } from 'ionic-angular';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';
import { ClaService } from '../services/cla.service';
import { EnvConfig } from '../services/cla.env.utils';

import { AuthService } from '../services/auth.service';
import { AuthPage } from '../pages/auth/auth';
import { HttpClient } from '../services/http-client';
import { KeycloakService } from '../services/keycloak/keycloak.service';
import { LfxHeaderService } from '../services/lfx-header.service';

@Component({
  templateUrl: 'app.html',
  providers: []
})
export class MyApp {
  @ViewChild(Nav) nav: Nav;

  rootPage: any = AuthPage;

  constructor(
    public platform: Platform,
    public statusBar: StatusBar,
    public splashScreen: SplashScreen,
    public claService: ClaService,
    public authService: AuthService,
    public httpClient: HttpClient,
    public keycloak: KeycloakService,
    private lfxHeaderService: LfxHeaderService
  ) {
    this.initializeApp();

    // Determine if we're running in a local services (developer) mode - the USE_LOCAL_SERVICES environment variable
    // will be set to 'true', otherwise we're using normal services deployed in each environment
    const localServicesMode = (process.env.USE_LOCAL_SERVICES || 'false').toLowerCase() === 'true';
    // Set true for local debugging using localhost (local ports set in claService)
    this.claService.isLocalTesting(localServicesMode);
    this.claService.setApiUrl(EnvConfig['cla-api-url']);
    this.claService.setHttp(httpClient);
  }

  ngOnInit() {
    if (EnvConfig['lfx-header-enabled'] === "true") {
      this.mounted();
    }
  }

  initializeApp() {
    this.platform.ready().then(() => { });
  }

  mounted() {
    const script = document.createElement('script');
    script.setAttribute(
      'src',
      EnvConfig['lfx-header']
    );
    document.head.appendChild(script);
  }

  openPage(page) {
    // Set the nav root so back button doesn't show
    this.nav.setRoot(page.component);
  }
}
