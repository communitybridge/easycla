// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import {Component} from '@angular/core';
import {Events, IonicPage, NavParams, ViewController} from 'ionic-angular';
import {ClaService} from '../../services/cla.service';

@IonicPage({
  segment: 'cla-configure-github-repositories-modal'
})
@Component({
  selector: 'cla-configure-github-repositories-modal',
  templateUrl: 'cla-configure-github-repositories-modal.html'
})
export class ClaConfigureGithubRepositoriesModal {
  loading: any;
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
    this.getDefaults();

    events.subscribe('modal:close', () => {
      this.dismiss();
    });
  }

  getDefaults() {
    this.loading = {
      repositories: true,
      activateSpinner: false
    }
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
        this.loading.activateSpinner = false;
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
    this.loading.repositories = false
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
    this.loading.activateSpinner = true;

    payload.repository_project_id = this.claProjectId;
    payload.repository_external_id = payload.repository_github_id;
    delete payload.repository_github_id;

    this.claService.postProjectRepository(payload)
      .subscribe(() => this.getOrgRepositories())
  }

  removeRepository(repository) {
    this.loading.activateSpinner = true;
    this.claService.removeProjectRepository(repository.repository_id)
      .subscribe(() => this.getOrgRepositories())
  }

  addAll() {
    console.log(this.orgAndRepositories);
    for (const orgRepo in this.orgAndRepositories) {
      const theOrg = this.orgAndRepositories[orgRepo];
      // No data?
      if (theOrg == null || theOrg.repositories == null || theOrg.repositories == 0) {
        continue;
      }

      for (const repo in theOrg.repositories) {
        const theRepo = theOrg.repositories[repo];
        // No data? - move on
        if (theRepo == null || theRepo.status == null) {
          continue;
        }

        if (theRepo.status == 'free') {
          console.log('Adding repo: ' + theRepo.repository_name);
          this.assignRepository(theRepo);
        } else {
          console.log('Skipping repo: ' + theRepo.repository_name + ', status is: ' + theRepo.status);
        }
      }
    }
  }

  anyAvailableRepos(): boolean {
    let retVal: boolean = false;
    for (const orgRepo in this.orgAndRepositories) {
      const theOrg = this.orgAndRepositories[orgRepo];
      // No data?
      if (theOrg == null || theOrg.repositories == null || theOrg.repositories == 0) {
        continue;
      }

      for (const repo in theOrg.repositories) {
        const theRepo = theOrg.repositories[repo];
        // No data? - move on
        if (theRepo == null || theRepo.status == null) {
          continue;
        }

        if (theRepo.status == 'free') {
          retVal = true;
        }
      }
    }

    return retVal;
  }
}
