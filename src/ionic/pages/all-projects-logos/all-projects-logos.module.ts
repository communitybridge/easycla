import { NgModule } from '@angular/core';

import { IonicPageModule } from 'ionic-angular';

import { AllProjectsLogosPage } from './all-projects-logos';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';

@NgModule({
  declarations: [
    AllProjectsLogosPage
  ],
  imports: [
    SortingDisplayComponentModule,
    IonicPageModule.forChild(AllProjectsLogosPage)
  ],
})
export class AllProjectsLogosPageModule {}
