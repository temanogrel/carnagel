import { httpClient } from 'api';
import { UserEntity } from 'api/user/entity';
import { store } from 'store/configure-store';

export class UserService {
  static getCurrentUser(): Promise<UserEntity> {
    const { session } = store.getState();

    return httpClient
      .get(`backend://users/${session.user}`)
      .then(({ data }) => UserEntity.create(data));
  }

  static isAvailable(field, value): Promise<void> {
    return httpClient.post('backend://rpc/user.available', {
      field: field,
      value: value,
    });
  }

  static register(details): Promise<void> {
    return httpClient.post('backend://users', details);
  }

  static requestResetPassword(usernameOrEmail): Promise<void> {
    return httpClient.post('backend://rpc/user.password-reset', { usernameOrEmail });
  }

  static setNewPassword(password, token): Promise<void> {
    return httpClient.post('backend://rpc/user.new-password', { password, token });
  }
}
