import json

HOST = "127.0.0.1"
PORT = 5000

def assert_successful_response(response, status, data=None):
    body = response.read()
    body = json.loads(body)

    assert status == response.status
    assert status == body["status"]
    assert True == body["successful"]
    if data:
        assert data == body["data"]

def assert_failed_response(response, status, data=None):
    body = response.read()
    body = json.loads(body)

    assert status == response.status
    assert status == body["status"]
    assert False == body["successful"]
    if data:
        assert data == body["data"]
