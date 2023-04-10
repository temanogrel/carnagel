import React from 'react';
import { connect } from 'react-redux';
import { popModal } from 'store/ducks/modals';

const render = ({ close }) => (
  <div>
    <h3 className='text-red'>Another session detected</h3>
    <div className='text-center mt20'>
      Only one session can be active on a account at a given point.
      Another session was detected on your account and you have been signed out.
    </div>
    <button type='button' className='red small user-button width-half mt30 mb30 margin-auto' onClick={ close }>
      I understand
    </button>
  </div>
);

const mapDispatchToProps = (dispatch) => ({
  close: () => dispatch(popModal()),
});

export default connect(null, mapDispatchToProps)(render);
