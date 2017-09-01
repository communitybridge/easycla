import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaEmployeeRequestAccessModal } from './cla-employee-request-access-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';

@NgModule({
  declarations: [
    ClaEmployeeRequestAccessModal,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    IonicPageModule.forChild(ClaEmployeeRequestAccessModal)
  ],
  entryComponents: [
    ClaEmployeeRequestAccessModal,
  ]
})
export class ClaEmployeeRequestAccessModalModule {}
