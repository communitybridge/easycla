import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { CompanyPage } from './company-page';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { SortingDisplayComponentModule } from '../../components/sorting-display/sorting-display.module';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    CompanyPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    SortingDisplayComponentModule,
    LayoutModule,
    IonicPageModule.forChild(CompanyPage)
  ],
  entryComponents: [
    CompanyPage,
  ]
})
export class CompanyPageModule {}
