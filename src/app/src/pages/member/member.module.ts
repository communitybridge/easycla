import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { MemberPage } from './member';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';

@NgModule({
  declarations: [
    MemberPage
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    IonicPageModule.forChild(MemberPage)
  ],
  entryComponents: [
    MemberPage
  ]
})
export class MemberPageModule {}
