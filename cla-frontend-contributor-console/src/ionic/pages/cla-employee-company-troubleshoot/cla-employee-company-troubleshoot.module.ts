// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaEmployeeCompanyTroubleshootPage } from './cla-employee-company-troubleshoot';
import { LayoutModule } from '../../layout/layout.module';
import { GetHelpComponentModule } from '../../components/get-help/get-help.module';

@NgModule({
  declarations: [ClaEmployeeCompanyTroubleshootPage],
  imports: [IonicPageModule.forChild(ClaEmployeeCompanyTroubleshootPage), LayoutModule, GetHelpComponentModule],
  entryComponents: [ClaEmployeeCompanyTroubleshootPage]
})
export class ClaEmployeeCompanyTroubleshootPageModule {}
