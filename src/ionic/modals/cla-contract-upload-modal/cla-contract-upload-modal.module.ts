import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaContractUploadModal } from './cla-contract-upload-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    ClaContractUploadModal
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(ClaContractUploadModal)
  ],
  entryComponents: [
    ClaContractUploadModal
  ]
})
export class ClaContractUploadModalModule {}
