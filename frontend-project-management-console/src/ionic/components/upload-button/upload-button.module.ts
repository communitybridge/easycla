// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { UploadButtonComponent } from './upload-button';

@NgModule({
  declarations: [
    UploadButtonComponent,
  ],
  imports: [
    IonicModule,
  ],
  exports: [
    UploadButtonComponent,
  ],
})
export class UploadButtonComponentModule {}
