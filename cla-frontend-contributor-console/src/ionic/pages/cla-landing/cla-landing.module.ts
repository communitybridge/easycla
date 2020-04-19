// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaLandingPage } from './cla-landing';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { LayoutModule } from '../../layout/layout.module';
import { GetHelpComponentModule } from '../../components/get-help/get-help.module';

@NgModule({
  declarations: [ClaLandingPage],
  imports: [
    GetHelpComponentModule,
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(ClaLandingPage),
    LayoutModule
  ],
  entryComponents: [ClaLandingPage]
})
export class ClaLandingPageModule {}
