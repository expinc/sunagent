import common
import http.client
import os
import shutil
import urllib
from assertpy import assert_that
from http import HTTPStatus

class TestFileAppend:

    @classmethod
    def setup_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR, ignore_errors=True)
        os.makedirs(common.TEST_TMP_DIR, exist_ok=True)

    @classmethod
    def teardown_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR)

    def test_ordinary(self):
        conn = http.client.HTTPConnection(common.HOST, common.PORT)
        try:
            # read original content
            file_path = os.path.join(common.TEST_TMP_DIR, "binary")
            shutil.copyfile(os.path.join(common.TEST_DATA_PATH, "binary"), file_path)            
            with open(file_path, "rb") as f:
                originContent = f.read()

            # send request
            file_path = os.path.join(common.TEST_TMP_DIR, "binary")
            params = urllib.parse.urlencode({"path":file_path})
            url = "/api/v1/file/append?" + params
            append_content = b'\x22\xbb'
            conn.request("POST", url, append_content, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)

            # verify response
            assert_that("binary").is_equal_to(data["name"])
            assert_that(4+2).is_equal_to(data["size"])
            assert_that(data["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            assert_that(data["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)

            # verify file content
            with open(file_path, "rb") as f:
                content = f.read()
            assert_that(content).is_equal_to(originContent + append_content)
        finally:
            os.remove(file_path)
            conn.close()

    def test_nonexist_file(self):
        conn = http.client.HTTPConnection(common.HOST, common.PORT)
        try:
            # specify file path
            file_path = os.path.join(common.TEST_TMP_DIR, "binary")

            # send request
            params = urllib.parse.urlencode({"path":file_path})
            url = "/api/v1/file/append?" + params
            append_content = b'\x22\xbb'
            conn.request("POST", url, append_content, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            data = common.get_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()

    def test_empty_file(self):
        conn = http.client.HTTPConnection(common.HOST, common.PORT)
        try:
            # read original content
            file_path = os.path.join(common.TEST_TMP_DIR, "empty")
            shutil.copyfile(os.path.join(common.TEST_DATA_PATH, "empty"), file_path)            
            with open(file_path, "rb") as f:
                originContent = f.read()

            # send request
            file_path = os.path.join(common.TEST_TMP_DIR, "empty")
            params = urllib.parse.urlencode({"path":file_path})
            url = "/api/v1/file/append?" + params
            append_content = b'\x22\xbb'
            conn.request("POST", url, append_content, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)

            # verify response
            assert_that("empty").is_equal_to(data["name"])
            assert_that(2).is_equal_to(data["size"])
            assert_that(data["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            assert_that(data["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)

            # verify file content
            with open(file_path, "rb") as f:
                content = f.read()
            assert_that(content).is_equal_to(originContent + append_content)
        finally:
            os.remove(file_path)
            conn.close()

    def test_append_nothing(self):
        conn = http.client.HTTPConnection(common.HOST, common.PORT)
        try:
            # read original content
            file_path = os.path.join(common.TEST_TMP_DIR, "binary")
            shutil.copyfile(os.path.join(common.TEST_DATA_PATH, "binary"), file_path)
            with open(file_path, "rb") as f:
                originContent = f.read()

            # send request
            file_path = os.path.join(common.TEST_TMP_DIR, "binary")
            params = urllib.parse.urlencode({"path":file_path})
            url = "/api/v1/file/append?" + params
            append_content = b''
            conn.request("POST", url, append_content, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)

            # verify response
            assert_that("binary").is_equal_to(data["name"])
            assert_that(4).is_equal_to(data["size"])
            assert_that(data["lastModifiedTime"]).matches(common.TIMESTAMP_PATTERN)
            assert_that(data["mode"]).matches(common.REGULAR_FILE_MODE_PATTERN)

            # verify file content
            with open(file_path, "rb") as f:
                content = f.read()
            assert_that(content).is_equal_to(originContent + append_content)
        finally:
            os.remove(file_path)
            conn.close()

    def test_oversize(self):
        conn = http.client.HTTPConnection(common.HOST, common.PORT)
        try:
            # prepare original file
            file_path = os.path.join(common.TEST_TMP_DIR, "binary")
            shutil.copyfile(os.path.join(common.TEST_DATA_PATH, "binary"), file_path)

            # send request
            params = urllib.parse.urlencode({"path":file_path, "isDir":False})
            url = "/api/v1/file/append?" + params
            content = bytes(101 * 1024 * 1024)
            conn.request("POST", url, content, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify response
            common.assert_failed_response(response, HTTPStatus.BAD_REQUEST)
        finally:
            conn.close()
