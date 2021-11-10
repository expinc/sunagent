import common
import http.client
import os
import os.path
import platform
import urllib
from assertpy import assert_that
from http import HTTPStatus

NON_EXIST_FILE = "dummy"

class TestFileGetMeta:

    def test_text_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "text.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file/meta?" + params
            conn.request("GET", url)
            response = conn.getresponse()
            
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that("text.txt").is_equal_to(data[0]["name"])
            assert_that(23).is_equal_to(data[0]["size"])
            assert_that(data[0]["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data[0]["owner"]).is_true()
            assert_that(data[0]["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)
        finally:
            conn.close()

    def test_unicode_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "中文.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file/meta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that("中文.txt").is_equal_to(data[0]["name"])
            assert_that(20).is_equal_to(data[0]["size"])
            assert_that(data[0]["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data[0]["owner"]).is_true()
            assert_that(data[0]["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)
        finally:
            conn.close()

    def test_binary_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "binary")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file/meta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that("binary").is_equal_to(data[0]["name"])
            assert_that(4).is_equal_to(data[0]["size"])
            assert_that(data[0]["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data[0]["owner"]).is_true()
            assert_that(data[0]["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)
        finally:
            conn.close()

    def test_directory(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":common.TEST_DATA_PATH})
            url = "/api/v1/file/meta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that(common.TEST_DATA_DIR).is_equal_to(data[0]["name"])
            dir_size = 0
            if "Linux" == platform.system():
                dir_size = common.LINUX_DIR_SIZE
            assert_that(dir_size).is_equal_to(data[0]["size"])
            assert_that(data[0]["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data[0]["owner"]).is_true()
            assert_that(data[0]["mode"]).matches(common.DIRECTORY_MODE_PATTERN)
        finally:
            conn.close()

    def test_directory_list(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":common.TEST_DATA_PATH, "list":True})
            url = "/api/v1/file/meta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(5).is_equal_to(len(data))
            files = [f["name"] for f in data]
            expected_files = ["binary", "empty", "shell.sh", "text.txt", "中文.txt"]
            for expected_file in expected_files:
                assert_that(files).contains(expected_file)
        finally:
            conn.close()

    def test_non_exist(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, NON_EXIST_FILE)
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file/meta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()

    def test_empty_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_DATA_PATH, "empty")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file/meta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that("empty").is_equal_to(data[0]["name"])
            assert_that(0).is_equal_to(data[0]["size"])
            assert_that(data[0]["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data[0]["owner"]).is_true()
            assert_that(data[0]["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)
        finally:
            conn.close()
