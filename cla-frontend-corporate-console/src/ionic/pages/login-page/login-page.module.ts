// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { LoginPage } from './login-page';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    LoginPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    LayoutModule,
    IonicPageModule.forChild(LoginPage),
  ],
  entryComponents: [
    LoginPage
  ]
})
export class LoginPageModule {}
