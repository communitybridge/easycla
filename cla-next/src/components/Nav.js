// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import React from 'react'
import Link from 'next/link'
import HamburgerIcon from './HamburgerIcon';

const leftLinks = [
  {
    href: 'https://funding.communitybridge.org/', label: 'Funding'
  },
  {
    href: 'https://people.communitybridge.org/', label: 'Mentorship'
  },
  {
    href: 'https://docs.linuxfoundation.org/', label: 'Docs'
  },
  {
    href: 'https://support.communitybridge.org/', label: 'Support'
  }
].map(link => {
  link.key = `nav-link-${link.href}-${link.label}`
  return link
})

const rightLinks = [
  {
    href: 'https://twitter.com/linuxfoundation', label: 'twitter'
  },
  {
    href: 'https://www.facebook.com/TheLinuxFoundation', label: 'facebook'
  },
  {
    href: 'https://www.linkedin.com/company/208777', label: 'linkedin'
  },
  {
    href: 'https://www.youtube.com/user/TheLinuxFoundation', label: 'youtube-play'
  }
].map(link => {
  link.key = `nav-link-${link.href}-${link.label}`
  return link
})




const Nav = () => (
  <header>
    <div className="row header-row">
      <div className="container">
        <div className="header-wrapper">
          <div className="header-left-column">
            <nav>
              <ul>
                <li>
                  <Link prefetch href="/">
                    <a className="header-brand">
                      <img src="../static/logo.svg" alt="" />
                    </a>
                  </Link>
                </li>
              </ul>
              <ul className="header-navs d-none d-lg-flex">
                {leftLinks.map(({ key, href, label }) => (
                  <li key={key}>
                    <Link href={href}>
                      <a>{label}</a>
                    </Link>
                  </li>
                ))}
              </ul>
            </nav>
          </div>
          <div className="header-right-column">
            <nav className="d-none d-lg-flex">
              <ul className="header-navs">
                {rightLinks.map(({ key, href, label }) => (
                  <li key={key}>
                    <Link href={href}>
                      <a>
                        <i class={`fa fa-${label} faa-bounce animated-hover`}></i>
                      </a>
                    </Link>
                  </li>
                ))}
              </ul>
              <button className="linux-button">
                <Link href="https://www.linuxfoundation.org/">
                  <a> THE LINUX FOUNDATION</a>
                </Link>
              </button>
            </nav>
            <div className="d-flex d-md-none">
              <HamburgerIcon />
            </div>
          </div>
        </div>
      </div>
    </div>
  </header >
)

export default Nav
