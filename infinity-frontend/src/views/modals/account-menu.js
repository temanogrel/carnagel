import React from 'react';
import { Link } from 'react-router-dom';
import { popModal } from 'store/ducks/modals';
import { logout } from 'store/ducks/session';
import { connect } from 'react-redux';

export const render = ({ logout, closeModal }) => (
  <ul className='text-center'>
    <li className='text-large text-large-mobile block pb30'>
      <Link onClick={closeModal} to='/favorites' className='link'>My collection</Link>
    </li>
    <li className='text-large text-large-mobile block pb30'>
      <Link onClick={closeModal} to='/payment-plans' className='link'>Account-status</Link>
    </li>
    <li className='text-large text-large-mobile block pb30'>
      <a onClick={logout} className='link'>Log out</a>
    </li>
  </ul>
);

const mapDispatchToProps = (dispatch) => ({
  closeModal: () => dispatch(popModal()),
  logout: () => {
    dispatch(logout());
    dispatch(popModal());
  },
});

export default connect(null, mapDispatchToProps)(render);
