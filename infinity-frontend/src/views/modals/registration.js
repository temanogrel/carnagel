import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { SubmissionError, Field, reduxForm } from 'redux-form';
import { RenderInput } from 'components';
import { popModal } from 'store/ducks/modals';
import { displayToaster } from 'store/ducks/toaster';
import { SessionService, UserService } from 'api';

const onSubmit = (values, dispatch) =>
  UserService
    .register({ email: values.email, username: values.username, password: values.password })
    .then(() =>
      SessionService
        .createAuthenticatedSession({ identity: values.username, password: values.password, endOtherSessions: false })
        .then(session => {
          if (!session.blackListedToday) {
            dispatch(popModal());
          }

          dispatch(displayToaster('Your account was created'));
        }),
    );

const validate = (values) => {
  const errors = {};

  if (!values.email) {
    errors.email = 'This field is required';
  } else if (values.email.indexOf('@') === -1) {
    errors.email = 'Invalid email';
  }

  if (!values.username) {
    errors.username = 'This field is required';
  } else if (values.username.length < 4) {
    errors.username = 'Too short username';
  } else if (values.username.length > 20) {
    errors.username = 'Too long username';
  }

  if (!values.password) {
    errors.password = 'This field is required';
  } else if (values.password.length < 5) {
    errors.password = 'Too short password';
  }

  if (values.password !== values.passwordRepeated) {
    errors.passwordRepated = 'Does not match previous password';
  }

  return errors;
};

const asyncValidate = (values) => {
  const promises = [];

  if (values.email) {
    promises.push(
      UserService.isAvailable('email', values.email).catch(() => {
        return new SubmissionError({ email: 'Email is already in use.' });
      }),
    );
  }

  if (values.username) {
    promises.push(
      UserService.isAvailable('username', values.username).catch(() => {
        return new SubmissionError({ username: 'The username is already in use.' });
      }),
    );
  }

  return Promise
    .all(promises)
    .then(result => {
      const errors = result
        .filter(resp => resp instanceof SubmissionError)
        .map(submissionError => submissionError.errors);

      if (errors.length > 0) {
        return Promise.reject(Object.assign({}, ...errors));
      }
    });
};

const render = ({ handleSubmit, close }) => (
  <form onSubmit={ handleSubmit(onSubmit) }>
    <Field name="username" label="Username" component={ RenderInput } type="text" required/>
    <Field name="email" label="Email" component={ RenderInput } type="text" required/>
    <Field name="password" label="Password" component={ RenderInput } type="password" required/>
    <Field name="passwordRepeated" label="Repeat password" component={ RenderInput } type="password" required/>

    <button type="submit" className="red small user-button width-full mt15 mb15">Sign up</button>
    <button onClick={ close } className="button red transparent small user-button width-full">Back</button>
  </form>
);

render.propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  submitSucceeded: PropTypes.bool.isRequired,
  close: PropTypes.func.isRequired,
};

const mapDispatchToProps = dispatch => ({
  close: () => dispatch(popModal()),
});

export default connect(null, mapDispatchToProps)(reduxForm({
  form: 'user:register',
  validate,
  asyncValidate,
  asyncBlurFields: ['username', 'email'],
})(render));
