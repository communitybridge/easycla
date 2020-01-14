// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
import {Http} from '@angular/http';
import { Component, ViewChild } from '@angular/core';
import { Nav, Platform, App, Events } from 'ionic-angular';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';

import { CincoService } from '../services/cinco.service';
import { KeycloakService } from '../services/keycloak/keycloak.service';
import { KeycloakHttp } from '../services/keycloak/keycloak.http';
import { RolesService } from '../services/roles.service';
import { S3Service } from '../services/s3.service';
import { ClaService } from '../services/cla.service';
import { HttpClient } from '../services/http-client';

import { AuthService } from '../services/auth.service';
import { AuthPage } from '../pages/auth/auth';
import { EnvConfig } from '../services/cla.env.utils';

@Component({
  templateUrl: 'app.html'
})
export class MyApp {
  @ViewChild(Nav) nav: Nav;

  rootPage: any = AuthPage;

  userRoles: any;
  pages: Array<{
    icon?: string;
    access: boolean;
    title: string;
    component: any;
  }>;

  users: any[];

  constructor(
    public platform: Platform,
    public app: App,
    public statusBar: StatusBar,
    public splashScreen: SplashScreen,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    private keycloakHttp: KeycloakHttp,
    private rolesService: RolesService,
    public claService: ClaService,
    public s3service: S3Service,
    public http: Http,
    public httpClient: HttpClient,
    public authService: AuthService,
    public events: Events
  ) {
    this.getDefaults();
    this.initializeApp();

    // initialize s3service with dependency injected version of httpClient
    this.s3service.setHttp(httpClient);

    // manually create new httpClient with keycloakHttp
    let kcHttpClient = new HttpClient(http, keycloak, authService);

    // Determine if we're running in a local services (developer) mode - the USE_LOCAL_SERVICES environment variable
    // will be set to 'true', otherwise we're using normal services deployed in each environment
    const localServicesMode = (process.env.USE_LOCAL_SERVICES || 'false').toLowerCase() === 'true';
    // Set true for local debugging using localhost (local ports set in claService)
    this.claService.isLocalTesting(localServicesMode);

    // set authd services to use kcHttpClient
    this.claService.setApiUrl(EnvConfig['cla-api-url']);
    this.claService.setS3LogoUrl(EnvConfig['cla-logo-s3-url']);
    this.claService.setHttp(kcHttpClient);

    this.cincoService.setApiUrl(EnvConfig['cinco-api-url']);
    this.cincoService.setHttp(kcHttpClient);

    this.authService.handleAuthentication();

    events.subscribe('nav:allProjects', () => {
      this.nav.setRoot('AllProjectsPage');
    });
  }

  getDefaults() {
    this.pages = [];
    this.userRoles = this.rolesService.userRoles;
    this.regeneratePagesMenu();
  }

  ngOnInit() {
    this.rolesService.getUserRolesPromise().then((userRoles) => {
      this.userRoles = userRoles;
      this.regeneratePagesMenu();
    });
  }

  initializeApp() {
    this.platform.ready().then(() => {
      // Okay, so the platform is ready and our plugins are available.
      // Here you can do any higher level native things you might need.
      // this.statusBar.styleDefault();
      // this.splashScreen.hide();
    });
  }

  openHomePage() {
    this.nav.setRoot(this.rootPage);
  }

  openPage(page) {
    // Set the nav root so back button doesn't show
    this.nav.setRoot(page.component);
  }

  regeneratePagesMenu() {
    this.pages = [
      {
        title: 'All Projects',
        access: this.userRoles.isPmcUser,
        component: 'AllProjectsPage'
      },
      // {
      //   title: "Member Companies",
      //   access: this.userRoles.isPmcUser,
      //   component: "AllMembersPage"
      // },
      // {
      //   title: "All Invoices Status",
      //   access: this.userRoles.isPmcUser,
      //   component: "AllInvoicesPage"
      // },
      // {
      //   title: "All Projects Logos",
      //   access: this.userRoles.isPmcUser,
      //   component: "AllProjectsLogosPage"
      // },
      // {
      //   title: "Activity Log",
      //   access: this.userRoles.isAdmin,
      //   component: "ActivityLogPage"
      // },
      {
        title: 'Sign Out',
        access: true,
        component: 'LogoutPage'
      },
      {
        title: 'Sign In',
        access: !this.userRoles.isAuthenticated,
        component: 'LoginPage'
      }
    ];
  }
}
