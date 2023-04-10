import React from 'react';
import { BarLoader } from 'react-spinners';
import { isMobile } from 'utils/device';

export const Spinner = ({ loading }) => (
  <div className="spinner-container">
    <div className="logo"/>
    <BarLoader width={ isMobile() ? 100 : 200 } height={ 8 } color={ '#fb5255' } loading={ loading }/>
  </div>
);
