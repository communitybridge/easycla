// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { SearchAddContactModal } from './search-add-contact-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';

@NgModule({
  declarations: [
    SearchAddContactModal,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    IonicPageModule.forChild(SearchAddContactModal)
  ],
  entryComponents: [
    SearchAddContactModal,
  ]
})
export class SearchAddContactModalModule {}
