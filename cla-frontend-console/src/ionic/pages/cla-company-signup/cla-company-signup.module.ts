// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaCompanySignupPage } from './cla-company-signup';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    ClaCompanySignupPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaCompanySignupPage),
    LayoutModule
  ],
  entryComponents: [
    ClaCompanySignupPage
  ]
})
export class ClaCompanySignupPageModule {}
