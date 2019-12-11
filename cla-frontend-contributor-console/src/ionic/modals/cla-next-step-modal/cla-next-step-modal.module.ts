// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaNextStepModal } from './cla-next-step-modal';

@NgModule({
  declarations: [ClaNextStepModal],
  imports: [IonicPageModule.forChild(ClaNextStepModal)],
  entryComponents: [ClaNextStepModal]
})
export class ClaNextStepModalModule {}
