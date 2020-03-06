// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

export class CLAGroupModel {
  // Internal project ID
  projectID: string;
  projectName: string;
  // External SFDC project ID
  projectExternalID: string;
  projectACL: string[];
  projectICLAEnabled: boolean;
  projectCCLAEnabled: boolean;
  projectCCLARequiresICLA: boolean;

  /**
   * Constructor for this class.
   *
   * @param projectID the internal project ID
   * @param projectName the project name
   * @param projectExternalID the external project ID (SFDC)
   * @param projectACL the project ACL list
   * @param projectICLAEnabled the ICLA enabled flag
   * @param projectCCLAEnabled the CCLA enabled flag
   * @param projectCCLARequiresICLA the flag to indicate that CCLAs also require ICLAs
   */
  constructor(projectID: string, projectName: string, projectExternalID: string, projectACL: string[], projectICLAEnabled: boolean, projectCCLAEnabled: boolean, projectCCLARequiresICLA: boolean) {
    this.projectID = projectID;
    this.projectName = projectName;
    this.projectExternalID = projectExternalID;
    this.projectACL = projectACL;
    this.projectICLAEnabled = projectICLAEnabled;
    this.projectCCLAEnabled = projectCCLAEnabled;
    this.projectCCLARequiresICLA = projectCCLARequiresICLA;
  }
}
