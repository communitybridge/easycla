// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaOrganizationProviderModal } from './cla-organization-provider-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { ModalHeaderComponentModule } from "../../components/modal-header/modal-header.module";

@NgModule({
  declarations: [
    ClaOrganizationProviderModal
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    ModalHeaderComponentModule,
    IonicPageModule.forChild(ClaOrganizationProviderModal)
  ],
  entryComponents: [
    ClaOrganizationProviderModal
  ]
})
export class ClaOrganizationProviderModalModule {}
