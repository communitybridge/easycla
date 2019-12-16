// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { EditCompanyModal } from './edit-company-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';

@NgModule({
  declarations: [EditCompanyModal],
  imports: [IonicPageModule.forChild(EditCompanyModal), LoadingSpinnerComponentModule],
  entryComponents: [EditCompanyModal]
})
export class EditCompanyModule {}
