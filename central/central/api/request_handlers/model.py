from central.api.request_handlers.core import JsonRequestHandler


class ModelIntersectionHandler(JsonRequestHandler):

    def post(self, service: str):
        self.get_service(service).intersect_performers(self.request.json)


class ModelCollectionHandler(JsonRequestHandler):

    def get(self, service: str):
        """
        Receive a list of performers
        """

        data = {}

        for uid, performer in self.get_service(service).performers.items():
            data[uid] = performer.__dict__

        self.write(data)


class ModelResourceHandler(JsonRequestHandler):

    def patch(self, service: str, uid: str):
        """
        Update a performer

        :param service: str
        :param uid: str

        :return:
        """

        service = self.get_service(service)

        if uid not in service:
            return self.send_error(404, message='Model not found')

        service.update_performer(service[uid], self.request.json)

        # No point in returning the data
        self.set_status(204)

    def delete(self, service: str, uid: str):
        """
        Remove a performer from the modelserver

        :param service:
        :param uid:

        :return:
        """

        service = self.get_service(service)

        if uid not in service:
            return self.write_error(404, message='Model not found')

        del service[uid]

        self.set_status(204)


