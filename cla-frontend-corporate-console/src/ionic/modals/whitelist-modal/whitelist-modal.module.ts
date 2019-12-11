// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { WhitelistModal } from './whitelist-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';

@NgModule({
  declarations: [WhitelistModal],
  imports: [IonicPageModule.forChild(WhitelistModal), LoadingSpinnerComponentModule],
  entryComponents: [WhitelistModal]
})
export class WhitelistModalModule {}
