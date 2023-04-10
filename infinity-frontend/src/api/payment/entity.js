import {
  PAYMENT_TRANSACTION_STATE_FULLY_PAID, PAYMENT_TRANSACTION_STATE_PARTIALLY_PAID,
  PAYMENT_TRANSACTION_STATE_PENDING,
  PAYMENT_TRANSACTION_STATE_TOO_MUCH_PAID, SATOSHIS_PER_BITCOIN,
} from 'store/ducks/payments';

export class PaymentPlanEntity {
  constructor(data) {
    this.uuid = data.uuid;
    this.name = data.name;
    this.description = data.description;
    this.bandwidth = data.bandwidth;
    this.devices = data.devices;
    this.price = data.price;
    this.duration = data.duration;

    return Object.freeze(this);
  }

  isUpgradedPlan(userPlan: PaymentPlanEntity): boolean {
    return this.price > userPlan.price;
  }

  static create(data): PaymentPlanEntity {
    return new PaymentPlanEntity(data);
  }
}

export class TransactionEntity {
  constructor(data) {
    this.uuid = data.uuid;
    this.userUuid = data.userUuid;
    this.paymentPlanUuid = data.paymentPlanUuid;
    this.paymentAddress = data.paymentAddress;
    this.conversionRate = data.conversionRate;
    this.unconfirmedReceivedAmountInSatoshis = data.unconfirmedReceivedAmountInSatoshis;
    this.confirmedFullyPaid = data.confirmedFullyPaid;
    this.receivedAmountInSatoshis = data.receivedAmountInSatoshis;
    this.expectedAmountInSatoshis = data.expectedAmountInSatoshis;
    this.state = data.state;
    this.expiresAt = data.expiresAt;
  }

  static create(data): TransactionEntity {
    return new TransactionEntity(data);
  }

  get isPending(): boolean {
    return !this.fullyPaid;
  }

  get isPartiallyPaid(): boolean {
    return this.state === PAYMENT_TRANSACTION_STATE_PARTIALLY_PAID;
  }

  get amountReceived(): number {
    return this.receivedAmountInSatoshis + this.unconfirmedReceivedAmountInSatoshis;
  }

  get remainingAmount(): number {
    return this.expectedAmountInSatoshis - this.amountReceived;
  }

  get remainingAmountAsBitcoin(): number {
    return this.remainingAmount / SATOSHIS_PER_BITCOIN;
  }

  get fullyPaid(): boolean {
    return this.state === PAYMENT_TRANSACTION_STATE_FULLY_PAID ||
      this.state === PAYMENT_TRANSACTION_STATE_TOO_MUCH_PAID;
  }
}
