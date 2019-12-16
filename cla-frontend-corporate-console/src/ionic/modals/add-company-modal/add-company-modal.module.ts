// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AddCompanyModal } from './add-company-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';

@NgModule({
  declarations: [AddCompanyModal],
  imports: [IonicPageModule.forChild(AddCompanyModal), LoadingSpinnerComponentModule],
  entryComponents: [AddCompanyModal]
})
export class AddCompanyModalModule {}
