// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';

import { IonicPageModule } from 'ionic-angular';

import { AllInvoicesPage } from './all-invoices';

@NgModule({
  declarations: [
    AllInvoicesPage
  ],
  imports: [
    IonicPageModule.forChild(AllInvoicesPage)
  ],
})
export class AllInvoicesPageModule {}
