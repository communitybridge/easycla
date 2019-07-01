// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaLoadingPage } from './cla-loading';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    ClaLoadingPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaLoadingPage),
    LayoutModule
  ],
  entryComponents: [
    ClaLoadingPage
  ]
})
export class ClaLoadingPageModule {}
