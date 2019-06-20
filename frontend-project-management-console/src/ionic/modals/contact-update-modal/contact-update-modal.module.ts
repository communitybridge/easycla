// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ContactUpdateModal } from './contact-update-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { UploadButtonComponentModule } from '../../components/upload-button/upload-button.module';

@NgModule({
  declarations: [
    ContactUpdateModal,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    UploadButtonComponentModule,
    IonicPageModule.forChild(ContactUpdateModal),
  ],
  entryComponents: [
    ContactUpdateModal,
  ]
})
export class ContactUpdateModalModule {}
