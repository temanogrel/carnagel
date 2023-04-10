import { SessionService } from 'api/session';
import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Field, SubmissionError, reduxForm } from 'redux-form';
import {
  MODAL_FORGOT_PASSWORD,
  MODAL_OVERRIDE_EXISTING_SESSION,
  MODAL_REGISTRATION,
  popModal, pushModal,
  replaceModal,
} from 'store/ducks/modals';
import { RenderInput } from 'components';

const onSubmit = (values, dispatch) =>
  SessionService
    .createAuthenticatedSession({ ...values, endOtherSessions: false })
    .then(session => {
      if (!session.blackListedToday) {
        dispatch(popModal());
      }
    })
    .catch((error) => {
      if (error.response && error.response.data.errors && error.response.data.errors['userHasActiveSocketSession']) {
        return dispatch(replaceModal(MODAL_OVERRIDE_EXISTING_SESSION, values));
      }

      throw new SubmissionError({ _error: 'Invalid username or password provided' });
    });

const validate = (values) => {
  const errors = {};

  if (!values.username || values.username.length === 0) {
    errors.username = 'Field is required';
  }

  if (!values.password || values.password.length === 0) {
    errors.password = 'Field is required';
  }

  return errors;
};

const render = ({ handleSubmit, error, submitting, invalid, openForgotPassword, openRegistration }) => (
  <form onSubmit={ handleSubmit(onSubmit) }>
    <Field name="identity" component={ RenderInput } type="text" label="Username"/>
    <Field name="password" component={ RenderInput } type="password" label="Password"/>
    { error && <p className="text-red pb10">{ error }</p> }
    <button type="submit" className="red small user-button width-full mt15 mb15" disabled={ invalid || submitting }>
      Login
    </button>
    <button onClick={openForgotPassword} className="button blue transparent small user-button width-full mb15">
      Forgot password
    </button>
    <button onClick={openRegistration} className="button red transparent small user-button width-full">
      Create account
    </button>
  </form>
);

render.propTypes = {
  invalid: PropTypes.bool.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  error: PropTypes.string,
  submitting: PropTypes.bool.isRequired,
  openForgotPassword: PropTypes.func.isRequired,
  openRegistration: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => ({
  openForgotPassword: () => dispatch(pushModal(MODAL_FORGOT_PASSWORD)),
  openRegistration: () => dispatch(replaceModal(MODAL_REGISTRATION)),
});

export default connect(null, mapDispatchToProps)(reduxForm({
  form: 'login',
  values: ['identity', 'password'],
  validate,
})(render));
