import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { CollectAuthorityEmailModal } from './collect-authority-email-modal';

@NgModule({
  declarations: [
    CollectAuthorityEmailModal,
  ],
  imports: [
    IonicPageModule.forChild(CollectAuthorityEmailModal)
  ],
  entryComponents: [
    CollectAuthorityEmailModal,
  ]
})
export class CollectAuthorityEmailModalModule {}
