import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { MembershipsModal } from './memberships-modal';

@NgModule({
  declarations: [
    MembershipsModal,
  ],
  imports: [
    IonicPageModule.forChild(MembershipsModal),
  ],
  entryComponents: [
    MembershipsModal,
  ]
})
export class MembershipsModalModule {}
