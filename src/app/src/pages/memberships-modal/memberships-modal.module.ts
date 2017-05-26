import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { MembershipsModal } from './memberships-modal';
// import { ComponentsModule } from '../../components/components.modules';

@NgModule({
  declarations: [
    MembershipsModal,
  ],
  imports: [
    IonicPageModule.forChild(MembershipsModal),
    // ComponentsModule
  ],
  entryComponents: [
    MembershipsModal,
  ]
})
export class MembershipsModalModule {}
