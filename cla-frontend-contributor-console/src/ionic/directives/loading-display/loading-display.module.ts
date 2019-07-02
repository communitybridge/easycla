// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { NgModule } from '@angular/core';
import { LoadingDisplayDirective  } from './loading-display';
@NgModule({
  declarations: [
    LoadingDisplayDirective
  ],
  exports: [
    LoadingDisplayDirective
  ]
})
export class LoadingDisplayDirectiveModule { }
