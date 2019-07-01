// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ConsoleUsersPage } from './console-users';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    ConsoleUsersPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(ConsoleUsersPage)
  ],
  entryComponents: [
    ConsoleUsersPage,
  ],
})
export class ConsoleUsersPageModule {}
