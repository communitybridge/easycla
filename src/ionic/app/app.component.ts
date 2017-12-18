import { Component, ViewChild } from '@angular/core';
import { Nav, Platform, App } from 'ionic-angular';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';

import { CincoService } from '../services/cinco.service';
import { KeycloakService } from '../services/keycloak/keycloak.service';
import { RolesService } from '../services/roles.service';
import { ClaService } from 'cla-service';
import { CLA_API_URL } from '../services/constants';
import { HttpClient } from '../services/http-client';

@Component({
  templateUrl: 'app.html',
})
export class MyApp {
  @ViewChild(Nav) nav: Nav;

  rootPage: any = 'AllProjectsPage';

  userRoles: any;
  pages: Array<{
    icon?: string,
    access: boolean,
    title: string,
    component: any
  }>;

  users: any[];

  constructor(
    public platform: Platform,
    public app: App,
    public statusBar: StatusBar,
    public splashScreen: SplashScreen,
    private cincoService: CincoService,
    private keycloak: KeycloakService,
    private rolesService: RolesService,
    public claService: ClaService,
    public httpClient: HttpClient,
  ) {
    this.getDefaults();
    this.initializeApp();
    this.claService.setApiUrl(CLA_API_URL);
    this.claService.setHttp(httpClient);
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
      {
        title: 'Member Companies',
        access: this.userRoles.isPmcUser,
        component: 'AllMembersPage'
      },
      {
        title: 'All Invoices Status',
        access: this.userRoles.isPmcUser,
        component: 'AllInvoicesPage'
      },
      {
        title: 'All Projects Logos',
        access: this.userRoles.isPmcUser,
        component: 'AllProjectsLogosPage'
      },
      {
        icon: 'settings',
        title: 'Account Settings',
        access: this.userRoles.isPmcUser,
        component: 'AccountSettingsPage'
      },
      {
        title: 'Linux Console Users',
        access: this.userRoles.isAdmin,
        component: 'ConsoleUsersPage'
      },
      {
        title: 'Sign Out',
        access: this.userRoles.isAuthenticated,
        component: 'LoginPage'
      },
      {
        title: 'Sign In',
        access: !this.userRoles.isAuthenticated,
        component: 'LoginPage'
      },
    ];
  }
}
