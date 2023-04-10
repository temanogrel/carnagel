import { RecordingImageEntity } from 'api';

export class ProxyService {
  static proxies = {
    1: 'p2.camtube.co',
    2: 'p2.camtube.co',
    3: 'p3.camtube.co',
    4: 'p4.camtube.co',
    5: 'p5.camtube.co',
    6: 'p6.camtube.co',
    7: 'p7.camtube.co',
    8: 'p8.camtube.co',
    9: 'p9.camtube.co',
    0: 'p10.camtube.co'
  };

  static getImageUrl(image: RecordingImageEntity) {
    const index = parseInt(image.uuid[0]);
    switch (index) {
      case 1:
      case 2:
      case 3:
      case 4:
      case 5:
      case 6:
      case 7:
      case 8:
      case 9:
      case 0:
        return `//${ProxyService.proxies[index]}/c/${image.uuid}`;
      default:
        return `//${ProxyService.proxies[image.uuid.charCodeAt(0) % 10]}/c/${image.uuid}`;
    }
  }
}
