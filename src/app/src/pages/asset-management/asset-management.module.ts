import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AssetManagementModal } from './asset-management';

@NgModule({
  declarations: [
    AssetManagementModal
  ],
  imports: [
    IonicPageModule.forChild(AssetManagementModal)
  ],
  entryComponents: [
    AssetManagementModal
  ]
})
export class AssetManagementModalModule {}
