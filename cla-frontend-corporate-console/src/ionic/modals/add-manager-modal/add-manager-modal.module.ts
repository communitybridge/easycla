import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AddManagerModal } from './add-manager-modal';

@NgModule({
  declarations: [
    AddManagerModal,
  ],
  imports: [
    IonicPageModule.forChild(AddManagerModal)
  ],
  entryComponents: [
    AddManagerModal,
  ]
})
export class AddManagerModalModule {}
