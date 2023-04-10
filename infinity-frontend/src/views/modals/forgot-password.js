import { UserService } from 'api/user';
import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Field, reduxForm } from 'redux-form';
import { RenderInput } from 'components';
import { popModal } from 'store/ducks/modals';

const onSubmit = (values) =>
  UserService
    .requestResetPassword(values.usernameOrEmail)
    .catch(err => console.error(err));

const validate = (values) => {
  const errors = {};

  if (!values.usernameOrEmail || values.usernameOrEmail.length === 0) {
    errors.usernameOrEmail = 'Field is required';
  }

  return errors;
};

const render = ({ handleSubmit, error, submitSucceeded, close }) => (
  <form onSubmit={ handleSubmit(onSubmit) }>
    <Field name="usernameOrEmail" component={ RenderInput } type="text" label="Username or email"/>
    { submitSucceeded && (
      <p className="pb10">
        If an account with the provided username or email exists, a confirmation email has been sent to
        the associated email.
      </p>
    ) }
    { !submitSucceeded && <button className="red small user-button width-full mt15 mb15">Reset password</button> }
    <button onClick={close} className="button red transparent small user-button width-full">
      Back
    </button>
  </form>
);

render.propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  error: PropTypes.string,
  submitSucceeded: PropTypes.bool.isRequired,
  close: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => ({
  close: () => dispatch(popModal()),
});

export default connect(null, mapDispatchToProps)(reduxForm({
  form: 'lost-password',
  fields: ['usernameOrEmail'],
  validate,
})(render));
