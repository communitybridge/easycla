import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaManagerOnboardingPage } from './cla-manager-onboarding';

@NgModule({
  declarations: [
    ClaManagerOnboardingPage,
  ],
  imports: [
    IonicPageModule.forChild(ClaManagerOnboardingPage),
  ],
})
export class ClaManagerOnboardingPageModule {}
