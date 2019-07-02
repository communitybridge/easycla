// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaGerritModal } from './cla-gerrit-modal';
import { ModalHeaderComponentModule } from "../../components/modal-header/modal-header.module";

@NgModule({
  declarations: [
    ClaGerritModal
  ],
  imports: [
    ModalHeaderComponentModule,
    IonicPageModule.forChild(ClaGerritModal)
  ],
  entryComponents: [
    ClaGerritModal
  ]
})
export class ClaGerritModalModule {}
