// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { Events, IonicPage, NavParams, ViewController } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

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
    };
  }

  ngOnInit() {
    this.getOrgRepositories();
  }

  getOrgRepositories() {
    this.claService
      .getProjectConfigurationAndRepos(this.claProjectId)
      .subscribe((data) => this.checkAssignedRepositories(data));
  }

  checkAssignedRepositories(data) {
    //console.log('Received response: ');
    //console.log(data);
    this.assignedRepositories = data['repositories'];
    this.orgAndRepositories = data['orgs_and_repos']
      .map((organization) => {
        this.loading.activateSpinner = false;
        organization.repositories.map((repository) => {
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
      })
      .sort((a, b) => {
        return a.organization_name.trim().localeCompare(b.organization_name.trim());
      });
    this.loading.repositories = false;
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

  isAssignedToLocalContractGroup(repository) {
    return String(repository.repository_project_id) === String(this.claProjectId);
  }

  isTaken(repository) {
    return this.assignedRepositories.some((assignedRepository) => {
      if (String(assignedRepository.repository_external_id) === String(repository.repository_github_id)) {
        repository.repository_project_id = assignedRepository.repository_project_id;
        repository.repository_id = assignedRepository.repository_id;
        return true;
      }
      return false;
    });
  }

  assignRepository(repository) {
    const payload = { ...repository };
    delete payload.status;
    this.loading.activateSpinner = true;

    payload.repository_project_id = this.claProjectId;
    payload.repository_external_id = payload.repository_github_id;
    delete payload.repository_github_id;

    this.claService.postProjectRepository(payload).subscribe(() => this.getOrgRepositories());
  }

  removeRepository(repository) {
    this.loading.activateSpinner = true;
    this.claService.removeProjectRepository(repository.repository_id).subscribe(() => this.getOrgRepositories());
  }

  /**
   * Add all available repositories.
   */
  addAll(organizationName: string) {
    //console.log(this.orgAndRepositories);
    for (const orgRepo in this.orgAndRepositories) {
      const theOrg = this.orgAndRepositories[orgRepo];
      // No data?
      if (theOrg == null || theOrg.repositories == null || theOrg.repositories == 0) {
        //console.log('addAll() - skipping organization/repo at: ' + orgRepo);
        continue;
      }

      //console.log('addAll() - processing organization: ' + theOrg.organizationName);

      if (theOrg.organization_name != organizationName) {
        //console.log('addAll() - skipping organization: ' + organizationName + ', does not match: ' + theOrg.organization_name);
        continue;
      }

      for (const repo in theOrg.repositories) {
        const theRepo = theOrg.repositories[repo];
        // No data? - move on
        if (theRepo == null || theRepo.status == null) {
          //console.log('addAll() - skipping organization: ' + organizationName + ', repository or repository status is empty.');
          console.log(theOrg);
          continue;
        }

        if (theRepo.status == 'free') {
          console.log('addAll() - Adding repo: ' + theRepo.repository_name);
          this.assignRepository(theRepo);
        } else {
          console.log('addAll() - Skipping repo: ' + theRepo.repository_name + ', status is: ' + theRepo.status);
        }
      }
    }
  }

  /**
   * Remove all available repositories.
   */
  removeAll(organizationName: string) {
    //console.log(this.orgAndRepositories);
    for (const orgRepo in this.orgAndRepositories) {
      const theOrg = this.orgAndRepositories[orgRepo];
      // No data?
      if (theOrg == null || theOrg.repositories == null || theOrg.repositories == 0) {
        //console.log('removeAll() - skipping organization/repo at: ' + orgRepo);
        continue;
      }

      //console.log('removeAll() - skipping organization: ' + organizationName + ', does not match: ' + theOrg.organization_name);

      if (theOrg.organization_name != organizationName) {
        //console.log('removeAll() - skipping organization: ' + organizationName);
        continue;
      }

      for (const repo in theOrg.repositories) {
        const theRepo = theOrg.repositories[repo];
        // No data? - move on
        if (theRepo == null || theRepo.status == null) {
          //console.log('removeAll() - skipping organization: ' + organizationName + ', repository or repository status is empty.');
          console.log(theOrg);
          continue;
        }

        if (theRepo.status == 'assigned') {
          console.log('removeAll() - Removing repo: ' + theRepo.repository_name);
          this.removeRepository(theRepo);
        } else {
          console.log('removeAll() - Skipping repo: ' + theRepo.repository_name + ', status is: ' + theRepo.status);
        }
      }
    }
  }

  anyAvailableRepos(organizationName: string): boolean {
    let retVal: boolean = false;
    for (const orgRepo in this.orgAndRepositories) {
      const theOrg = this.orgAndRepositories[orgRepo];
      // No data?
      if (theOrg == null || theOrg.repositories == null || theOrg.repositories == 0) {
        continue;
      }

      if (theOrg.organization_name != organizationName) {
        //console.log('anyAvailableRepos() - skipping organization: ' + organizationName + ', does not match: ' + theOrg.organization_name);
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

  anyAvailableRemoveRepos(organizationName: string): boolean {
    let retVal: boolean = false;
    for (const orgRepo in this.orgAndRepositories) {
      const theOrg = this.orgAndRepositories[orgRepo];
      // No data?
      if (theOrg == null || theOrg.repositories == null || theOrg.repositories == 0) {
        continue;
      }

      if (theOrg.organization_name != organizationName) {
        //console.log('anyAvailableRemoveRepos() - skipping organization: ' + organizationName + ', does not match: ' + theOrg.organization_name);
        continue;
      }

      for (const repo in theOrg.repositories) {
        const theRepo = theOrg.repositories[repo];
        // No data? - move on
        if (theRepo == null || theRepo.status == null) {
          continue;
        }

        if (theRepo.status == 'assigned') {
          retVal = true;
        }
      }
    }

    return retVal;
  }
}
