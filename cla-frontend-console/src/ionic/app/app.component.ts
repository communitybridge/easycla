import { Component, ViewChild } from "@angular/core";
import { Nav, Platform } from "ionic-angular";
import { StatusBar } from "@ionic-native/status-bar";
import { SplashScreen } from "@ionic-native/splash-screen";
import { ClaService } from "../services/cla.service";
import { EnvConfig } from "../services/cla.env.utils";

import { AuthService } from "../services/auth.service";
import { AuthPage } from "../pages/auth/auth";
import { HttpClient } from "../services/http-client";
import { KeycloakService } from "../services/keycloak/keycloak.service";

@Component({
  templateUrl: "app.html",
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
    public keycloak: KeycloakService 
  ) {
    this.getDefaults();
    this.initializeApp();
    this.claService.setApiUrl(EnvConfig['cla-api-url']);
    this.claService.setHttp(httpClient);

    this.authService.handleAuthentication();
  }

  getDefaults() {}

  ngOnInit() {}

  initializeApp() {
    this.platform.ready().then(() => {
    });
  }

  openPage(page) {
    // Set the nav root so back button doesn't show
    this.nav.setRoot(page.component);
  }
}
