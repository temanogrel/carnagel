import '@babel/polyfill';
import 'api';
import 'rxjs';
import './styles/style.scss';
import { SessionService } from 'api';
import React, { Fragment } from 'react';
import ReactDOM from 'react-dom';
import { store } from 'store/configure-store';
import { Provider } from 'react-redux';
import { ConnectedRouter } from 'react-router-redux';
import { BandwidthPollInterval, setBandwidthPollInterval } from 'store/ducks/bandwidth';
import history from 'utils/history';
import Router from 'routes';

// Load initial data
SessionService.initSession();

// Set interval for polling bandwidth, every 5 seconds outside recording view
store.dispatch(setBandwidthPollInterval(BandwidthPollInterval.DEFAULT));

// ========================================================
// Render Setup
// ========================================================
const MOUNT_NODE = document.getElementById('root');

let render = () =>
  ReactDOM.render(
    <Fragment>
      <Provider store={ store }>
        <ConnectedRouter history={ history }>
          <Router/>
        </ConnectedRouter>
      </Provider>
    </Fragment>, MOUNT_NODE,
  );

// This code is excluded from production
if (__DEV__) {
  if (module.hot) {
    // Development render functions
    const renderApp = render;
    const renderError = (error) => {
      const RedBox = require('redbox-react').default;

      ReactDOM.render(<RedBox error={ error }/>, MOUNT_NODE);
    };

    // Wrap render in try/catch
    render = () => {
      try {
        renderApp();
      } catch (error) {
        console.error(error);
        renderError(error);
      }
    };

    // Setup hot module replacement
    module.hot.accept('./routes/index', () =>
      setImmediate(() => {
        ReactDOM.unmountComponentAtNode(MOUNT_NODE);
        render();
      }),
    );
  }
}

render();
