import { NgModule } from '@angular/core';

import { IonicPageModule } from 'ionic-angular';

import { ActivityLogPage } from './activity-log';

@NgModule({
  declarations: [
    ActivityLogPage
  ],
  imports: [
    IonicPageModule.forChild(ActivityLogPage)
  ],
})
export class ActivityLogPageModule {}
