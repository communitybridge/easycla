// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaNewCompanyModal } from './cla-new-company-modal';

// import { ClipboardModule } from 'ngx-clipboard';

@NgModule({
  declarations: [
    ClaNewCompanyModal,
  ],
  imports: [
    // ClipboardModule,
    IonicPageModule.forChild(ClaNewCompanyModal)
  ],
  entryComponents: [
    ClaNewCompanyModal,
  ]
})
export class ClaNewCompanyModalModule {}
