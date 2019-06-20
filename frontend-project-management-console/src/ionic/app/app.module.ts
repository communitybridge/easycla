// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { BrowserModule } from "@angular/platform-browser";
import { NgModule, ErrorHandler } from "@angular/core";
import { HttpModule } from "@angular/http";
import { CurrencyPipe } from "@angular/common";
import { DatePipe } from "@angular/common";
import { IonicApp, IonicModule, IonicErrorHandler } from "ionic-angular";
import { StatusBar } from "@ionic-native/status-bar";
import { SplashScreen } from "@ionic-native/splash-screen";

import { HttpClient } from "../services/http-client";
import { CincoService } from "../services/cinco.service";
import { AnalyticsService } from "../services/analytics.service";
import { S3Service } from "../services/s3.service";
import { RolesService } from "../services/roles.service";
import { XHRBackend, RequestOptions } from "@angular/http";
import { KeycloakService } from "../services/keycloak/keycloak.service";
import {
  KeycloakHttp,
  keycloakHttpFactory
} from "../services/keycloak/keycloak.http";
import { SortService } from "../services/sort.service";
import { FilterService } from "../services/filter.service";
import { ClaService } from "../services/cla.service";
import { AuthService } from "../services/auth.service";
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
    DatePipe,
    HttpClient,
    CincoService,
    AnalyticsService,
    S3Service,
    RolesService,
    KeycloakService,
    SortService,
    FilterService,
    ClaService,
    AuthService,
    { provide: ErrorHandler, useClass: IonicErrorHandler },
    {
      provide: KeycloakHttp,
      useFactory: keycloakHttpFactory,
      deps: [XHRBackend, RequestOptions, AuthService]
    }
  ]
})
export class AppModule {}
