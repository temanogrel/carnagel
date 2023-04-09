from central.api.request_handlers.core import JsonRequestHandler


class CredentialCollectionHandler(JsonRequestHandler):
    def post(self, service: str):
        """
        Update all the existing credentials

        :param service:

        :return:
        """

        service = self.service_container.get_service(service)

        if not hasattr(service, 'set_credential'):
            return self.send_error(400, message='This services does not support credentials')

        for identity, api_token in self.request.json.items():
            service.set_credential(identity, api_token)

    def get(self, service: str):
        """
        Return a list of all the credentials

        :param service:

        :return:
        """

        service = self.service_container.get_service(service)

        if not hasattr(service, 'get_credentials'):
            return self.send_error(400, message='This service does not support credentials')

        self.write(service.get_credentials())


class SessionIdResourceHandler(JsonRequestHandler):
    """
    MyFreeCams requires a session id to allow downloading
    """

    def get(self, service: str):
        """
        Retrieve the session id

        If no session id is available then we should return a 400 error

        :param service:
        :return:
        """
        service = self.get_service(service)

        if not hasattr(service, 'session_id'):
            return self.send_error(400, message='Service does not support session id\'s')

        self.write(dict(session_id=service.session_id))

    def post(self, service: str):
        """
        Update the session id for my freecams

        :param service:
        :return:
        """
        service = self.get_service(service)

        if not hasattr(service, 'session_id'):
            return self.send_error(400, message='Service does not support session id\'s')

        setattr(service, '_session_id', self.request.json['session_id'])

        self.set_status(204)


class ServiceCollectionHandler(JsonRequestHandler):
    """
    Get meta data related to the services
    """

    def get(self, *args, **kwargs):
        self.write(dict(self.service_container.get_meta()))
