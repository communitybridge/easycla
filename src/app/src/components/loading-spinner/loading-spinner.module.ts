import { NgModule } from '@angular/core';
// import { IonicPageModule } from 'ionic-angular';
import { IonicModule } from 'ionic-angular';
import { LoadingSpinnerComponent } from './loading-spinner';

@NgModule({
  declarations: [
    LoadingSpinnerComponent,
  ],
  imports: [
    // IonicPageModule.forChild(LoadingSpinnerComponent),
    IonicModule,
  ],
  exports: [
    LoadingSpinnerComponent
  ]
})
export class LoadingSpinnerComponentModule {}
