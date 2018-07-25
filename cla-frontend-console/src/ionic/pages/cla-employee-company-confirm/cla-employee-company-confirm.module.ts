import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaEmployeeCompanyConfirmPage } from './cla-employee-company-confirm';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';

@NgModule({
  declarations: [
    ClaEmployeeCompanyConfirmPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    IonicPageModule.forChild(ClaEmployeeCompanyConfirmPage),
  ],
  entryComponents: [
    ClaEmployeeCompanyConfirmPage
  ]
})
export class ClaEmployeeCompanyConfirmPageModule {}
