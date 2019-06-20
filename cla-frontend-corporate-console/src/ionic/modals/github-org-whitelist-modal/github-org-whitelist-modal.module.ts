// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { GithubOrgWhitelistModal } from './github-org-whitelist-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { ModalHeaderComponentModule } from "../../components/modal-header/modal-header.module";

@NgModule({
  declarations: [
    GithubOrgWhitelistModal
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    ModalHeaderComponentModule,
    IonicPageModule.forChild(GithubOrgWhitelistModal)
  ],
  entryComponents: [
    GithubOrgWhitelistModal
  ]
})
export class GithubOrgWhitelistModalModule {}
