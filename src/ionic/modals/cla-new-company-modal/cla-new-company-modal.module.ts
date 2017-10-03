import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaNewCompanyModal } from './cla-new-company-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';

// import { ClipboardModule } from 'ngx-clipboard';

@NgModule({
  declarations: [
    ClaNewCompanyModal,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    // ClipboardModule,
    IonicPageModule.forChild(ClaNewCompanyModal)
  ],
  entryComponents: [
    ClaNewCompanyModal,
  ]
})
export class ClaNewCompanyModalModule {}
