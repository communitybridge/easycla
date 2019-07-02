// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { BrowserModule } from "@angular/platform-browser";
import { NgModule, ErrorHandler } from "@angular/core";
import { HttpModule } from "@angular/http";
import { CurrencyPipe } from "@angular/common";
import { IonicApp, IonicModule, IonicErrorHandler } from "ionic-angular";
import { StatusBar } from "@ionic-native/status-bar";
import { SplashScreen } from "@ionic-native/splash-screen";
import { ClaService } from "../services/cla.service";
import { AuthService } from "../services/auth.service";
import { RolesService } from "../services/roles.service";
import { HttpClient } from "../services/http-client";
import { KeycloakService } from "../services/keycloak/keycloak.service";
import {
  KeycloakHttp,
  KEYCLOAK_HTTP_PROVIDER
} from "../services/keycloak/keycloak.http";
import { AuthPage } from "../pages/auth/auth";
import { MyApp } from "./app.component";

@NgModule({
  declarations: [MyApp, AuthPage],
  imports: [BrowserModule, HttpModule, IonicModule.forRoot(MyApp)],
  bootstrap: [IonicApp],
  entryComponents: [MyApp, AuthPage],
  providers: [
    StatusBar,
    SplashScreen,
    CurrencyPipe,
    HttpClient,
    ClaService,
    AuthService,
    KeycloakService,
    KEYCLOAK_HTTP_PROVIDER,
    RolesService,
    { provide: ErrorHandler, useClass: IonicErrorHandler },
  ]
})
export class AppModule {}
