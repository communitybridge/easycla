// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

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
