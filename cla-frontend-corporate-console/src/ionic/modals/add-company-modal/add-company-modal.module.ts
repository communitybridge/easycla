// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AddCompanyModal } from './add-company-modal';

@NgModule({
  declarations: [
    AddCompanyModal,
  ],
  imports: [
    IonicPageModule.forChild(AddCompanyModal)
  ],
  entryComponents: [
    AddCompanyModal,
  ]
})
export class AddCompanyModalModule {}
