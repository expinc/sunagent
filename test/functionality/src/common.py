import json
import os
from assertpy import assert_that

HOST = "127.0.0.1"
PORT = 5000

# where test data is read
TEST_DATA_DIR = "dir"
TEST_DATA_PATH = os.path.join(os.getcwd(), "test", "functionality", "data", TEST_DATA_DIR)

# where temporary file directory
TEST_TMP_DIR = os.path.join(os.getcwd(), "tmp")

# sample: "2021-07-31T22:44:17.7724489+08:00"
TIMESTAMP_PATTERN = r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+\+\d{2}:\d{2}$"
# sample: "-rw-rw-rw-"
REGULAR_FILE_MODE_PATTERN = r"^-((r|-)(w|-)(x|-)){3}$"
# sample: "drwxrwxrwx"
DIRECTORY_MODE_PATTERN = r"^d((r|-)(w|-)(x|-)){3}$"

LINUX_DIR_SIZE = 4096

def assert_successful_response(response, status, data=None):
    assert_that(status).is_equal_to(response.status)

    body = response.read()
    body = json.loads(body)

    assert_that(status).is_equal_to(body["status"])
    assert_that(True).is_equal_to(body["successful"])
    if data:
        assert_that(data).is_equal_to(body["data"])

def assert_failed_response(response, status, data=None):
    assert_that(status).is_equal_to(response.status)

    body = response.read()
    body = json.loads(body)

    assert_that(status).is_equal_to(body["status"])
    assert_that(False).is_equal_to(body["successful"])
    if data:
        assert_that(data).is_equal_to(body["data"])

def get_successful_response(response, status):
    assert_that(status).is_equal_to(response.status)
    body = response.read()
    body = json.loads(body)
    return body["data"]

def get_binary_response(response, status):
    assert_that(status).is_equal_to(response.status)
    return response.read()
