import common
import http.client
import os
import os.path
import shutil
import urllib
from assertpy import assert_that
from http import HTTPStatus

LINUX_DIR_SIZE = 4096
NON_EXIST_FILE = "dummy"

class TestFileGetContent:

    @classmethod
    def setup_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR, ignore_errors=True)
        os.makedirs(common.TEST_TMP_DIR, exist_ok=True)

    @classmethod
    def teardown_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR)

    def test_text_file(self):
        try:
            # Send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "text.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # Check response body
            data = common.get_binary_response(response, HTTPStatus.OK)
            with open(path, "rb") as f:
                content = f.read()
                assert_that(data).is_equal_to(content)

            # Check response header
            content_length = int(response.getheader("Content-Length"))
            assert_that(content_length).is_equal_to(os.stat(path).st_size)
        finally:
            conn.close()

    def test_unicode_file(self):
        try:
            # Send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "中文.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # Check response body
            data = common.get_binary_response(response, HTTPStatus.OK)
            with open(path, "rb") as f:
                content = f.read()
                assert_that(data).is_equal_to(content)

            # Check response header
            content_length = int(response.getheader("Content-Length"))
            assert_that(content_length).is_equal_to(os.stat(path).st_size)
        finally:
            conn.close()

    def test_binary_file(self):
        try:
            # Send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "binary")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # Check response body
            data = common.get_binary_response(response, HTTPStatus.OK)
            with open(path, "rb") as f:
                content = f.read()
                assert_that(data).is_equal_to(content)

            # Check response header
            content_length = int(response.getheader("Content-Length"))
            assert_that(content_length).is_equal_to(os.stat(path).st_size)
        finally:
            conn.close()

    def test_empty_file(self):
        try:
            # Send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "empty")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # Check response body
            data = common.get_binary_response(response, HTTPStatus.OK)
            with open(path, "rb") as f:
                content = f.read()
                assert_that(data).is_equal_to(content)

            # Check response header
            content_length = int(response.getheader("Content-Length"))
            assert_that(content_length).is_equal_to(os.stat(path).st_size)
        finally:
            conn.close()

    def test_not_exist(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, NON_EXIST_FILE)
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()

    def test_file_stream(self):
        try:
            # Prepare file
            path = os.path.join(common.TEST_TMP_DIR, "mybinary")
            # Default streaming chunk size is 100 MB (see config CORE.FileUploadMaxBytes)
            # Make file length as 250 MB to make the file transferred by multiple streaming chunks (100 MB, 100 MB, 50 MB)
            file_length = 262144000
            with open(path, "wb") as f:
                f.write(b"\xff" * file_length)

            # Send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # Check response body
            data = common.get_binary_response(response, HTTPStatus.OK)
            with open(path, "rb") as f:
                content = f.read()
                assert_that(data).is_equal_to(content)

            # Check response header
            content_length = int(response.getheader("Content-Length"))
            assert_that(content_length).is_equal_to(os.stat(path).st_size)
        finally:
            conn.close()
