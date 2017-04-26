import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { AddProjectPage } from './add-project';

@NgModule({
  declarations: [
    AddProjectPage
  ],
  imports: [
    IonicPageModule.forChild(AddProjectPage)
  ],
})
export class AddProjectPageModule {}
