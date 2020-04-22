// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Input, Component, ElementRef, ViewChild, AfterViewChecked } from '@angular/core';
import { ViewController, NavController } from 'ionic-angular';
import { RolesService } from '../../services/roles.service';

@Component({
  selector: 'project-navigation',
  templateUrl: 'project-navigation.html'
})
export class ProjectNavigationComponent implements AfterViewChecked {
  userRoles: any;
  navItems: any;

  scrollChange: number;

  prevMaxScrollOffset = 0;

  @Input('context')
  private context: string;

  @Input('projectId')
  private projectId: any;

  @ViewChild('scrollArea') scrollArea: ElementRef;
  @ViewChild('scrollLeft') scrollLeft: ElementRef;
  @ViewChild('scrollRight') scrollRight: ElementRef;

  constructor(private navCtrl: NavController, private rolesService: RolesService, private viewCtrl: ViewController) {
    this.context;
    this.scrollChange = 40;
    this.getDefaults();
  }

  getDefaults() {
    this.userRoles = this.rolesService.userRoles;
  }

  ngOnInit() {
    this.scrollArea.nativeElement.onscroll = function() {
      this.checkScroll();
    }.bind(this);

    this.rolesService.getUserRolesPromise().then((userRoles) => {
      this.userRoles = userRoles;
      this.generateNavItems();
    });
  }

  ngOnDestroy() {
    // remove scroll listener
  }

  ngAfterViewChecked() {
    let scrollElement = this.scrollArea.nativeElement;
    let maxScrollOffset = scrollElement.scrollWidth - scrollElement.clientWidth;
    if (maxScrollOffset === this.prevMaxScrollOffset) {
    } else {
      this.prevMaxScrollOffset = maxScrollOffset;
      this.checkScroll();
    }
  }

  checkScroll() {
    let scrollElement = this.scrollArea.nativeElement;
    let scrollOffset = scrollElement.scrollLeft;
    if (scrollOffset <= 0) {
      this.scrollLeft.nativeElement.classList.add('disabled');
    } else {
      this.scrollLeft.nativeElement.classList.remove('disabled');
    }
    let maxScrollOffset = scrollElement.scrollWidth - scrollElement.clientWidth;
    if (scrollOffset >= maxScrollOffset) {
      this.scrollRight.nativeElement.classList.add('disabled');
    } else {
      this.scrollRight.nativeElement.classList.remove('disabled');
    }
  }

  navAction(item) {
    this.openPage(item);
  }

  scrollTowardsEnd() {
    this.scrollArea.nativeElement.scrollLeft += this.scrollChange;
  }

  scrollTowardsStart() {
    this.scrollArea.nativeElement.scrollLeft -= this.scrollChange;
  }

  openPage(item) {
    let index = this.navCtrl.indexOf(this.viewCtrl);
    if (index === 0) {
      this.navCtrl.setRoot(item.page, {
        projectId: this.projectId
      });
    } else {
      this.navCtrl
        .push(
          item.page,
          {
            projectId: this.projectId
          },
          {
            animate: false
          }
        )
        .then(() => {
          this.navCtrl.remove(index);
        });
    }
  }

  generateNavItems() {
    this.navItems = [
      {
        label: 'CLA',
        page: 'ProjectClaPage',
        access: this.userRoles.isAdmin
      }
    ];
  }
}
