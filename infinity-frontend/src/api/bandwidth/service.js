import { BandwidthStatusEntity, websocketService } from 'api';

export class BandwidthService {
  static getRemainingBandwidth(): Promise<BandwidthStatusEntity> {
    return websocketService
      .sendRpc('bandwidth:get-remaining')
      .then(BandwidthStatusEntity.create);
  }
}
