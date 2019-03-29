import { Component,  } from '@angular/core';
import { NavController, NavParams, ViewController, ModalController, IonicPage } from 'ionic-angular';
import { ClaService } from '../../services/cla.service';

@IonicPage({
  segment: 'cla-next-step-modal'
})
@Component({
  selector: 'cla-next-step-modal',
  templateUrl: 'cla-next-step-modal.html',
  providers: [
  ]
})
export class ClaNextStepModal {
  projectId: string;
  userId: string;

  project: any;
  signature: any;

  userIsDone: boolean;
  loading: any;

  signingType: string; // "Gerrit" / "Github" 

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private modalCtrl: ModalController,
    private claService: ClaService,
  ) {
    this.projectId = navParams.get('projectId');
    this.userId = navParams.get('userId');
    this.project = navParams.get('project');
    this.signature = navParams.get('signature');
    this.signingType = navParams.get('signingType');
    this.getDefaults();
  }

  getDefaults() {

    this.loading = {
      icla: true,
    };
  }

  ngOnInit() {
    let requiresIcla = this.project.project_ccla_requires_icla_signature;
    if (!requiresIcla) {
      this.userIsDone = true;
      this.loading.icla = false;
      console.log("no icla required. redirect.");
    } else {
      console.log("icla required");
      this.claService.getLastIndividualSignature(this.userId, this.projectId).subscribe(response => {
        console.log(response);
        console.log('need to get the value for if the latest icla is valid');
        if (response === null) {
          // User has no icla, they need one
          this.userIsDone = false;
        } else {
          // get whether icla is up to date
          if (response.requires_resigning) {
            this.userIsDone = false;
          } else {
            this.userIsDone = true;
          }
        }
        this.loading.icla = false;
      });
    }
  }

  // ClaNextStepsModal modal dismiss
  dismiss() {
    this.viewCtrl.dismiss();
  }

  openIclaPage() {
    this.navCtrl.push('ClaIndividualPage', {
      projectId: this.projectId,
      userId: this.userId,
    });
  }

  gotoRepo() {
    window.open(this.signature.signature_return_url, '_blank');
  }

  closeClaApp() {
    // Scripts may close only the windows that were opened by it.
    window.close();
  }

}
