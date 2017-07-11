import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { ProjectPage } from './project';

@NgModule({
  declarations: [
    ProjectPage,
  ],
  imports: [
    IonicPageModule.forChild(ProjectPage)
  ],
  entryComponents: [
    ProjectPage,
  ]
})
export class ProjectPageModule {}
