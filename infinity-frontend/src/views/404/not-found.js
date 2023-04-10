import React from 'react';
import { Link } from 'react-router-dom';
import { PerformerSearch } from 'components/search/performer';

export default () => (
  <div className='page-not-found'>
    <div className='text-center'>
      <h3>404!</h3>
      <div className='clear mb20'>This page does not exist.</div>
      <div className='clear mb10'>Go <Link to='/'>home</Link> or try a search:</div>
    </div>
    <PerformerSearch reduxId="embedded" className={'search button red large transparent text-left'}/>
  </div>
);

