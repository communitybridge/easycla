import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AssignUserModal } from './assign-user';

@NgModule({
  declarations: [
    AssignUserModal
  ],
  imports: [
    IonicPageModule.forChild(AssignUserModal)
  ],
  entryComponents: [
    AssignUserModal
  ]
})
export class AssignUserModalModule {}
