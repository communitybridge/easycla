import React from 'react';

const CLASummary = () => {
  return (
    <div className="section">
      <div className="row">
        <div className="container">
          <div className="row">
            <div className="col-12 col-md-6">
              <h2 className="cla-title">CommunityBridge: CLA</h2>
              <h3 className="cla-subtitle">Ship More Code. Chase Fewer Signatures.</h3>
              <p className="cla-paragraph pb-4">
                For contributors, maintainers, and the companies supporting their developers, Contributor License Agreements (CLAs) can seem
                like they just get in the way of growing the community around a project.
              </p>
              <p className="cla-paragraph pb-3">
                <span><b>EasyCLA</b></span>
                streamlines the process for everyone.
              </p>
              <ul className="cla-list-wrapper">
                <li className="cla-list">
                  <span>
                    <span>Coders can code more quickly by reducing manual steps to get themselves authorized.</span>
                  </span>
                </li>

                <li className="cla-list">
                  <span>
                    <span>Corporations and projects can save time by reducing manual steps managing CLAs and their signatures</span>
                  </span>
                </li>

                <li className="cla-list">
                  <span>
                    <span>Currently available on all Linux Foundation hosted projects<strong>
                      <br /></strong></span>
                  </span>
                </li>
              </ul>
            </div>
            <div className="col-12 col-md-6">
              <div className="">
                <picture>
                  <img src="../static/img/undraw_easycla.svg" alt="" />
                </picture>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CLASummary;