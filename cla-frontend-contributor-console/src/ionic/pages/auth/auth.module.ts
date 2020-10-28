// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { LayoutModule } from '../../layout/layout.module';
import { AuthPage } from './auth';

@NgModule({
  declarations: [],
  imports: [IonicPageModule.forChild(AuthPage), LayoutModule]
})
export class AuthPageModule { }
