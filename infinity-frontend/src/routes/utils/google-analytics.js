import { Component } from 'react';
import ReactGA from 'react-ga';
import { withRouter } from 'react-router-dom';

@withRouter
export class GoogleAnalytics extends Component {
  componentDidUpdate(prevProps) {
    if (this.props.location.pathname !== prevProps.location.pathname) {
      ReactGA.set({ page: window.location.pathname });
      ReactGA.pageview(window.location.pathname);
    }
  }

  render() {
    return this.props.children;
  }
}
