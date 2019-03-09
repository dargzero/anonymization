import http.client
import json
import unittest


class CommonTest(unittest.TestCase):

    def set_dataset_name(self, dataset_name):
        self.dataset_name = dataset_name

    def setUp(self):
        self.connection = http.client.HTTPConnection('anonymization_server', 9137)

    def tearDown(self):
        self.empty_request('DELETE', '/v1/datasets/' + self.dataset_name)
        self.connection.close()

    def request(self, method, url, body=None):
        json_body_string = json.dumps(body)
        self.connection.request(method, url, json_body_string.encode('utf8'), {'Content-Type': 'application/json'})

        return self.connection.getresponse()

    def empty_request(self, method, url, body=None):
        response = self.request(method, url, body)
        response.read()
        return response.status

    def json_request(self, method, url, body=None):
        response = self.request(method, url, body)
        return response.status, json.load(response)

    def create_dataset(self, create_dataset_request):
        status = self.empty_request('PUT', '/v1/datasets/' + self.dataset_name, create_dataset_request)
        self.assertEqual(status, 201)
