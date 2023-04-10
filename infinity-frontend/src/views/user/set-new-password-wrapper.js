import React, { Component } from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';
import queryString from 'query-string';
import { MODAL_SET_NEW_PASSWORD, pushModal } from 'store/ducks/modals';

type Props = {
  children: Object;
  openSetNewPasswordModal: pushModal;
};

const mapDispatchToProps = dispatch => ({
  openSetNewPasswordModal: token => dispatch(pushModal(MODAL_SET_NEW_PASSWORD, { token })),
});

@withRouter
@connect(null, mapDispatchToProps)
export default class extends Component<Props> {
  componentDidMount() {
    const params = queryString.parse(this.props.location.search);

    if (params.token && params.token !== '') {
      this.props.openSetNewPasswordModal(params.token);
    }
  }

  render() {
    return this.props.children;
  }
}
