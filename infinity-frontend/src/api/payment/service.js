import { httpClient } from 'api';
import { PaymentPlanEntity, PurchaseEntity, TransactionEntity } from 'api/payment/entity';

export class PaymentService {
  static getPaymentPlans(): Promise<PaymentPlanEntity[]> {
    return httpClient
      .get(`backend://payment-plans`)
      .then(({ data }) => data.data.map(PaymentPlanEntity.create));
  }

  static purchasePaymentPlan(uuid: string): Promise<TransactionEntity> {
    return httpClient
      .post(`backend://rpc/payments-plans.purchase`, { uuid })
      .then(({ data }) => TransactionEntity.create(data));
  }

  static getTransaction(uuid: string): Promise<TransactionEntity> {
    return httpClient
      .get(`backend://payment-transactions/${uuid}`)
      .then(({ data }) => TransactionEntity.create(data));
  }
}
