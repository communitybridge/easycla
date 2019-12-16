// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaCorporateMemberOptionsModal } from './cla-corporate-member-options-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [ClaCorporateMemberOptionsModal],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(ClaCorporateMemberOptionsModal)
  ],
  entryComponents: [ClaCorporateMemberOptionsModal]
})
export class ClaCorporateMemberOptionsModalModule {}
