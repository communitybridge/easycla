// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: AGPL-3.0-or-later

import {Component} from '@angular/core';
import {NavParams, ViewController, IonicPage, Events} from 'ionic-angular';
import {ClaService} from '../../services/cla.service';

@IonicPage({
  segment: 'cla-configure-github-repositories-modal'
})
@Component({
  selector: 'cla-configure-github-repositories-modal',
  templateUrl: 'cla-configure-github-repositories-modal.html'
})
export class ClaConfigureGithubRepositoriesModal {
  responseErrors: string[] = [];
  claProjectId: any;
  orgAndRepositories: any[];
  assignedRepositories: any[];

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public claService: ClaService,
    public events: Events
  ) {
    this.claProjectId = this.navParams.get('claProjectId');

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  ngOnInit() {
    this.getOrgRepositories();
  }

  getOrgRepositories() {
    this.claService.getProjectConfigurationAndRepos(this.claProjectId)
      .subscribe(data => this.checkAssignedRepositories(data));
  }

  checkAssignedRepositories(data) {
    this.assignedRepositories = data['repositories'];
    this.orgAndRepositories = data['orgs_and_repos']
      .map(organization => {
        organization.repositories.map(repository => {
          repository.status = 'free';
          repository.repository_organization_name = organization.organization_name;

          if (this.isTaken(repository)) {
            repository.status = 'taken';

            if (this.isAssignedToLocalContractGroup(repository)) {
              repository.status = 'assigned';
            }
          }

          return repository;
        });

        return organization;
      });
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  isAssignedToLocalContractGroup(repository) {
    return String(repository.repository_project_id) === String(this.claProjectId);
  }

  isTaken(repository) {
    return this.assignedRepositories.some(assignedRepository => {
      if (String(assignedRepository.repository_external_id) === String(repository.repository_github_id)) {
        repository.repository_project_id = assignedRepository.repository_project_id;
        repository.repository_id = assignedRepository.repository_id;

        return true;
      }

      return false;
    });
  }

  assignRepository(repository) {
    const payload = {...repository};
    delete payload.status;

    payload.repository_project_id = this.claProjectId;
    payload.repository_external_id = payload.repository_github_id;
    delete payload.repository_github_id;

    this.claService.postProjectRepository(payload)
      .subscribe(() => this.getOrgRepositories())
  }

  removeRepository(repository) {
    this.claService.removeProjectRepository(repository.repository_id)
      .subscribe(() => this.getOrgRepositories())
  }

}
