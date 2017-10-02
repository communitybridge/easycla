import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaOrganizationProviderModal } from './cla-organization-provider-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    ClaOrganizationProviderModal
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(ClaOrganizationProviderModal)
  ],
  entryComponents: [
    ClaOrganizationProviderModal
  ]
})
export class ClaOrganizationProviderModalModule {}
