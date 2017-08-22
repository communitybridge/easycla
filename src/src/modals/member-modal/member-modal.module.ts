import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { MemberModal } from './member-modal';

@NgModule({
  declarations: [
    MemberModal,
  ],
  imports: [
    IonicPageModule.forChild(MemberModal),
  ],
  entryComponents: [
    MemberModal,
  ]
})
export class MemberModalModule {}
