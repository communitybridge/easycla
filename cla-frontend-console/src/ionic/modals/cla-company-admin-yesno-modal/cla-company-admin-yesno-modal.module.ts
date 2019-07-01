// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaCompanyAdminYesnoModal } from './cla-company-admin-yesno-modal';

@NgModule({
  declarations: [
    ClaCompanyAdminYesnoModal,
  ],
  imports: [
    // ClipboardModule,
    IonicPageModule.forChild(ClaCompanyAdminYesnoModal)
  ],
  entryComponents: [
    ClaCompanyAdminYesnoModal,
  ]
})
export class ClaCompanyAdminYesnoModalModule {}
