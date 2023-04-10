import React from 'react';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SessionService } from 'api';
import { popModal } from 'store/ducks/modals';
import { resolveCurrentUser } from 'store/ducks/user';
import { connect } from 'react-redux';

const onLoginSubmit = (values, dispatch) =>
  SessionService
    .createAuthenticatedSession({ ...values, endOtherSessions: true })
    .then(() => {
      dispatch(resolveCurrentUser());
      dispatch(popModal());
    });

const render = ({ handleSubmit, submitting, close }) => (
  <form onSubmit={ handleSubmit(onLoginSubmit) }>
    <div className="text-center text-medium">
      Another session is already active on this account. Do you wish to sign out the other session out and sign
      in here instead?
    </div>
    <button type="submit" className="red medium user-button width-full mt15 mb15" disabled={ submitting }>
      Sign me in here
    </button>
    <button
      type="button"
      className="blue medium user-button width-full mt15 mb15"
      disabled={ submitting }
      onClick={ close }
    >
      Cancel
    </button>
  </form>
);

render.propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  error: PropTypes.string,
  submitting: PropTypes.bool.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  close: () => dispatch(popModal()),
});

export default connect(null, mapDispatchToProps)(reduxForm({
  form: 'existing-session:override',
  values: ['identity', 'password'],
})(render));
