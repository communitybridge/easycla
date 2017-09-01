import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaEmployeeCompanyPage } from './cla-employee-company';

@NgModule({
  declarations: [
    ClaEmployeeCompanyPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaEmployeeCompanyPage),
  ],
  entryComponents: [
    ClaEmployeeCompanyPage
  ]
})
export class ClaEmployeeCompanyPageModule {}
