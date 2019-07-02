// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { ActionPopoverComponent } from './action-popover';
import { IonicPageModule } from 'ionic-angular';

@NgModule({
  declarations: [
    ActionPopoverComponent,
  ],
  imports: [
    IonicPageModule.forChild(ActionPopoverComponent),
  ],
  entryComponents: [
    ActionPopoverComponent,
  ]
})
export class ActionPopoverComponentModule {}
