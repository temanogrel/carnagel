import React, { Fragment } from 'react';
import { HeaderContainer } from 'components';
import { Toaster } from 'containers/toaster';
import { Seo } from 'containers/seo';
import { Modal } from 'views/modals/modal';

export const AppContainer = ({ children }) => (
  <Fragment>
    <Seo/>
    <Toaster/>
    <Modal/>
    <HeaderContainer/>
    { children }
  </Fragment>
);
