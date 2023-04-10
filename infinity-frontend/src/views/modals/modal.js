import React, { Component } from 'react';
import Loadable from 'react-loadable';
import { connect } from 'react-redux';
import {
  MODAL_PURCHASE_PAYMENT_PLAN,
  MODAL_LOGIN,
  MODAL_ACCOUNT,
  MODAL_FORGOT_PASSWORD,
  MODAL_REGISTRATION,
  MODAL_BANDWIDTH_EXPLAINED,
  MODAL_SESSION_EXPIRED,
  MODAL_OVERRIDE_EXISTING_SESSION,
  MODAL_SET_NEW_PASSWORD,
  MODAL_SESSION_BLACKLISTED_TODAY, getCurrentModal,
  popModal,
  Modal as ModalInterface,
} from 'store/ducks/modals';

const loading = () => (<div/>);

const Login = Loadable({
  loading,
  loader: () => import(/* webpackChunkName: 'modal-login' */ './login')
});

const AccountMenu = Loadable({
  loading,
  loader: () => import(/* webpackChunkName: 'modal-account-menu' */ './account-menu')
});

const Registration = Loadable({
  loading,
  loader: () => import(/* webpackChunkName: 'modal-registration' */ './registration')
});

const ForgotPassword = Loadable({
  loading,
  loader: () => import(/* webpackChunkName: 'modal-forgot-password' */ './forgot-password')
});

const OverrideExistingSession = Loadable({
  loading,
  loader: () => import(/* webpackChunkName: 'modal-override-existing-session' */ './override-existing-session')
});

const PurchasePaymentPlan = Loadable({
  loading,
  loader: () => import(/* webpackChunkName: 'modals-purchase-payment-plan' */ './purchase-payment-plan')
});

const BandwidthExplained = Loadable({
  loading,
  loader: () => import(/* webpackChunkName: 'modal-bandwidth-explained' */ './bandwidth-explained')
});

const SessionExpired = Loadable({
  loading,
  loader: () => import(/* webpackChunkName: 'modal-session-expired' */ './session-expired')
});

const SetNewPassword = Loadable({
  loading,
  loader: () => import(/* webpackChunkName: 'modal-set-new-password' */ './set-new-password')
});

const SessionBlacklistedToday = Loadable({
  loading,
  loader: () => import(/* webpackChunkName: 'modal-session-blacklisted-today' */ './session-blacklisted-today')
});


type Props = {
  modal: ModalInterface;
};

const mapStateToProps = (state) => ({
  modal: getCurrentModal(state),
});

const mapDispatchToProps = (dispatch: Function) => ({
  pop: () => dispatch(popModal()),
  dispatch: action => dispatch(action),
});

@connect(mapStateToProps, mapDispatchToProps)
export class Modal extends Component<Props> {
  componentDidMount() {
    document.addEventListener('keydown', this.onEscape);
  }

  componentDidUpdate(prevProps) {
    if (!prevProps.modal && this.props.modal) {
      document
        .querySelector('body')
        .setAttribute('class', 'modal-active');
    } else if (prevProps.modal && !this.props.modal) {
      document
        .querySelector('body')
        .removeAttribute('class');
    }
  }

  onEscape = (evt: KeyboardEvent) => {
    if (!this.props.modal) {
      return;
    }

    if (evt.which === 27) {
      evt.preventDefault();
      this.dismiss();
    }
  };

  dismiss = () => {
    if (!this.props.modal.onClose) {
      return this.props.pop();
    }

    this.props.dispatch(this.props.modal.onClose);
    this.props.pop();
  };

  renderModal() {
    switch (this.props.modal.type) {
      case MODAL_LOGIN:
        return <Login/>;

      case MODAL_ACCOUNT:
        return <AccountMenu/>;

      case MODAL_REGISTRATION:
        return <Registration/>;

      case MODAL_FORGOT_PASSWORD:
        return <ForgotPassword/>;

      case MODAL_PURCHASE_PAYMENT_PLAN:
        return <PurchasePaymentPlan/>;

      case MODAL_BANDWIDTH_EXPLAINED:
        return <BandwidthExplained/>;

      case MODAL_SESSION_EXPIRED:
        return <SessionExpired/>;

      case MODAL_OVERRIDE_EXISTING_SESSION:
        return <OverrideExistingSession initialValues={this.props.modal.payload}/>;

      case MODAL_SET_NEW_PASSWORD:
        return <SetNewPassword initialValues={this.props.modal.payload}/>;

      case MODAL_SESSION_BLACKLISTED_TODAY:
        return <SessionBlacklistedToday/>;

      default:
        return null;
    }
  }

  render() {
    const { modal } = this.props;
    if (!modal) {
      return null;
    }

    return (
      <div className='tint'>
        <div onClick={this.dismiss} className='close-modal'/>
        <div className='modal' style={{'maxWidth': '500px'}}>
          {this.renderModal()}
        </div>
      </div>
    );
  }
}
