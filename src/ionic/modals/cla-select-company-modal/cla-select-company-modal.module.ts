import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaSelectCompanyModal } from './cla-select-company-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';

@NgModule({
  declarations: [
    ClaSelectCompanyModal,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    IonicPageModule.forChild(ClaSelectCompanyModal)
  ],
  entryComponents: [
    ClaSelectCompanyModal,
  ]
})
export class ClaSelectCompanyModalModule {}
