import React from 'react';
import { push } from 'react-router-redux';
import { popModal } from 'store/ducks/modals';
import { connect } from 'react-redux';

const mapDispatchToProps = dispatch => ({
  redirectToPaymentPlans: () => {
    dispatch(popModal());
    dispatch(push('/'));
  },
});

export default connect(null, mapDispatchToProps)(({ redirectToPaymentPlans }) => (
  <div>
    <h3 className='text-red text-center'>Too many sessions</h3>
    <div className='text-center mt20'>
      A maximum of two accounts may be created or logged into each day,
      any additional accounts will have zero bandwidth available for 24 hours.
      To remove this restriction purchase a payment plan.
    </div>
    <button
      type='button' className='red small user-button mobile-grid-12 width-half mt30 mb30 margin-auto'
      onClick={ redirectToPaymentPlans }
    >
      View payment plans
    </button>
  </div>
));
