// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ConsoleUserUpdateModal } from './console-user-update-modal';

@NgModule({
  declarations: [ConsoleUserUpdateModal],
  imports: [IonicPageModule.forChild(ConsoleUserUpdateModal)],
  entryComponents: [ConsoleUserUpdateModal]
})
export class ConsoleUserUpdateModalModule {}
