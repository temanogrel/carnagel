import { SessionEntity } from 'api/session';
import { UserEntity } from 'api/user';
import { push } from 'react-router-redux';
import { Spinner } from 'components/spinner/spinner';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import { displayToaster } from 'store/ducks/toaster';

type Props = {
  children: Object;
  user: UserEntity;
  session: SessionEntity;
  push: push;
  displayToaster: displayToaster;
};

const mapStateToProps = state => ({
  session: state.session,
  user: state.user,
});

const mapDispatchToProps = dispatch => ({
  push: url => dispatch(push(url)),
  displayToaster: message => dispatch(displayToaster(message)),
});

@withRouter
@connect(mapStateToProps, mapDispatchToProps)
export default class extends Component<Props> {
  checkAuthenticatedSession() {
    if (this.props.session.role === 'guest') {
      this.props.displayToaster('Must be signed in to access this url');
      this.props.push('/');
    }
  }

  componentDidMount() {
    this.checkAuthenticatedSession();
  }

  componentDidUpdate() {
    this.checkAuthenticatedSession();
  }

  render() {
    if (!this.props.user) {
      return <div className="container"><Spinner loading/></div>;
    }

    return this.props.children;
  }
}
