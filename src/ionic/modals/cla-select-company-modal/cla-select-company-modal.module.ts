import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaSelectCompanyModal } from './cla-select-company-modal';

@NgModule({
  declarations: [
    ClaSelectCompanyModal,
  ],
  imports: [
    IonicPageModule.forChild(ClaSelectCompanyModal)
  ],
  entryComponents: [
    ClaSelectCompanyModal,
  ]
})
export class ClaSelectCompanyModalModule {}
