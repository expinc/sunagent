import common
import http.client
import os
import platform
import shutil
import urllib.parse
from assertpy import assert_that
from http import HTTPStatus

class TestFileCreate:

    @classmethod
    def setup_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR, ignore_errors=True)
        os.makedirs(common.TEST_TMP_DIR, exist_ok=True)

    @classmethod
    def teardown_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR)

    def test_text_file(self):
        try:
            # read file content
            originPath = os.path.join(common.TEST_DATA_PATH, "text.txt")
            with open(originPath, "rb") as f:
                originContent = f.read()

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_TMP_DIR, "text.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("POST", url, originContent)
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)

            # verify response
            assert_that("text.txt").is_equal_to(data["name"])
            assert_that(23).is_equal_to(data["size"])
            assert_that(data["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data["owner"]).is_true()
            assert_that(data["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)

            # verify file content
            with open(path, "rb") as f:
                content = f.read()
            assert_that(content).is_equal_to(originContent)
        finally:
            conn.close()

    def test_unicode_file(self):
        try:
            # read file content
            originPath = os.path.join(common.TEST_DATA_PATH, "中文.txt")
            with open(originPath, "rb") as f:
                originContent = f.read()

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_TMP_DIR, "中文.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("POST", url, originContent)
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)

            # verify response
            assert_that("中文.txt").is_equal_to(data["name"])
            assert_that(20).is_equal_to(data["size"])
            assert_that(data["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data["owner"]).is_true()
            assert_that(data["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)

            # verify file content
            with open(path, "rb") as f:
                content = f.read()
            assert_that(content).is_equal_to(originContent)
        finally:
            conn.close()

    def test_binary_file(self):
        try:
            # read file content
            originPath = os.path.join(common.TEST_DATA_PATH, "中文.txt")
            with open(originPath, "rb") as f:
                originContent = f.read()

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_TMP_DIR, "binary")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("POST", url, originContent)
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)

            # verify response
            assert_that("binary").is_equal_to(data["name"])
            assert_that(20).is_equal_to(data["size"])
            assert_that(data["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data["owner"]).is_true()
            assert_that(data["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)

            # verify file content
            with open(path, "rb") as f:
                content = f.read()
            assert_that(content).is_equal_to(originContent)
        finally:
            conn.close()

    def test_directory(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_TMP_DIR, "dir", "subdir")
            params = urllib.parse.urlencode({"path":path, "isDir":True})
            url = "/api/v1/file?" + params
            conn.request("POST", url)
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)

            # verify response
            assert_that("subdir").is_equal_to(data["name"])
            dir_size = 0
            if "Linux" == platform.system():
                dir_size = common.LINUX_DIR_SIZE
            assert_that(dir_size).is_equal_to(data["size"])
            assert_that(data["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data["owner"]).is_true()
            assert_that(data["mode"]).matches(common.DIRECTORY_MODE_PATTERN)
        finally:
            conn.close()

    def test_empty_file(self):
        try:
            # prepare file content
            originContent = bytearray()

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_TMP_DIR, "empty")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("POST", url, originContent)
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)

            # verify response
            assert_that("empty").is_equal_to(data["name"])
            assert_that(0).is_equal_to(data["size"])
            assert_that(data["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data["owner"]).is_true()
            assert_that(data["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)

            # verify file content
            with open(path, "rb") as f:
                content = f.read()
            assert_that(content).is_equal_to(originContent)
        finally:
            conn.close()

    def test_exist_file(self):
        try:
            # prepare file content
            path = os.path.join(common.TEST_TMP_DIR, "exist.txt")
            with open(path, "w") as f:
                f.write("exist")

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("POST", url, "")
            response = conn.getresponse()

            # verify response
            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()

    def test_exist_dir(self):
        try:
            # prepare file content
            path = os.path.join(common.TEST_TMP_DIR, "exist-dir")
            os.makedirs(path)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":path, "isDir":True})
            url = "/api/v1/file?" + params
            conn.request("POST", url)
            response = conn.getresponse()

            # verify response
            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()
