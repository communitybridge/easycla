import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ClaManagerOnboardingPage } from './cla-manager-onboarding';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';

@NgModule({
  declarations: [ClaManagerOnboardingPage],
  imports: [IonicPageModule.forChild(ClaManagerOnboardingPage)]
})
export class ClaManagerOnboardingPageModule {}
