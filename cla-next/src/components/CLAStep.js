// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import React from 'react';
import Button from './Button';

const CLAStep = () => {
  return (
    <div className="section">
      <div className="section-title">
        <h4>Whom Is EasyCLA For?</h4>
      </div>
      <div className="container">
        <div className="row">
          <div className="col-12 col-md-4">
            <div className="cla-step-col">
              <div className="center-item">
                <div className="nectar_icon_wrap" data-style="default" data-draw="" data-border-thickness="2px" data-padding="20px" data-color="accent-color">
                  <div className="nectar_icon">
                    <i className="iconsmind-Business-ManWoman"></i>
                  </div>
                </div>
                <h4>Developers</h4>
              </div>
              <p>Get started contributing code faster and with less friction.</p>
              <ul>
                <li className="cla-list">Receive an automatic notification in GitHub or Gerrit if you need to be whitelisted</li>
                <li className="cla-list">Sign your Individual CLA with an e-signature</li>
                <li className="cla-list">Start contributing faster with a streamlined authorization workflow for Corporate CLAs</li>
              </ul>
              <div className="center-item">
                <Button variant="primary">Read blog post</Button>
              </div>
            </div>
          </div>
          <div className="col-12 col-md-4">
            <div className="cla-step-col">
              <div className="center-item">
                <div className="nectar_icon_wrap" data-style="default" data-draw="" data-border-thickness="2px" data-padding="20px" data-color="accent-color"
                >
                  <div className="nectar_icon">
                    <i className="iconsmind-Code-Window"></i>
                  </div>
                </div>
                <h4>Projects</h4>
              </div>
              <p>Reduce administrative hassles of supporting the CLA for your project.</p>
              <ul>
                <li className="cla-list">Look in one place for companies and individuals who have signed the CLA</li>
                <li className="cla-list">Support both Individual and Corporate Contributors within a single portal</li>
                <li className="cla-list">Enable Companies to manage authorization of their own developers</li>
              </ul>
              <div className="center-item">
                <Button variant="primary">sign in</Button>
                <Button variant="success">add your project</Button>
              </div>
            </div>
          </div>
          <div className="col-12 col-md-4">
            <div className="cla-step-col">
              <div className="center-item">
                <div className="nectar_icon_wrap" data-style="default" data-draw="" data-border-thickness="2px" data-padding="20px" data-color="accent-color"
                >
                  <div className="nectar_icon">
                    <i className="iconsmind-Building"></i>
                  </div>
                </div>
                <h4>Corporations</h4>
              </div>
              <p>Enable all your developers to contribute code easily and quickly while remaining compliant with contribution policies:</p>
              <ul>
                <li className="cla-list">Whitelist developers based on email, domain, GitHub handle, or GitHub organization</li>
                <li className="cla-list">Enable your signatories and contributors to authorize with e-signature via DocuSign</li>
                <li className="cla-list">
                  <span>Enforce signing of the Corporate CLA by your developers without slowing them down with manual bureaucracy</span>
                </li>
              </ul>
              <div className="center-item">
                <Button variant="primary">Corporations</Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CLAStep;
