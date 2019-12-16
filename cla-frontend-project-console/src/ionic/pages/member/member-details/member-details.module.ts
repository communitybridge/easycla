// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { MemberDetailsPage } from './member-details';
import { LoadingSpinnerComponentModule } from '../../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../../components/sorting-display/sorting-display.module';
import { MemberHeaderComponentModule } from '../../../components/member-header/member-header.module';
import { MemberNavigationComponentModule } from '../../../components/member-navigation/member-navigation.module';

@NgModule({
  declarations: [MemberDetailsPage],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    MemberHeaderComponentModule,
    MemberNavigationComponentModule,
    IonicPageModule.forChild(MemberDetailsPage)
  ],
  entryComponents: [MemberDetailsPage]
})
export class MemberDetailsPageModule {}
