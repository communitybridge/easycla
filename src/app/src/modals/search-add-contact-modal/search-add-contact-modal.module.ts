import { NgModule } from '@angular/core';
import { IonicPageModule } from 'ionic-angular';
import { SearchAddContactModal } from './search-add-contact-modal';
import { LoadingSpinnerComponentModule } from '../../components/loading-spinner/loading-spinner.module';
import { LoadingDisplayDirectiveModule } from '../../directives/loading-display/loading-display.module';

@NgModule({
  declarations: [
    SearchAddContactModal,
  ],
  imports: [
    LoadingSpinnerComponentModule,
    LoadingDisplayDirectiveModule,
    IonicPageModule.forChild(SearchAddContactModal)
  ],
  entryComponents: [
    SearchAddContactModal,
  ]
})
export class SearchAddContactModalModule {}
