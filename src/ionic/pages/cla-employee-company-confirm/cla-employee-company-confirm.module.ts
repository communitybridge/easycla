import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaEmployeeCompanyConfirmPage } from './cla-employee-company-confirm';

@NgModule({
  declarations: [
    ClaEmployeeCompanyConfirmPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaEmployeeCompanyConfirmPage),
  ],
  entryComponents: [
    ClaEmployeeCompanyConfirmPage
  ]
})
export class ClaEmployeeCompanyConfirmPageModule {}
