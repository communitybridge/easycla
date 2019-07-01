// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaSendClaManagerEmailModal } from './cla-send-cla-manager-email-modal';

@NgModule({
  declarations: [
    ClaSendClaManagerEmailModal,
  ],
  imports: [
    IonicPageModule.forChild(ClaSendClaManagerEmailModal)
  ],
  entryComponents: [
    ClaSendClaManagerEmailModal,
  ]
})
export class ClaSendClaManagerEmailModalModalModule {}
