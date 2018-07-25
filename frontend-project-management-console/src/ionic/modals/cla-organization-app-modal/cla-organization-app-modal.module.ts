import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaOrganizationAppModal } from './cla-organization-app-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    ClaOrganizationAppModal
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(ClaOrganizationAppModal)
  ],
  entryComponents: [
    ClaOrganizationAppModal
  ]
})
export class ClaOrganizationAppModalModule {}
