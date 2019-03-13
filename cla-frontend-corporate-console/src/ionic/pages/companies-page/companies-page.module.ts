import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { CompaniesPage } from './companies-page';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';
import { LayoutModule } from "../../layout/layout.module";

@NgModule({
  declarations: [
    CompaniesPage,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(CompaniesPage),
    LayoutModule
  ],
  entryComponents: [
    CompaniesPage
  ]
})
export class CompaniesPageModule {}
