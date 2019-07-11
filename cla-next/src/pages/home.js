// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

import React from 'react'
import Link from 'next/link'
import Head from '../components/Head'
import Nav from '../components/Nav'
import Hero from '../components/Hero';
import ShapeDivider from '../components/shape-divider';
import CLASummary from '../components/CLASummary';
import CLAStep from '../components/CLAStep';
import HowCLAWorks from '../components/HowCLAWorks';
import Diagram from '../components/Diagram';
import Footer from '../components/Footer';


const Home = () => (
  <div>
    <Head title="Home" />
    <Nav />
    <Hero />
    <CLASummary />
    <CLAStep />
    <HowCLAWorks />
    <Diagram />
    <Footer />
  </div>
)

export default Home
