import { Component } from '@angular/core';

import { NavController, IonicPage } from 'ionic-angular';

import { CincoService } from '../../app/services/cinco.service'

@IonicPage({
  segment: 'all-invoices'
})
@Component({
  selector: 'all-invoices',
  templateUrl: 'all-invoices.html'
})
export class AllInvoicesPage {

  constructor(public navCtrl: NavController, private cincoService: CincoService) {

  }

}
