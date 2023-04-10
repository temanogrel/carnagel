// @flow
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { hideToaster } from 'store/ducks/toaster';

const mapStateToProps = (state) => ({
  toaster: state.toaster,
});

const mapDispatchToProps = (dispatch) => ({
  close: id => dispatch(hideToaster(id)),
});

@connect(mapStateToProps, mapDispatchToProps)
export class Toaster extends Component {
  static propTypes = {
    close: PropTypes.func.isRequired,
    toaster: PropTypes.object,
  };

  close = () => this.props.close(this.props.toaster.id);

  render() {
    const { toaster } = this.props;
    const containerClasses = 'popup-container ' + (toaster ? 'active' : '');

    return (
      <div className={containerClasses}>
        <div className='popup'>
          <div className='popup-indicator'>{ toaster ? toaster.message : '' }</div>
          <a className='popup-close' onClick={this.close}/>
        </div>
      </div>
    );
  }
}
