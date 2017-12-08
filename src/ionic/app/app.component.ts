import { Component, ViewChild } from '@angular/core';
import { Nav, Platform, App } from 'ionic-angular';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';

import { CincoService } from '../services/cinco.service';
import { KeycloakService } from '../services/keycloak/keycloak.service';
import { RolesService } from '../services/roles.service';
import { ClaService } from 'cla-service';
import { CLA_API_URL } from '../services/constants';

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
  ) {
    this.getDefaults();
    this.initializeApp();
    this.claService.setApiUrl(CLA_API_URL);
  }

  getDefaults() {
    this.pages = [];
    this.userRoles = this.rolesService.userRoles;
    this.regeneratePagesMenu();
  }

  ngOnInit() {
    console.log('inital roles');
    console.log(this.rolesService.userRoles);
    // this.rolesService.getData.subscribe((userRoles) => {
    //   this.userRoles = userRoles;
    //   console.log('roles after subscribe');
    //   console.log(this.rolesService.userRoles);
    //   this.regeneratePagesMenu();
    // });
    // this.rolesService.getUserRoles();
    this.rolesService.getUserRolesPromise().then((userRoles) => {
      this.userRoles = this.rolesService.userRoles;
      console.log('roles after subscribe');
      console.log(this.rolesService.userRoles);
      this.regeneratePagesMenu();
    });

  }

  // authenticated(): boolean {
  //   return this.keycloak.authenticated();
  // }
  //
  // login() {
  //   this.keycloak.login();
  // }
  //
  // logout() {
  //   this.nav.setRoot('LoginPage');
  //   this.nav.popToRoot();
  //   this.keycloak.logout();
  // }
  //
  // account() {
  //   this.keycloak.account();
  // }

  initializeApp() {
    this.platform.ready().then(() => {
      // Okay, so the platform is ready and our plugins are available.
      // Here you can do any higher level native things you might need.
      // this.statusBar.styleDefault();
      // this.splashScreen.hide();
      // console.log('initialize app');
      // if(this.keycloak.authenticated()) {
      //   this.rootPage = 'AllProjectsPage';
      // } else {
      //   this.rootPage = 'LoginPage';
      // }

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
        access: this.userRoles.user,
        component: 'AllProjectsPage'
      },
      {
        title: 'Member Companies',
        access: this.userRoles.user,
        component: 'AllMembersPage'
      },
      {
        title: 'All Invoices Status',
        access: this.userRoles.user,
        component: 'AllInvoicesPage'
      },
      {
        title: 'All Projects Logos',
        access: this.userRoles.user,
        component: 'AllProjectsLogosPage'
      },
      {
        icon: 'settings',
        title: 'Account Settings',
        access: this.userRoles.user,
        component: 'AccountSettingsPage'
      },
      {
        title: 'Linux Console Users',
        access: this.userRoles.admin,
        component: 'ConsoleUsersPage'
      },
      {
        title: 'Sign Out',
        access: this.userRoles.authenticated,
        component: 'LoginPage'
      },
    ];
  }
}
