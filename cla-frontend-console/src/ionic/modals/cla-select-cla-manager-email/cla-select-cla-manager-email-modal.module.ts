import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaSelectClaManagerEmailModal } from './cla-employee-request-access-modal';

@NgModule({
  declarations: [
    ClaSelectClaManagerEmailModal,
  ],
  imports: [
    IonicPageModule.forChild(ClaSelectClaManagerEmailModal)
  ],
  entryComponents: [
    ClaSelectClaManagerEmailModal,
  ]
})
export class ClaSelectClaManagerEmailModalModule {}
