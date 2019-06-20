// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaEmployeeRequestAccessModal } from './cla-employee-request-access-modal';

@NgModule({
  declarations: [
    ClaEmployeeRequestAccessModal,
  ],
  imports: [
    IonicPageModule.forChild(ClaEmployeeRequestAccessModal)
  ],
  entryComponents: [
    ClaEmployeeRequestAccessModal,
  ]
})
export class ClaEmployeeRequestAccessModalModule {}
