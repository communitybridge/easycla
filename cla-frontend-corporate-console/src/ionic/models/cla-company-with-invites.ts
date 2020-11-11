// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

export class ClaCompanyWithInvitesModel {
  id: string;
  name: string;
  acl: string[];
  status: string;
  created: string;
  updated: string;

  constructor(id: string, name: string, acl: string[], status: string, created: string, updated: string) {
    this.id = id;
    this.name = name;
    this.acl = acl;
    this.status = status
    this.created = created;
    this.updated = updated;
  }
}
