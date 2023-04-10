// @flow
import { SessionService, PaymentService } from 'api';
import { combineEpics } from 'redux-observable';
import { MODAL_PURCHASE_PAYMENT_PLAN, pushModal } from 'store/ducks/modals';
import { displayToaster } from 'store/ducks/toaster';
import { push } from 'react-router-redux';
import { PURCHASE_UPDATE_BROADCAST } from 'store/ducks/websocket';
import { resolveCurrentUser } from 'store/ducks/user';
import { Observable } from 'rxjs/Observable';

// ------------------------------------
// Constants
// ------------------------------------
export const FETCH_PAYMENT_PLANS = 'payments:fetch-plans';
export const FETCH_PAYMENT_PLANS_SUCCESS = 'payments:fetch-plans-success';
export const FETCH_TRANSACTION = 'payments:fetch-transactions';
export const FETCH_TRANSACTION_SUCCESS = 'payments:fetch-transaction-success';
export const PURCHASE_PAYMENT_PLAN = 'transaction-plan:purchase';

export const PAYMENT_TRANSACTION_STATE_PENDING = 1;
export const PAYMENT_TRANSACTION_STATE_PARTIALLY_PAID = 2;
export const PAYMENT_TRANSACTION_STATE_FULLY_PAID = 3;
export const PAYMENT_TRANSACTION_STATE_TOO_MUCH_PAID = 4;

export const SATOSHIS_PER_BITCOIN = 100000000;

// ------------------------------------
// ------------------------------------
// Reducer
// ------------------------------------

const INITIAL_STATE = {
  plans: [],
  currentTransaction: null,
};

export default (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case FETCH_PAYMENT_PLANS_SUCCESS:
      return { ...state, plans: action.payload };

    case FETCH_TRANSACTION_SUCCESS:
      return { ...state, currentTransaction: action.payload };

    default:
      return state;
  }
};

/** *****************************************************************
 Selectors
 ****************************************************************** */

/** *****************************************************************
 Action creators
 ****************************************************************** */
export const fetchTransaction = uuid => ({ type: FETCH_TRANSACTION, payload: uuid });
export const fetchPaymentPlans = () => ({ type: FETCH_PAYMENT_PLANS });
export const purchasePaymentPlan = (plan) => ({ type: PURCHASE_PAYMENT_PLAN, payload: plan });

/** *****************************************************************
 Epics
 ****************************************************************** */
const purchasePaymentPlanEpic = action$ =>
  action$
    .ofType(PURCHASE_PAYMENT_PLAN)
    .switchMap(({ payload }) =>
      Observable
        .fromPromise(PaymentService.purchasePaymentPlan(payload.uuid))
        .flatMap(transaction => Observable.concat(
          Observable.of({ type: FETCH_TRANSACTION_SUCCESS, payload: transaction }),
          Observable.of(pushModal(MODAL_PURCHASE_PAYMENT_PLAN)),
        ))
        .catch(() => Observable.of(displayToaster('An unknown error occurred when initializing purchase'))),
    );


const fetchPaymentsPlansEpic = action$ =>
  action$
    .ofType(FETCH_PAYMENT_PLANS)
    .switchMap(() =>
      Observable
        .fromPromise(PaymentService.getPaymentPlans())
        .map(payload => ({ type: FETCH_PAYMENT_PLANS_SUCCESS, payload }))
        .catch(() => Observable.of(push('/'))),
    );

const purchaseUpdatedEpic = action$ =>
  action$
    .ofType(PURCHASE_UPDATE_BROADCAST)
    .flatMap(({ payload }) => Observable.of(fetchTransaction(payload.transactionUuid)));

const fetchTransactionEpic = action$ =>
  action$
    .ofType(FETCH_TRANSACTION)
    .flatMap(({ payload }) =>
      Observable
        .fromPromise(PaymentService.getTransaction(payload))
        .flatMap(transaction => {
          const observables = [
            Observable.of(pushModal(MODAL_PURCHASE_PAYMENT_PLAN)),
            Observable.of({ type: FETCH_TRANSACTION_SUCCESS, payload: transaction }),
          ];

          if (transaction.fullyPaid) {
            observables.push(
              Observable
                .fromPromise(SessionService.renewAuthenticatedSession())
                .mapTo(resolveCurrentUser()),
            );
          }

          return Observable.concat(
            ...observables,
          );
        }),
    );

export const paymentsEpics = combineEpics(
  purchasePaymentPlanEpic,
  fetchTransactionEpic,
  purchaseUpdatedEpic,
  fetchPaymentsPlansEpic,
);
