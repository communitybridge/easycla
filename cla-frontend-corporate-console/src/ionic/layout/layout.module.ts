// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import { NgModule } from '@angular/core';
import { ClaFooter } from "./cla-footer/cla-footer";
import {IonicModule} from "ionic-angular";

@NgModule({
  declarations: [
    ClaFooter,
  ],
  imports: [
    IonicModule,
    ],
  exports: [
    ClaFooter,
  ]
})
export class LayoutModule {}
