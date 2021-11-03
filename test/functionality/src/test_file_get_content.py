import common
import http.client
import os
import os.path
import urllib
from assertpy import assert_that
from http import HTTPStatus

LINUX_DIR_SIZE = 4096
NON_EXIST_FILE = "dummy"

class TestFileGetContent:

    def test_text_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "text.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_binary_response(response, HTTPStatus.OK)
            with open(path, "rb") as f:
                content = f.read()
                assert_that(data).is_equal_to(content)
        finally:
            conn.close()

    def test_unicode_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "中文.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_binary_response(response, HTTPStatus.OK)
            with open(path, "rb") as f:
                content = f.read()
                assert_that(data).is_equal_to(content)
        finally:
            conn.close()

    def test_binary_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "binary")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_binary_response(response, HTTPStatus.OK)
            with open(path, "rb") as f:
                content = f.read()
                assert_that(data).is_equal_to(content)
        finally:
            conn.close()

    def test_empty_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "empty")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_binary_response(response, HTTPStatus.OK)
            with open(path, "rb") as f:
                content = f.read()
                assert_that(data).is_equal_to(content)
        finally:
            conn.close()

    def test_not_exist(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, NON_EXIST_FILE)
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url)
            response = conn.getresponse()
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()
