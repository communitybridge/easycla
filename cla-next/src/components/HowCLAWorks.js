// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import React from 'react';

const HowCLAWorks = () => {
  return (
    <div className="section">
      <div className="section-title">
        <h4>How EasyCLA Works</h4>
      </div>

      <div className="container">
        <div className="row how-cla-works">
          <p>
          Below is a high-level flow of how EasyCLA works. The Project Manager (who can be a project maintainer or someone from the Linux Foundation, depending on how you’ve been set up) starts the process by setting up their preferred CLA.
          </p>

          <p>
          Once EasyCLA has been enabled to enforce agreements, the workflow starts when a Developer attempts to contribution to the project. Detailed user steps can be found in our public documentation on GitHub.
          </p>

        </div>
      </div>
    </div>
  );
};

export default HowCLAWorks;