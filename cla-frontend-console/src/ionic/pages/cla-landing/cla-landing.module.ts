import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaLandingPage } from './cla-landing';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    ClaLandingPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(ClaLandingPage),
  ],
  entryComponents: [
    ClaLandingPage
  ]
})
export class ClaLandingPageModule {}
