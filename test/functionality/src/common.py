import json
import urllib.parse
from assertpy import assert_that

HOST = "127.0.0.1"
PORT = 5000

def assert_successful_response(response, status, data=None):
    body = response.read()
    body = json.loads(body)

    assert_that(status).is_equal_to(response.status)
    assert_that(status).is_equal_to(body["status"])
    assert_that(True).is_equal_to(body["successful"])
    if data:
        assert_that(data).is_equal_to(body["data"])

def assert_failed_response(response, status, data=None):
    body = response.read()
    body = json.loads(body)

    assert_that(status).is_equal_to(response.status)
    assert_that(status).is_equal_to(body["status"])
    assert_that(False).is_equal_to(body["successful"])
    if data:
        assert_that(data).is_equal_to(body["data"])

def get_successful_response(response, status):
    body = response.read()
    body = json.loads(body)
    assert_that(status).is_equal_to(response.status)
    return body["data"]
