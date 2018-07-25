import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaCompanySignupPage } from './cla-company-signup';

@NgModule({
  declarations: [
    ClaCompanySignupPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaCompanySignupPage),
  ],
  entryComponents: [
    ClaCompanySignupPage
  ]
})
export class ClaCompanySignupPageModule {}
