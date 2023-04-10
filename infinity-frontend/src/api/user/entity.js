export class UserEntity{
  constructor(data): UserEntity {
    this.uuid = data.uuid;
    this.email = data.email;
    this.username = data.username;
    this.paymentPlanUuid = data.paymentPlanUuid;
    this.paymentPlanSubscribedAt = data.paymentPlanSubscribedAt;
    this.createdAt = data.createdAt;
    this.updatedAt = data.updatedAt;

    return Object.freeze(this);
  }

  static create(data): UserEntity {
    return new UserEntity(data);
  }
}
