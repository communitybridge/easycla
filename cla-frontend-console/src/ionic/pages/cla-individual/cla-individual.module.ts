// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaIndividualPage } from './cla-individual';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    ClaIndividualPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaIndividualPage),
    LayoutModule
  ],
  entryComponents: [
    ClaIndividualPage
  ]
})
export class ClaIndividualPageModule {}
