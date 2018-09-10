import { BrowserModule } from "@angular/platform-browser";
import { NgModule, ErrorHandler } from "@angular/core";
import { HttpModule } from "@angular/http";
import { CurrencyPipe } from "@angular/common";
import { IonicApp, IonicModule, IonicErrorHandler } from "ionic-angular";
import { StatusBar } from "@ionic-native/status-bar";
import { SplashScreen } from "@ionic-native/splash-screen";

import { HttpClient } from "../services/http-client";
import { CincoService } from "../services/cinco.service";
import { RolesService } from "../services/roles.service";
import { KeycloakService } from "../services/keycloak/keycloak.service";
import {
  KeycloakHttp,
  KEYCLOAK_HTTP_PROVIDER
} from "../services/keycloak/keycloak.http";
import { SortService } from "../services/sort.service";
import { ClaService } from "../services/cla.service";
import { AuthService } from "../services/auth.service";

import { MyApp } from "./app.component";

@NgModule({
  declarations: [MyApp],
  imports: [BrowserModule, HttpModule, IonicModule.forRoot(MyApp)],
  bootstrap: [IonicApp],
  entryComponents: [MyApp],
  providers: [
    StatusBar,
    SplashScreen,
    CurrencyPipe,
    HttpClient,
    CincoService,
    RolesService,
    KeycloakService,
    KEYCLOAK_HTTP_PROVIDER,
    SortService,
    ClaService,
    { provide: ErrorHandler, useClass: IonicErrorHandler }
  ]
})
export class AppModule {}
