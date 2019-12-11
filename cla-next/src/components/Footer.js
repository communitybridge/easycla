// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import React from 'react';
import { Link } from '../routes';

const Footer = () => {
  return (
    <div className="cla-footer">
      <div className="container">
        <p>
          Copyright &copy; 2019 The Linux Foundation®. All rights reserved. The Linux Foundation has registered
          trademarks and uses trademarks. For a list of trademarks of The Linux Foundation, please see our
          <Link to="https://www.linuxfoundation.org/trademark-usage/">
            <a target="_blank" href="https://www.linuxfoundation.org/trademark-usage/">
              Trademark Usage{' '}
            </a>
          </Link>
          page. Linux is a registered trademark of Linus Torvalds.
          <Link to="https://communitybridge.dev.platform.linuxfoundation.org/acceptable-use">
            <a target="_blank" href="https://communitybridge.dev.platform.linuxfoundation.org/acceptable-use">
              Acceptable Use Policy
            </a>
          </Link>
          |
          <Link to="https://communitybridge.dev.platform.linuxfoundation.org/service-terms">
            <a target="_blank" href="https://communitybridge.dev.platform.linuxfoundation.org/service-terms">
              Service-Specific Terms
            </a>
          </Link>
          |
          <Link to="https://communitybridge.dev.platform.linuxfoundation.org/platform-use-agreement">
            <a target="_blank" href="https://communitybridge.dev.platform.linuxfoundation.org/platform-use-agreement">
              {' '}
              Platform Use Agreement
            </a>
          </Link>
          |
          <Link to="https://communitybridge.org/guide/">
            <a target="_blank" href="https://communitybridge.org/guide/">
              {' '}
              CommunityBridge People Guide
            </a>
          </Link>
          |
          <Link to="https://www.linuxfoundation.org/privacy/">
            <a target="_blank" href="https://www.linuxfoundation.org/privacy/">
              Privacy Policy
            </a>
          </Link>
          <br />
          <br />
          DocuSign is a registered trademark of DocuSign, Inc
        </p>
      </div>
    </div>
  );
};

export default Footer;
