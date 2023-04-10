import React from 'react';
import { Field, reduxForm, SubmissionError } from 'redux-form';
import { RenderInput } from 'components';
import { UserService } from 'api';
import { connect } from 'react-redux';
import { popModal } from 'store/ducks/modals';
import { displayToaster } from 'store/ducks/toaster';

const onSubmit = (values, dispatch) =>
  UserService
    .setNewPassword(values.password, values.token)
    .then(() => {
      dispatch(displayToaster('Your password has been updated'));
      dispatch(popModal());
    })
    .catch(error => {
      throw new SubmissionError({ _error: 'Error setting new password' });
    });

const render = ({ handleSubmit, error, submitSucceeded, closeModal }) => (
  <form onSubmit={ handleSubmit(onSubmit) }>
    <Field name='password' component={ RenderInput } type='password' label='New password'/>
    { error && <p className='pb10 text-red'>Your password link has expired or is invalid.</p> }
    { !error && <button className='red small user-button width-full mb15'>Set password</button> }
    <button onClick={ closeModal } className='button red transparent small user-button width-full'>
      { error ? 'Close' : 'Back' }
    </button>
  </form>
);

const mapDispatchToProps = (dispatch) => ({
  closeModal: () => dispatch(popModal()),
});

export default connect(null, mapDispatchToProps)(reduxForm({
  form: 'set-new-password',
})(render));
