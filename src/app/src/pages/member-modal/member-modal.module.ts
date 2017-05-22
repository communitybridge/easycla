import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { MemberModal } from './member-modal';
import { ComponentsModule } from '../../components/components.modules';

@NgModule({
  declarations: [
    MemberModal,
  ],
  imports: [
    IonicPageModule.forChild(MemberModal),
    ComponentsModule
  ],
  entryComponents: [
    MemberModal,
  ]
})
export class MemberModalModule {}
