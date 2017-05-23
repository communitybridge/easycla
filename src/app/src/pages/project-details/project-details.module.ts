import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectDetailsPage } from './project-details';

@NgModule({
  declarations: [
    ProjectDetailsPage
  ],
  imports: [
    IonicPageModule.forChild(ProjectDetailsPage)
  ],
})
export class ProjectDetailsPageModule {}
