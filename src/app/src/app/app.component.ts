import { Component, ViewChild } from '@angular/core';
import { Nav, Platform } from 'ionic-angular';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';

import { CincoService } from '../services/cinco.service';

@Component({
  templateUrl: 'app.html',
  providers: [CincoService]
})
export class MyApp {
  @ViewChild(Nav) nav: Nav;

  rootPage: any = 'AllProjectsPage';

  thisUser: any;
  pages: Array<{
    icon?: string,
    access: boolean,
    title: string,
    component: any
  }>;

  constructor(
    public platform: Platform,
    public statusBar: StatusBar,
    public splashScreen: SplashScreen,
    private cincoService: CincoService
  ) {
    this.getDefaults();
    this.initializeApp();
  }

  getDefaults() {
    this.thisUser = {
      isAdmin: false,
      isProjectManager: false,
      isUser: true,
    };
    this.pages = [];
  }

  ngOnInit() {
    this.getUserAccess();
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

  getUserAccess() {
    this.cincoService.getSessionData().subscribe(response => {
      if(response) {
        this.thisUser = response;
        this.regeneratePagesMenu();
      }
    });
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
        access: this.thisUser.isAdmin,
        component: 'ConsoleUsersPage'
      }
    ];
  }
}
