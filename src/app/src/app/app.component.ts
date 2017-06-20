import { Component, ViewChild } from '@angular/core';
import { Nav, Platform } from 'ionic-angular';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';

import { CincoService } from './services/cinco.service';

@Component({
  templateUrl: 'app.html',
  providers: [CincoService]
})
export class MyApp {
  @ViewChild(Nav) nav: Nav;

  rootPage: any = 'ProjectsListPage';

  pages: Array<{
    icon?: string,
    title: string,
    component: any
  }>;

  constructor(public platform: Platform, public statusBar: StatusBar, public splashScreen: SplashScreen) {
    this.initializeApp();

    // used for an example of ngFor and navigation
    this.pages = [
      { title: 'All Projects', component: 'ProjectsListPage' },
      { title: 'Member Companies', component: 'AllMembersPage' },
      { title: 'All Invoices Status', component: 'AllInvoicesPage' },
      {
        icon: 'settings',
        title: 'Account Settings',
        component: 'AccountSettingsPage'
      },
      {
        title: 'Linux Console Users',
        component: 'ConsoleUsersPage'
      }
    ];

  }

  initializeApp() {
    this.platform.ready().then(() => {
      // Okay, so the platform is ready and our plugins are available.
      // Here you can do any higher level native things you might need.
      // this.statusBar.styleDefault();
      // this.splashScreen.hide();
    });
  }

  openPage(page) {
    // Reset the content nav to have just this page
    // we wouldn't want the back button to show in this scenario
    this.nav.setRoot(page.component);
  }
}
