import { SessionEntity } from 'api/session';
import { Spinner } from 'components/spinner/spinner';
import React, { Component } from 'react';
import { connect } from 'react-redux';

type Props = {
  children: Object;
  session: SessionEntity;
};

const mapStateToProps = state => ({
  session: state.session,
});

@connect(mapStateToProps)
export default class extends Component<Props> {
  render() {
    if (!this.props.session) {
      return <div className="container"><Spinner loading/></div>;
    }

    return this.props.children;
  }
}
