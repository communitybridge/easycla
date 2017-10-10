import { Component, ViewChild } from '@angular/core';
import { Nav, Platform } from 'ionic-angular';
import { StatusBar } from '@ionic-native/status-bar';
import { SplashScreen } from '@ionic-native/splash-screen';

import { SortService } from 'sort-service';

@Component({
  templateUrl: 'app.html',
  providers: [
    SortService,
  ]
})
export class MyApp {
  @ViewChild(Nav) nav: Nav;

  rootPage: any = 'ClaLandingPage';

  constructor(
    public platform: Platform,
    public statusBar: StatusBar,
    public splashScreen: SplashScreen,
  ) {
    this.getDefaults();
    this.initializeApp();
  }

  getDefaults() {
  }

  ngOnInit() {
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
    // Set the nav root so back button doesn't show
    this.nav.setRoot(page.component);
  }
}
