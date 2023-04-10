import { TransactionEntity } from 'api';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { popModal } from 'store/ducks/modals';
import QRCode from 'qrcode.react';
import CopyToClipboard from 'react-copy-to-clipboard';

type Props = {
  transaction: TransactionEntity;
  close: popModal;
};

type State = {
  displayHelp: boolean;
  copyHelpEmailText: string;
  copyAmountText: string;
  copyAddressText: string;
};

const mapStateToProps = (state) => ({
  transaction: state.payments.currentTransaction,
});

const mapDispatchToProps = (dispatch) => ({
  close: () => dispatch(popModal()),
});

@connect(mapStateToProps, mapDispatchToProps)
export default class extends Component<Props, State> {
  constructor(props: Props) {
    super(props);

    this.state = {
      copyAmountText: 'Copy amount',
      copyAddressText: 'Copy address',
      copyHelpEmailText: 'Copy email',
      displayHelp: false,
    };
  }

  render() {
    const { transaction } = this.props;

    return (
      <form>
        <h2 className="text-center">Pay with Bitcoin</h2>
        <div className="text-center mt20 mb10 text-medium">
          Bitcoin is a <span className="text-red"> digital currency.</span><br/>
        </div>
        <div className="text-center mb20">
          <a href="https://wirexapp.com/" target="_blank">
            Use Wirex to easily purchase and send Bitcoin
          </a>
        </div>
        <div className="purchase-payment-plan-details">
          { transaction.isPending && this.renderTransactionPending() }
          { transaction.fullyPaid && this.renderTransactionCompleted() }
        </div>
      </form>
    );
  }

  renderTransactionPending() {
    const { transaction } = this.props;

    return (
      <Fragment>
        { transaction.isPartiallyPaid && <div className="text-center text-red mt30">Partial payment received.</div> }
        <div className={ `text-center clear ${transaction.isPartiallyPaid ? 'mt10 text-red' : 'mt30'}` }>
          { transaction.isPartiallyPaid ? 'Send remaining amount to complete purchase.' : 'Send' }
        </div>
        <div className="text-center clear margin-auto">
          <span className="text-large mr10 purchase-bitcoin-amount">{ transaction.remainingAmountAsBitcoin }</span>
          <CopyToClipboard text={ transaction.remainingAmountAsBitcoin } onCopy={ this.onCopyAmount }>
            <span className="copy-to-clipboard" style={ { marginTop: '-10px' } }>{ this.state.copyAmountText }</span>
          </CopyToClipboard>
        </div>
        <div className="text-center clear mt20">to address</div>
        <div className="text-center mt10 clear">{ transaction.paymentAddress }</div>
        <div className="text-center pt10 mt20" style={ { backgroundColor: '#ffffff' } }>
          <QRCode value={ `bitcoin:${transaction.paymentAddress}?amount=${transaction.remainingAmountAsBitcoin}` }/>
          <CopyToClipboard text={ transaction.paymentAddress } onCopy={ this.onCopyAddress }>
            <span className="copy-to-clipboard">
              { this.state.copyAddressText }
            </span>
          </CopyToClipboard>
        </div>
        <div className="text-center mt20 mb20 pl20 pr20 clear margin-auto">
          After payment has been processed you will receive a confirmation message here and an email confirming
          &nbsp;your plan purchase.
        </div>
        { this.renderHelp() }
      </Fragment>
    );
  }

  renderHelp() {
    if (!this.state.displayHelp) {
      return (
        <button
          type="button"
          className="red small user-button width-half mt20 mb30 margin-auto"
          onClick={ this.displayHelp }>
          Help?
        </button>
      );
    }

    return (
      <div className="text-center mt20 mb20 pl20 pr20 clear margin-auto">
        For any issues or problems with payments please send an email to
        <span className="text-blue">
          &nbsp;payments@camtube.co
        </span>.
        <CopyToClipboard
          text="payments@camtube.co"
          onCopy={ this.onCopyHelpEmail }>
          <span className="copy-to-clipboard">
            { this.state.copyHelpEmailText }
          </span>
        </CopyToClipboard>
      </div>
    );
  }

  renderTransactionCompleted() {
    return (
      <Fragment>
        <div className="text-center mt20 text-large clear margin-auto">Completed</div>
        <div className="text-center mt20 text-medium clear margin-auto">The payment has been received.</div>
        <div className="text-center mt20 mb20 text-medium clear margin-auto">Your payment plan has been upgraded.</div>
        <div className="text-center mt20 mb20 pr20 pl20 clear margin-auto">
          An email has been sent to your registered email address confirming the purchase or upgrade of your plan.
        </div>
        <button
          type="button"
          className="red small user-button width-half mt30 mb30 margin-auto"
          onClick={ this.props.closeModal }>
          Close
        </button>
      </Fragment>
    );
  }

  displayHelp = () => this.setState({ displayHelp: true });

  onCopyAmount = () => {
    this.setState({ copyAmountText: 'Copied!', copyAddressText: 'Copy address', copyHelpEmailText: 'Copy email' });
  };

  onCopyAddress = () => {
    this.setState({ copyAddressText: 'Copied!', copyAmountText: 'Copy amount', copyHelpEmailText: 'Copy email' });
  };

  onCopyHelpEmail = () => {
    this.setState({ copyHelpEmailText: 'Copied!', copyAmountText: 'Copy amount', copyAddressText: 'Copy address' });
  };
}
