/**
 *
 *
 */

export function originServiceToString(originService: string) {
  switch (originService) {
    case 'cbc':
      return 'Chaturbate';
    case 'mfc':
      return 'MyFreeCams';
  }
}
