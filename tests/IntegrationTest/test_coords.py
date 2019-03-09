import http.client
import json
import tests
import http.client
import json
import unittest


class CoordsTest(unittest.TestCase):

    def test(self):
        tests.clear_server()
        self.assertEqual(tests.create_dataset().getcode(),201)
        session_id, res = tests.create_session()
        self.assertEqual(res.getcode(), 200)
        self.assertEqual(tests.upload_data(session_id).getcode(), 200)
        self.assertEqual(tests.test_db(), True)
        self.assertEqual(tests.test_db_form(), True)


if __name__ == '__main__':
    unittest.main()
