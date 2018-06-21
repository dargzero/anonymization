import bson.code
import json
import pprint
import pymongo
import statistics
import time
import unittest

import common


class UploadDataTest(common.CommonTest):

    k = 2

    upload_count = 1
    batch_size = 5
    batch_count = 1

    create_dataset_request = {
        'settings': {
            'k': k,
            'algorithm': 'mondrian',
            'mode': 'single' if upload_count == 1 else 'continuous'
        },
        'fields': [
            {
                'name': 'company',
                'mode': 'id'
            },
            {
                'name': 'cpuName',
                'mode': 'drop'
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
            },
            {
                'name': 'cpuSpeed',
                'mode': 'qid',
                'type': 'numeric'
            },
            {
                'name': 'freeMemory',
                'mode': 'qid',
                'type': 'numeric'
            },
            {
                'name': 'visibleMemory',
                'mode': 'qid',
                'type': 'numeric'
            },
            {
                'name': 'os',
                'mode': 'qid',
                'type': 'prefix'
            },
            {
                'name': 'message',
                'mode': 'qid',
                'type': 'prefix'
            }
        ]
    }

    @staticmethod
    def read_data(file):
        data = []
        for i in range(0, UploadDataTest.batch_size):
            line = file.readline()
            data.append(json.loads(line))
        return data

    @staticmethod
    def read_full_file(file):
        with open(file, encoding='utf-8') as f:
            return f.read()

    @staticmethod
    def print_statistics(name, min, max, avg):
        print('{0} - Min / Max / Avg: {1} / {2} / {3}'.format(name, min, max, avg))

    def setUp(self):
        self.set_dataset_name('exceptions')
        super().setUp()
        self.mongo_client = pymongo.MongoClient('anonymization_database', 27017)

    def tearDown(self):
        super().tearDown()
        self.mongo_client.close()

    def upload_data(self):
        with open('data/exceptions.json', encoding='utf-8') as f:
            for j in range(0, self.upload_count):
                status, response = self.json_request('POST', '/v1/upload', {'datasetName': self.dataset_name})
                self.assertEqual(status, 200)
                session_id = response['sessionId']

                for i in range(0, self.batch_count):
                    last = i == self.batch_count - 1
                    status, response = self.json_request('POST',
                                                         '/v1/upload/{0}?last={1}'.format(session_id, last),
                                                         self.read_data(f))
                    # pprint.pprint(response)
                    self.assertEqual(status, 200)
                    self.assertEqual(response['insertSuccessful'], True)
                    if last:
                        self.assertEqual(response['finalizeSuccessful'], True)

                not_anonymized = self.mongo_client.anondb.anon_exceptions.count({'__anonymized': False})
                self.assertEqual(not_anonymized, 0)

    def get_prefix_statistics(self, collection, name):
        pipeline = [
            {
                '$project': {
                    'strLen': {
                        '$cond': {
                            'if': {'$eq': ['$' + name, '-']},
                            'then': 0,
                            'else': {'$strLenCP': '$' + name}
                        }
                    }
                }
            },
            {
                '$group': {
                    '_id': None,
                    'min': {'$min': '$strLen'},
                    'max': {'$max': '$strLen'},
                    'avg': {'$avg': '$strLen'},
                }
            }
        ]
        db_response = list(collection.aggregate(pipeline))
        stats = db_response[0]
        self.print_statistics(name, stats['min'], stats['max'], stats['avg'])

    def get_numeric_statistics(self, collection, name):
        mapper = bson.code.Code(self.read_full_file('map.js').format(name))
        reducer = bson.code.Code(self.read_full_file('reduce.js'))
        finalizer = bson.code.Code(self.read_full_file('finalize.js'))

        db_response = collection.map_reduce(mapper, reducer, 'result', finalize=finalizer, jsMode=True)
        stats = db_response.find()[0]['value']
        self.print_statistics(name, stats['min'], stats['max'], stats['avg'])

    def calculate_statistics(self, collection, do_assert):
        qids = []
        for field in self.create_dataset_request['fields']:
            if field['mode'] == 'qid':
                qids.append(field)

        _id = {}
        for qid in qids:
            _id[qid['name']] = '$' + qid['name']
        pipeline = [
            {
                '$group': {
                    '_id': _id,
                    'count': {'$sum': 1}
                }
            }
        ]
        db_response = list(collection.aggregate(pipeline))
        sizes = []
        for group in db_response:
            if do_assert:
                self.assertTrue(group['count'] >= self.k, "The group {0} has count {1} that is lower than k."
                                .format(group['_id'], group['count']))
            sizes.append(group['count'])
        self.print_statistics('Count', min(sizes), max(sizes), statistics.mean(sizes))

        for qid in qids:
            if qid['type'] == 'numeric':
                self.get_numeric_statistics(collection, qid['name'])
            elif qid['type'] == 'prefix':
                self.get_prefix_statistics(collection, qid['name'])

    def test_upload(self):
        self.create_dataset(self.create_dataset_request)

        anon_start = time.perf_counter()
        self.upload_data()
        print('Anonymization took {0} s'.format(time.perf_counter() - anon_start))

        print('Original values:')
        self.calculate_statistics(self.mongo_client.anondb.data_exceptions, False)

        print('Anonymized values:')
        self.calculate_statistics(self.mongo_client.anondb.anon_exceptions, True)


if __name__ == '__main__':
    unittest.main()
