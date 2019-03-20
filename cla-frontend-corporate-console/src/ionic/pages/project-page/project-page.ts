import { Component } from "@angular/core";
import {
  NavController,
  ModalController,
  NavParams,
  IonicPage
} from "ionic-angular";
import { ClaService } from "../../services/cla.service";
import { ClaCompanyModel } from "../../models/cla-company";
import { ClaUserModel } from "../../models/cla-user";
import { ClaSignatureModel } from "../../models/cla-signature";
import { SortService } from "../../services/sort.service";
import { RolesService } from "../../services/roles.service";
import { Restricted } from "../../decorators/restricted";
import { WhitelistModal } from "../../modals/whitelist-modal/whitelist-modal";

@Restricted({
  roles: ["isAuthenticated"]
})
@IonicPage({
  segment: "company/:companyId/project/:projectId"
})
@Component({
  selector: "project-page",
  templateUrl: "project-page.html"
})
export class ProjectPage {
  cclaSignature: ClaSignatureModel;
  employeeSignatures: ClaSignatureModel[];
  loading: any;
  companyId: string;
  projectId: string;
  company: ClaCompanyModel;
  manager: ClaUserModel;

  project: any;
  users: any;

  sort: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    private claService: ClaService,
    public modalCtrl: ModalController,
    private rolesService: RolesService, // for @Restricted
    private sortService: SortService
  ) {
    this.companyId = navParams.get("companyId");
    this.projectId = navParams.get("projectId");
    this.getDefaults();
  }

  getDefaults() {
    this.loading = {};
    this.users = {};
    this.sort = {
      date: {
        arrayProp: "date_modified",
        sortType: "date",
        sort: null
      }
    };
    this.company = new ClaCompanyModel();
    this.cclaSignature = new ClaSignatureModel();
  }

  ngOnInit() {
    this.getProject();
    this.getProjectSignatures();
    this.getCompany();
  }

  getCompany() {
    this.claService.getCompany(this.companyId).subscribe(response => {
      this.company = response;
      this.getManager(this.company.company_manager_id);
    });
  }

  getProject() {
    this.claService.getProject(this.projectId).subscribe(response => {
      this.project = response;
    });
  }

  getProjectSignatures() {
    // get CCLA signatures
    this.claService
      .getCompanyProjectSignatures(this.companyId, this.projectId)
      .subscribe(response => {
        let cclaSignatures = response.filter(sig => sig.signature_type === 'ccla');
        if (cclaSignatures.length) {
          this.cclaSignature = cclaSignatures[0];
        }
      });

    // get employee signatures
    this.claService
      .getEmployeeProjectSignatures(this.companyId, this.projectId)
      .subscribe(response => {
        this.employeeSignatures = response;
        for (let signature of this.employeeSignatures) {
          this.getUser(signature.signature_reference_id);
        }
    });

  }

  getManager(userId) {
      this.claService.getUser(userId).subscribe(response => {
        this.manager = response;
      });
  }

  getUser(userId) {
    if (!this.users[userId]) {
      this.claService.getUser(userId).subscribe(response => {
        this.users[userId] = response;
      });
    }
  }

  openWhitelistEmailModal() {
    let modal = this.modalCtrl.create("WhitelistModal", {
      type: "email",
      signatureId: this.cclaSignature.signature_id,
      whitelist: this.cclaSignature.email_whitelist
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getProjectSignatures();
    });
    modal.present();
  }

  openWhitelistDomainModal() {
    let modal = this.modalCtrl.create("WhitelistModal", {
      type: "domain",
      signatureId: this.cclaSignature.signature_id,
      whitelist: this.cclaSignature.domain_whitelist
    });
    modal.onDidDismiss(data => {
      // A refresh of data anytime the modal is dismissed
      this.getProjectSignatures();
    });
    modal.present();
  }

  sortMembers(prop) {
    this.sortService.toggleSort(this.sort, prop, this.employeeSignatures);
  }
}
