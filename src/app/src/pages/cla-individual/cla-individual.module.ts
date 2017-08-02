import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaIndividualPage } from './cla-individual';

@NgModule({
  declarations: [
    ClaIndividualPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaIndividualPage),
  ],
  entryComponents: [
    ClaIndividualPage
  ]
})
export class ClaIndividualPageModule {}
