// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { LayoutModule } from '../../layout/layout.module';
import { Authorize } from './authorize';

@NgModule({
  declarations: [Authorize],
  imports: [
    IonicPageModule.forChild(Authorize),
    LayoutModule
  ],
  entryComponents: [Authorize]
})
export class AuthorizeModule { }
