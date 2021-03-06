import common
import http.client
import os
import platform
import shutil
import urllib.parse
from assertpy import assert_that
from http import HTTPStatus

class TestFileOverwrite:

    @classmethod
    def setup_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR, ignore_errors=True)
        os.makedirs(common.TEST_TMP_DIR, exist_ok=True)

    @classmethod
    def teardown_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR)

    def test_new_text_file(self):
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
            conn.request("PUT", url, originContent, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
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

    def test_exist_text_file(self):
        try:
            # prepare original file
            originPath = os.path.join(common.TEST_DATA_PATH, "text.txt")
            shutil.copy(originPath, common.TEST_TMP_DIR)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_TMP_DIR, "text.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            new_content = b"this is new content\n"
            conn.request("PUT", url, new_content, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)

            # verify response
            assert_that("text.txt").is_equal_to(data["name"])
            assert_that(data["size"]).is_equal_to(20)
            assert_that(data["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data["owner"]).is_true()
            assert_that(data["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)

            # verify file content
            with open(path, "rb") as f:
                content = f.read()
            assert_that(content).is_equal_to(new_content)
        finally:
            conn.close()

    def test_new_directory(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_TMP_DIR, "dir", "subdir")
            params = urllib.parse.urlencode({"path":path, "isDir":True})
            url = "/api/v1/file?" + params
            conn.request("PUT", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
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

    def test_exist_directory(self):
        try:
            # prepare directory
            dirPath = os.path.join(common.TEST_TMP_DIR, "dir", "subdir")
            os.makedirs(common.TEST_TMP_DIR, exist_ok=True)
            filePath = os.path.join(common.TEST_DATA_PATH, "text.txt")
            shutil.copy(filePath, dirPath)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_TMP_DIR, "dir", "subdir")
            params = urllib.parse.urlencode({"path":path, "isDir":True})
            url = "/api/v1/file?" + params
            conn.request("PUT", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
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

            # verify file
            assert_that(os.path.isfile(filePath)).is_true()
        finally:
            conn.close()

    def test_oversize(self):
        try:
            # prepare directory
            dirPath = os.path.join(common.TEST_TMP_DIR, "dir", "subdir")
            os.makedirs(common.TEST_TMP_DIR, exist_ok=True)
            filePath = os.path.join(common.TEST_DATA_PATH, "text.txt")
            shutil.copy(filePath, dirPath)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":filePath, "isDir":False})
            url = "/api/v1/file?" + params
            content = bytes(101 * 1024 * 1024)
            conn.request("PUT", url, content, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify response
            common.assert_failed_response(response, HTTPStatus.BAD_REQUEST)
        finally:
            conn.close()
