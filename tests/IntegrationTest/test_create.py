import unittest

import common


class CreateDatasetTest(common.CommonTest):

    def setUp(self):
        self.set_dataset_name('create_dataset_test')
        super().setUp()

    def test_create(self):
        create_dataset_request = {
            'settings': {
                'k': 10,
                'algorithm': 'mondrian',
                'mode': 'single'
            },
            'fields': [
                {
                    'name': 'company',
                    'mode': 'id'
                },
                {
                    'name': 'name',
                    'mode': 'id'
                },
                {
                    'name': 'serial',
                    'mode': 'id'
                },
                {
                    'name': 'stackTrace',
                    'mode': 'drop'
                }
            ]
        }
        self.create_dataset(create_dataset_request)

        status, response = self.json_request('GET', '/v1/datasets/' + self.dataset_name)
        self.assertEqual(status, 200)
        self.assertEqual(response['name'], self.dataset_name)


if __name__ == '__main__':
    unittest.main()
