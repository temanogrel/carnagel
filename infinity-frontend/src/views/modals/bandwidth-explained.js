import React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import fileSize from 'filesize';
import { popModal } from 'store/ducks/modals';

export const render = ({bandwidth, close}) => {
  // pretty print the remaining bandwidth
  const [size, notation] = fileSize(bandwidth.remaining, {output: 'array', round: 0});

  return (
    <div>
      <div className='clear text-center text-medium'>You currently have</div>
      <div className='text-center text-large text-red'>{size} {notation}</div>
      <div className='clear text-center text-medium'>daily bandwidth remaining.</div>
      <div className='text-center mt20'>
        When bandwidth is expended, wait until tomorrow to continue watching videos or
        increase your bandwidth by upgrading your payment plan.
      </div>
      <Link
        to='/payment-plans'
        onClick={close}
        className='button red medium user-button width-full mt30 mb30 margin-auto'>
        View available payment plans
      </Link>
    </div>
  );
};

const mapStateToProps = (state) => ({
  bandwidth: state.bandwidth,
});

const mapDispatchToProps = (dispatch) => ({
  close: () => dispatch(popModal()),
});

export default connect(mapStateToProps, mapDispatchToProps)(render);
