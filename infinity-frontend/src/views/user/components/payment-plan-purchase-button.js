import { PaymentPlanEntity, UserEntity } from 'api';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { MODAL_LOGIN, MODAL_REGISTRATION, pushModal } from 'store/ducks/modals';
import { purchasePaymentPlan } from 'store/ducks/payments';
import { displayToaster } from 'store/ducks/toaster';

type Props = {
  user: UserEntity;
  currentPlan: PaymentPlanEntity;
  plan: PaymentPlanEntity;

  openLogin: pushModal;
  openRegistration: pushModal;
  toaster: Function;
  popup: Function;
};

const mapStateToProps = (state) => ({
  user: state.user,
});

const mapDispatchToProps = (dispatch) => ({
  openLogin: () => dispatch(pushModal(MODAL_LOGIN)),
  openRegistration: () => dispatch(pushModal(MODAL_REGISTRATION)),
  toaster: (message) => dispatch(displayToaster(message)),
  purchase: (plan) => dispatch(purchasePaymentPlan(plan)),
});

@connect(mapStateToProps, mapDispatchToProps)
export class PaymentPlanPurchaseButton extends Component<Props> {
  render() {
    if (this.props.plan.uuid === this.props.currentPlan.uuid) {
      return <a className="button disabled large user-button clear pl30 pr30 payment-button">Current</a>;
    }

    if (!this.props.user && this.props.plan.price === 0) {
      return (
        <a
          className="button blue large user-button clear pl30 pr30 payment-button"
          onClick={ this.props.openRegistration }>
          Register
        </a>
      );
    }

    if (this.props.currentPlan.price < this.props.plan.price) {
      return (
        <a
          onClick={!this.props.user ? this.props.openLogin : this.purchase }
          className="button blue large user-button pl30 pr30 clear payment-button">
          { this.props.currentPlan.price > 0 ? 'Upgrade with Bitcoin' : 'Purchase with Bitcoin' }
          { /*<div className="purchase-bitcoin-amount"/>*/ }
        </a>
      );
    }

    return null;
  }

  purchase = () => {
    if (!this.props.user) {
      return this.props.popup('You must sign in before a purchase');
    }

    this.props.purchase(this.props.plan);
  };
}
