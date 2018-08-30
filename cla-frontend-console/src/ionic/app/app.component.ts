import { Component, ViewChild } from "@angular/core";
import { Nav, Platform } from "ionic-angular";
import { StatusBar } from "@ionic-native/status-bar";
import { SplashScreen } from "@ionic-native/splash-screen";
import { ClaService } from "../services/cla.service";

declare const webpackGlobalVars: any;

@Component({
  templateUrl: "app.html",
  providers: []
})
export class MyApp {
  @ViewChild(Nav) nav: Nav;

  rootPage: any = "ClaLandingPage";

  constructor(
    public platform: Platform,
    public statusBar: StatusBar,
    public splashScreen: SplashScreen,
    public claService: ClaService
  ) {
    this.getDefaults();
    this.initializeApp();
    this.claService.setApiUrl(webpackGlobalVars.CLA_API_URL);
  }

  getDefaults() {}

  ngOnInit() {}

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
