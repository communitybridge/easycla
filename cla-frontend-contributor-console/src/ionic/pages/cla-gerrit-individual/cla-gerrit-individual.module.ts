// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LayoutModule } from '../../layout/layout.module';
import { ClaGerritIndividualPage } from './cla-gerrit-individual';

@NgModule({
  declarations: [ClaGerritIndividualPage],
  imports: [IonicPageModule.forChild(ClaGerritIndividualPage),
    LayoutModule, LoadingSpinnerComponentModule],
  entryComponents: [ClaGerritIndividualPage]
})
export class ClaGerritIndividualPageModule { }
