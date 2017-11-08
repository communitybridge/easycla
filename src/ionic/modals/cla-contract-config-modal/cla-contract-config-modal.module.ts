import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaContractConfigModal } from './cla-contract-config-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    ClaContractConfigModal
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(ClaContractConfigModal)
  ],
  entryComponents: [
    ClaContractConfigModal
  ]
})
export class ClaContractConfigModalModule {}
