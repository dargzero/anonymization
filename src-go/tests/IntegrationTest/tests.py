import http.client
import json
from pymongo import MongoClient


dataSetName = "test_set__"
k = 3
n = 13

settings = json.dumps({
    "settings": {
        "k": k,
        "algorithm": "mondrian",
        "mode": "continuous"
    },
    "fields": [
        {
            "name": "location",
            "mode": "qid",
            "type": "coords"
        },
        {
            "name": "id",
            "mode": "drop",
            "type": ""
        }
    ]
})
data = json.dumps([
    {
        "location": "-136°, 69°",
        "id": "0"
    },
    {
        "location": "-145°, 19°",
        "id": "1"
    },
    {
        "location": "125°, 19°",
        "id": "2"
    },
    {
        "location": "-47°, 85°",
        "id": "3"
    },
    {
        "location": "122°, -114°",
        "id": "4"
    },
    {
        "location": "112°, 170°",
        "id": "5"
    },
    {
        "location": "-130°, 34°",
        "id": "6"
    },
    {
        "location": "52°, -174°",
        "id": "7"
    },
    {
        "location": "38°, 27°",
        "id": "8"
    },
    {
        "location": "141°, 77°",
        "id": "9"
    },
    {
        "location": "41°, 114°",
        "id": "10"
    },
    {
        "location": "-102°, 106°",
        "id": "11"
    },
    {
        "location": "-95°, -67°",
        "id": "12"
    }
])

def get_connect():
    return http.client.HTTPConnection('anonymization_server:9137')

headers = {'Content-type': 'application/json', 'Accept':'application/json'}


def clear_server():
    connection = get_connect()
    connection.request('DELETE', '/v1/datasets/test_set__', None, headers)
    return connection.getresponse()


def get_dataset():
    connection = get_connect()
    connection.request('GET', '/v1/datasets', None, headers)
    return connection.getresponse()


def create_dataset():
    connection = get_connect()

    connection.request('PUT', '/v1/datasets/test_set__', settings, headers)
    x = connection.getresponse()
    connection.close()
    return x

def create_session():
    connection = get_connect()
    connection.request('POST', '/v1/upload', json.dumps({"datasetName": "test_set__"}), headers)
    res = connection.getresponse()
    stri = str(res.read())
    connection.close()
    return stri[16:-3], res

def upload_data(sid):

    connection = get_connect()
    connection.request('POST', '/v1/upload/'+sid, data, headers)
    x = connection.getresponse()
    connection.close()
    return x

def test_db():
    client = MongoClient('anonymization_database', 27017)
    db = client['anondb']
    pipeline = [{"$match":{"__anonymized": True}},
                {"$group": {"_id": "$location", "count": {"$sum": 1}}},
                {"$sort":{"count": 1}},
                {"$limit": 1}]

    collection = db['anon_' + dataSetName]
    val = db.command('aggregate', 'anon_'+dataSetName, pipeline=pipeline, explain=False)
    if val['cursor']['firstBatch'][0]['count']<k:
       return  False
    if collection.find().count()!= n:
        return  False
    return True

def test_db_form():
    client = MongoClient('anonymization_database', 27017)
    db = client['anondb']
    pipeline = [{"$match":{"__anonymized": True,
                "location":{ "$regex":"(-?\d+(.\d+)?):(-?\d+(.\d+)?), (-?\d+(.\d+)?):(-?\d+(.\d+)?)"}}},
                {"$group" : {"_id": "null","count": { "$sum": 1 }}}
                 ]
    val = db.command('aggregate', 'anon_'+dataSetName, pipeline=pipeline, explain=False)
    if val['cursor']['firstBatch'][0]['count']!=n:
        return False
    collection = db['anon_' + dataSetName]
    if collection.find().count()!= n:
        return  False
    return True