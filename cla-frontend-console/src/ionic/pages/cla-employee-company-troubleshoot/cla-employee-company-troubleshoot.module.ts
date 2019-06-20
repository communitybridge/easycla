// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaEmployeeCompanyTroubleshootPage } from './cla-employee-company-troubleshoot';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    ClaEmployeeCompanyTroubleshootPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaEmployeeCompanyTroubleshootPage),
    LayoutModule
  ],
  entryComponents: [
    ClaEmployeeCompanyTroubleshootPage
  ]
})
export class ClaEmployeeCompanyTroubleshootPageModule {}
