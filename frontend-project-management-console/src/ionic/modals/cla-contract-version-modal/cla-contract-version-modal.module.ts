import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaContractVersionModal } from './cla-contract-version-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    ClaContractVersionModal
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(ClaContractVersionModal)
  ],
  entryComponents: [
    ClaContractVersionModal
  ]
})
export class ClaContractVersionModalModule {}
