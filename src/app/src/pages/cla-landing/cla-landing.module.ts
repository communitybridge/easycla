import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaLandingPage } from './cla-landing';

@NgModule({
  declarations: [
    ClaLandingPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaLandingPage),
  ],
  entryComponents: [
    ClaLandingPage
  ]
})
export class ClaLandingPageModule {}
