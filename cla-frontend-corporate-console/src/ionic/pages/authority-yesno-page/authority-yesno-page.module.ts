// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AuthorityYesnoPage } from './authority-yesno-page';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    AuthorityYesnoPage,
  ],
  imports: [
    IonicPageModule.forChild(AuthorityYesnoPage),
    LayoutModule
  ],
  entryComponents: [
    AuthorityYesnoPage
  ]
})
export class AuthorityYesnoPageModule {}
