import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaLoadingPage } from './cla-loading';

@NgModule({
  declarations: [
    ClaLoadingPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaLoadingPage),
  ],
  entryComponents: [
    ClaLoadingPage
  ]
})
export class ClaLoadingPageModule {}
