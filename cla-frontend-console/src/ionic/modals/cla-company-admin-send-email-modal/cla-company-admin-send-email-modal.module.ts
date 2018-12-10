import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaCompanyAdminSendEmailModal } from './cla-company-admin-send-email-modal';

@NgModule({
  declarations: [
    ClaCompanyAdminSendEmailModal,
  ],
  imports: [
    IonicPageModule.forChild(ClaCompanyAdminSendEmailModal)
  ],
  entryComponents: [
    ClaCompanyAdminSendEmailModal,
  ]
})
export class ClaCompanyAdminSendEmailModalModule {}
