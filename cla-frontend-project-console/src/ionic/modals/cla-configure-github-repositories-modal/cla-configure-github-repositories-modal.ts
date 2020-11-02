// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import { Component } from '@angular/core';
import { Events, IonicPage, NavParams, ViewController } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';
import { PlatformLocation } from '@angular/common';

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
  sfdcProjectId: any;
  orgAndRepositories: any[];
  assignedRepositories: any[];
  errorMessage: string = '';

  constructor(
    public navParams: NavParams,
    public viewCtrl: ViewController,
    public claService: ClaService,
    public events: Events,
    private location: PlatformLocation
  ) {
    this.claProjectId = this.navParams.get('claProjectId');
    this.sfdcProjectId = this.navParams.get('sfdcProjectId');
    this.getDefaults();

    this.location.onPopState(() => {
      this.viewCtrl.dismiss(false);
    });

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
    // console.log(data);
    this.assignedRepositories = data['repositories'];
    this.orgAndRepositories = data['orgs_and_repos']
      .map((organization) => {
        this.loading.activateSpinner = false;
        organization.repositories.map((repository) => {
          repository.status = 'free';
          repository.repository_organization_name = organization.organization_name;

          if (this.isTaken(repository)) {
            repository.status = 'taken';
            this.isAssignedToLocalContractGroup(repository);
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
    if (repository.repository_project_id === String(this.claProjectId)) {
      if (repository.enabled) {
        repository.status = 'assigned';
      } else {
        repository.status = 'free';
      }
    }
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
    this.loading.activateSpinner = true;
    this.errorMessage = '';
    const payload = {
      repositoryExternalID: repository.repository_github_id + '',
      repositoryName: repository.repository_name,
      repositoryOrganizationName: repository.repository_organization_name,
      repositoryProjectID: this.claProjectId,
      repositoryType: repository.repository_type,
      repositoryUrl: repository.repository_url
    };
    this.claService.postProjectRepository(this.sfdcProjectId, payload).subscribe(
      () => {
        this.getOrgRepositories();
      },
      (error) => {
        this.loading.activateSpinner = false;
        if (error._body) {
          this.errorMessage = JSON.parse(error._body).Message;
        }
      }
    );
  }

  removeRepository(repository) {
    this.loading.activateSpinner = true;
    this.errorMessage = '';
    this.claService.removeProjectRepository(this.sfdcProjectId, repository.repository_id).subscribe(
      () => {
        this.getOrgRepositories();
      },
      (error) => {
        this.loading.activateSpinner = false;
        if (error._body) {
          this.errorMessage = JSON.parse(error._body).Message;
        }
      }
    );
  }

  /**
   * Add all available repositories.
   */
  addAll(organizationName: string) {
    for (const orgRepo in this.orgAndRepositories) {
      const theOrg = this.orgAndRepositories[orgRepo];
      if (theOrg == null || theOrg.repositories == null || theOrg.repositories == 0) {
        continue;
      }

      if (theOrg.organization_name != organizationName) {
        continue;
      }

      for (const repo in theOrg.repositories) {
        const theRepo = theOrg.repositories[repo];
        if (theRepo == null || theRepo.status == null) {
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
    for (const orgRepo in this.orgAndRepositories) {
      const theOrg = this.orgAndRepositories[orgRepo];
      if (theOrg == null || theOrg.repositories == null || theOrg.repositories == 0) {
        continue;
      }

      if (theOrg.organization_name != organizationName) {
        continue;
      }

      for (const repo in theOrg.repositories) {
        const theRepo = theOrg.repositories[repo];
        if (theRepo == null || theRepo.status == null) {
          continue;
        }

        if (theRepo.status == 'assigned') {
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
