import { PaymentPlanEntity, SessionEntity, UserEntity } from 'api';
import React, { PureComponent } from 'react';
import fileSize from 'filesize';
import { connect } from 'react-redux';
import moment from 'moment';
import { fetchPaymentPlans } from 'store/ducks/payments';
import { seoUpdate } from 'store/ducks/seo';
import { PaymentPlanPurchaseButton } from './components/payment-plan-purchase-button';

type Props = {
  user: UserEntity;
  session: SessionEntity;
  plans: PaymentPlanEntity[];
  bandwidth: { total: number, remaining: number };
  seo: Function;
};

type State = {
  heightOfPlans: Object;
};

const mapStateToProps = (state) => ({
  session: state.session,
  plans: state.payments.plans,
  user: state.user,
  bandwidth: state.bandwidth,
});

const mapDispatchToProps = (dispatch) => ({
  fetchPaymentPlans: () => dispatch(fetchPaymentPlans()),
  seo: (...args: string) => dispatch(seoUpdate(...args)),
});

@connect(mapStateToProps, mapDispatchToProps)
export default class extends PureComponent<Props, State> {
  constructor(props) {
    super(props);

    this.state = { heightOfPlans: {} };
  }

  componentDidMount() {
    this.props.seo('Payment plans - camtube.co');

    if (this.props.plans.length === 0) {
      this.props.fetchPaymentPlans();
    }

    window.addEventListener('resize', this.calculateHeightsOfElementsOnSameRow);
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.calculateHeightsOfElementsOnSameRow);
  }

  componentDidUpdate(prevProps) {
    if (prevProps.plans.length === 0 && this.props.plans.length > 0) {
      this.calculateHeightsOfElementsOnSameRow();
    }
  }

  calculateHeightsOfElementsOnSameRow = () => {
    // Has type HTMLCollection so array functions are not available
    const elements = document.getElementsByClassName('payment-option');
    const heightOfElementsAtPosition = {};

    // Init
    for (let i = 0; i < elements.length; i++) {
      const element = elements.item(i);
      const pos = element.getBoundingClientRect();

      heightOfElementsAtPosition[pos.top] = { height: 0, elementIds: [] };
    }

    // Process
    for (let i = 0; i < elements.length; i++) {
      const element = elements.item(i);
      const pos = element.getBoundingClientRect();
      heightOfElementsAtPosition[pos.top].elementIds.push(element.id);

      if (element.offsetHeight > heightOfElementsAtPosition[pos.top].height) {
        heightOfElementsAtPosition[pos.top].height = element.offsetHeight;
      }
    }

    const heightOfPlans = {};
    this.props.plans.forEach(p => {
      Object.keys(heightOfElementsAtPosition).forEach(key => {
        const data = heightOfElementsAtPosition[key];

        if (data.elementIds.indexOf(p.uuid)) {
          heightOfPlans[p.uuid] = data.elementIds.length === 1 ? 'auto' : `${data.height}px`;
        }
      });
    });

    this.setState({ heightOfPlans });
  };

  renderPlan = (plan: PaymentPlanEntity) => {
    const userPlan = this.props.plans.find(plan => plan.uuid === this.props.session.paymentPlan);

    return (
      <div className="grid-4 tablet-grid-6 mobile-grid-12" key={ plan.uuid }>
        <div
          id={ plan.uuid }
          className={ `payment-option ${plan.isUpgradedPlan(userPlan) ? 'mobile-upgraded-plan-info' : null}` }
          style={ { height: this.state.heightOfPlans[plan.uuid] } }>
          <h2>{ plan.name }</h2>
          { this.renderDuration(plan.duration) }
          { this.renderGauge(plan.bandwidth) }
          { this.renderPrice(plan) }
          <div className="payment-text">{ plan.description }</div>
          <div className="text-center pb20 pl20 pr20 clear">
            { plan.isUpgradedPlan(userPlan) && 'When upgrading your plan, existing days will be added to your account' }
          </div>
          <PaymentPlanPurchaseButton currentPlan={ userPlan } plan={ plan }/>
        </div>
      </div>
    );
  };

  renderGauge(bandwidth) {
    const percentage = Math.floor(((this.props.bandwidth.total - this.props.bandwidth.remaining) / bandwidth) * 100);
    const [size, notation] = fileSize(bandwidth, { output: 'array', round: 0 });

    return (
      <div className={ 'radial-loader-container p-' + percentage }>
        <div className="radial-loader">
          <span>{ size }</span>
          <div>{ notation }</div>
        </div>
      </div>
    );
  }

  renderDuration(duration: number) {
    if (duration > 0) {
      return <span className="payment-duration">{ duration } days</span>;
    }

    return <span className="payment-duration"/>;
  }

  renderPrice(plan) {
    if (plan.price > 0) {
      return (
        <div>
          <h4 className="no-coin clear">{ plan.price }USD</h4>
          <div className="clear margin-auto mb20 text-center">
            <strong>{ Math.floor(plan.price * 100 / plan.duration) } cents / day</strong>
          </div>
        </div>
      );
    }

    return <h4 className="no-coin">Free</h4>;
  }

  /**
   * Render a sub-header with the remaining days left on the current plan.
   *
   * @returns {*}
   */
  renderPlanExpiration() {
    if (!this.props.user) {
      return false;
    }

    if (this.props.user.paymentPlanSubscribedAt === null) {
      return false;
    }

    const plan = this.props.plans.find((plan: PaymentPlanEntity) => plan.uuid === this.props.session.paymentPlan);
    if (plan === undefined) {
      return false;
    }

    if (plan.duration === 0) {
      return false;
    }

    // convert to moment object
    const subscribedAt = moment(this.props.user.paymentPlanSubscribedAt);

    // calculate the day different between now and the date they subscribed at.
    const daysDiff = subscribedAt.diff(moment(), 'days');

    if (daysDiff === 1) {
      return <h4 className="text-center">You have 1 day remaining on your current plan.</h4>;
    }

    return <h4 className="text-center">You have { plan.duration - daysDiff } days remaining on your current plan.</h4>;
  }

  render() {
    return (
      <div className="container">
        <h2 className="pt15 pb15 text-center thin">Payment plans</h2>
        { this.renderPlanExpiration() }
        { this.props.plans.map(this.renderPlan) }
      </div>
    );
  }
}
