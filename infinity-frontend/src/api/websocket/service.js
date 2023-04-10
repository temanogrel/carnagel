import apiEnvironment from 'api/environment';
import { Observable } from 'rxjs/Observable';
import { store } from 'store/configure-store';
import { websocketBroadcast, websocketConnectionChanged } from 'store/ducks/websocket';

export class WebsocketService {
  rpcCounter = 0;
  requests: { [id: number]: { request: { name: string, payload: Object}, resolve: Function, reject: Function } } = {};

  constructor() {
    this.createConnection();
  }

  createConnection = () => {
    this.socket = Observable.webSocket({
      url: apiEnvironment.websocketUri,
      openObserver: {
        next: () => store.dispatch(websocketConnectionChanged(true)),
      },
      closeObserver: {
        next: () => store.dispatch(websocketConnectionChanged(false)),
      },
    });

    this
      .socket
      .subscribe(this.onMessage, () => setTimeout(this.createConnection, 500));
  };

  sendRpc(name: string, payload: Object): Promise<Object> {
    const id = ++this.rpcCounter;

    this.socket.next({ id, name, payload });

    return new Promise((resolve, reject) => {
      this.requests[id] = { request: { name, payload }, resolve, reject };
    });
  }

  onMessage = (msg: Object) => {
    // Check if rpc
    if ('id' in msg) {
      if (!this.requests[msg.id]) {
        return console.warn('Received response for unknown rpc request', msg);
      }

      const { resolve, reject } = this.requests[msg.id];

      if (msg.success) {
        resolve(msg.payload);
      } else {
        reject(msg.payload);
      }

      return delete this.requests[msg.id];
    }

    store.dispatch(websocketBroadcast(msg.name, msg.payload));
  };
}
