import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { SearchAddContact } from './search-add-contact';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    SearchAddContact,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(SearchAddContact)
  ],
  entryComponents: [
    SearchAddContact,
  ]
})
export class SearchAddContactModule {}
