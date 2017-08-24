import { Component, ViewChild } from '@angular/core';
import { Nav, Platform, App } from 'ionic-angular';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';

import { CincoService } from '../services/cinco.service';
import { KeycloakService } from '../services/keycloak/keycloak.service';

@Component({
  templateUrl: 'app.html',
})
export class MyApp {
  @ViewChild(Nav) nav: Nav;

  rootPage: any = 'LoginPage';

  userRoles: any;
  pages: Array<{
    icon?: string,
    access: boolean,
    title: string,
    component: any
  }>;

  constructor(
    public platform: Platform,
    public app: App,
    public statusBar: StatusBar,
    public splashScreen: SplashScreen,
    private cincoService: CincoService,
    private keycloak: KeycloakService
  ) {
    this.getDefaults();
    this.initializeApp();
  }

  authenticated(): boolean {
    return this.keycloak.authenticated();
  }

  login() {
    this.keycloak.login();
  }

  logout() {
    this.nav.setRoot('LoginPage');
    this.nav.popToRoot();
    this.keycloak.logout();
  }

  account() {
    this.keycloak.account();
  }

  getDefaults() {
    this.userRoles = {
      isAdmin: false,
      isProgramManager: false,
      isUser: false,
    };
    this.pages = [];
  }

  ngOnInit() {
    this.getUserRoles();
  }

  initializeApp() {
    this.platform.ready().then(() => {
      // Okay, so the platform is ready and our plugins are available.
      // Here you can do any higher level native things you might need.
      // this.statusBar.styleDefault();
      // this.splashScreen.hide();

      if(this.keycloak.authenticated()) {
        this.rootPage = 'AllProjectsPage';
      } else {
        this.rootPage = 'LoginPage';
      }

    });
  }

  openHomePage() {
    this.nav.setRoot(this.rootPage);
  }

  openPage(page) {
    // Set the nav root so back button doesn't show
    this.nav.setRoot(page.component);
  }

  isInArray(roles, role) {
    for(let i=0; i<roles.length; i++) {
      if (roles[i].toLowerCase() === role.toLowerCase()) return true;
    }
    return false;
  }

  getUserRoles() {
    this.keycloak.getTokenParsed().then(
      (tokenParsed) => {
        if(tokenParsed){
          let isAdmin = this.isInArray(tokenParsed.realm_access.roles, 'admin');
          let isProgramManager = this.isInArray(tokenParsed.realm_access.roles, 'program_manager');
          let isUser = this.isInArray(tokenParsed.realm_access.roles, 'user');
          this.userRoles = {
            isAdmin: isAdmin,
            isProgramManager: isProgramManager,
            isUser: isUser,
          };
          this.regeneratePagesMenu();
        }
      }
    )
  }

  regeneratePagesMenu() {
    this.pages = [
      {
        title: 'All Projects',
        access: true,
        component: 'AllProjectsPage'
      },
      {
        title: 'Member Companies',
        access: true,
        component: 'AllMembersPage'
      },
      {
        title: 'All Invoices Status',
        access: true,
        component: 'AllInvoicesPage'
      },
      {
        icon: 'settings',
        title: 'Account Settings',
        access: true,
        component: 'AccountSettingsPage'
      },
      {
        title: 'Linux Console Users',
        access: this.userRoles.isAdmin,
        component: 'ConsoleUsersPage'
      }
    ];
  }
}
