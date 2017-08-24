import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AllProjectsPage } from './all-projects';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    AllProjectsPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(AllProjectsPage),
  ],
  entryComponents: [
    AllProjectsPage
  ]
})
export class AllProjectsPageModule {}
