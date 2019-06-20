// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaMessageSentPage } from './cla-message-sent';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    ClaMessageSentPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaMessageSentPage),
    LayoutModule
  ],
  entryComponents: [
    ClaMessageSentPage
  ]
})
export class ClaMessageSentPageModule {}
