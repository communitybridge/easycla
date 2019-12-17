// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaIndividualPage } from './cla-individual';
import { LayoutModule } from '../../layout/layout.module';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';

@NgModule({
  declarations: [ClaIndividualPage],
  imports: [IonicPageModule.forChild(ClaIndividualPage), LayoutModule, LoadingSpinnerComponentModule],
  entryComponents: [ClaIndividualPage]
})
export class ClaIndividualPageModule {}
