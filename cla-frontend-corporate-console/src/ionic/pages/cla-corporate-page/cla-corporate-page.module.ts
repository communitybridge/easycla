// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaCorporatePage } from './cla-corporate-page';
import { LayoutModule } from '../../layout/layout.module';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { GetHelpComponentModule } from '../../components/get-help/get-help.module';

@NgModule({
  declarations: [ClaCorporatePage],
  imports: [IonicPageModule.forChild(ClaCorporatePage), LayoutModule, GetHelpComponentModule, LoadingSpinnerComponentModule],
  entryComponents: [ClaCorporatePage]
})
export class ClaCorporatePageModule { }
