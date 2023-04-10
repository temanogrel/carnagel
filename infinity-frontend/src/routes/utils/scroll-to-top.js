import { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { scrollTo } from 'utils/scroll';

@withRouter
export class ScrollToTop extends Component {
  componentDidUpdate(prevProps) {
    if (this.props.location.pathname !== prevProps.location.pathname) {
      scrollTo(0, 0);
    } else if (this.props.location.search !== prevProps.location.search) {
      scrollTo(0, 0);
    }
  }

  render() {
    return this.props.children;
  }
}
