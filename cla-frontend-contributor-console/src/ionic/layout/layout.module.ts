// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { ClaFooter } from './cla-footer/cla-footer';
import { IonicModule } from 'ionic-angular';
import { ClaHeader } from './cla-header/cla-header';
import { lfxHeader } from './lfx-header/lfx-header';
import { GetHelpComponentModule } from '../components/get-help/get-help.module';

@NgModule({
  declarations: [ClaFooter, ClaHeader, lfxHeader],
  imports: [IonicModule, GetHelpComponentModule],
  exports: [ClaFooter, ClaHeader]
})
export class LayoutModule { }
