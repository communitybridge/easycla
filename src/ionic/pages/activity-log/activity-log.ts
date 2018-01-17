import { Component } from '@angular/core';
import { NavController, IonicPage } from 'ionic-angular';
import { CincoService } from '../../services/cinco.service';
import { RolesService } from '../../services/roles.service';
import { Restricted } from '../../decorators/restricted';

@Restricted({
  roles: ['isAdmin'],
})
@IonicPage({
  segment: 'activity-log'
})
@Component({
  selector: 'activity-log',
  templateUrl: 'activity-log.html'
})
export class ActivityLogPage {
  events: any;

  constructor(
    public navCtrl: NavController,
    private cincoService: CincoService,
    private rolesService: RolesService, // for @Restricted
  ) {
    this.getDefaults();
  }

  getDefaults() {
    this.events = [];
  }

  ngOnInit(){
    this.getEvents();
  }

  getEvents() {
    this.cincoService.getEventsForProject('a0941000002wByYAAU').subscribe(response => {
      if (response) {
        console.log(response);
      }
    });
  }

}
