// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaCorporatePage } from './cla-corporate-page';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    ClaCorporatePage,
  ],
  imports: [
    IonicPageModule.forChild(ClaCorporatePage),
    LayoutModule
  ],
  entryComponents: [
    ClaCorporatePage
  ]
})
export class ClaCorporatePageModule {}
