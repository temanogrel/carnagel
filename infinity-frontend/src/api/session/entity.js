export class SessionEntity {
  constructor(data) {
    this.rawToken = data.rawToken;
    this.role = data.role;
    this.user = data.userUuid;
    this.session = data.session;
    this.expiresAt = data.exp;
    this.paymentPlan = data.paymentPlan;
    this.blackListedToday = data.blackListedToday;

    return Object.freeze(this);
  }

  get isExpired(): boolean {
    return this.expiresAt < new Date().getTime() / 1000;
  }

  get expiresIn(): number {
    return this.expiresAt - new Date().getTime() / 1000;
  }

  static createFromJwtToken(token: string): SessionEntity {
    let parts = token.split('.');
    if (parts.length !== 3) {
      throw new Error(`Invalid JWT token "${token}" provided`);
    }

    const data = JSON.parse(atob(parts[1]));
    data.rawToken = token;

    return SessionEntity.create(data);
  }

  static create(data): SessionEntity {
    return new SessionEntity(data);
  }
}
