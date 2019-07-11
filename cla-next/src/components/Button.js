// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import React from 'react';

const Button = ({ children , variant }) => {
  return (
    <div className="btn-wrapper">
      <button className={`btn ${variant}`}>{children}</button>
    </div>
  );
};

export default Button;
