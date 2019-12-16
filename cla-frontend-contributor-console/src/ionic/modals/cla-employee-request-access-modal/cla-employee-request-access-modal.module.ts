// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaEmployeeRequestAccessModal } from './cla-employee-request-access-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';

@NgModule({
  declarations: [ClaEmployeeRequestAccessModal],
  imports: [LoadingSpinnerComponentModule, IonicPageModule.forChild(ClaEmployeeRequestAccessModal)],
  entryComponents: [ClaEmployeeRequestAccessModal]
})
export class ClaEmployeeRequestAccessModalModule {}
