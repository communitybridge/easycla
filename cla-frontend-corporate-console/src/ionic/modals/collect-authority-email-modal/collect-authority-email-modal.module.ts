// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { CollectAuthorityEmailModal } from './collect-authority-email-modal';

@NgModule({
  declarations: [
    CollectAuthorityEmailModal,
  ],
  imports: [
    IonicPageModule.forChild(CollectAuthorityEmailModal)
  ],
  entryComponents: [
    CollectAuthorityEmailModal,
  ]
})
export class CollectAuthorityEmailModalModule {}
