import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AssetManagementModal } from './asset-management';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    AssetManagementModal
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(AssetManagementModal)
  ],
  entryComponents: [
    AssetManagementModal
  ]
})
export class AssetManagementModalModule {}
