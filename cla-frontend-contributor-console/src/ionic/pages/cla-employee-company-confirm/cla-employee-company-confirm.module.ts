// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaEmployeeCompanyConfirmPage } from './cla-employee-company-confirm';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LayoutModule } from '../../layout/layout.module';

@NgModule({
  declarations: [ClaEmployeeCompanyConfirmPage],
  imports: [LoadingSpinnerComponentModule, IonicPageModule.forChild(ClaEmployeeCompanyConfirmPage), LayoutModule],
  entryComponents: [ClaEmployeeCompanyConfirmPage]
})
export class ClaEmployeeCompanyConfirmPageModule {}
