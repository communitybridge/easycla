// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicModule } from 'ionic-angular';
import { MemberHeaderComponent } from './member-header';
import { LoadingSpinnerComponentModule } from '../loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    MemberHeaderComponent,
  ],
  imports: [
    IonicModule,
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
  ],
  exports: [
    MemberHeaderComponent,
  ]
})
export class MemberHeaderComponentModule {}
