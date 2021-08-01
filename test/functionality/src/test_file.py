#!/usr/bin/env python

import common
import http.client
import os
import os.path
import platform
import re
import urllib
from assertpy import assert_that
from http import HTTPStatus

TEST_DATA_DIR = "dir"
TEST_DATA_PATH = os.path.join(os.getcwd(), "test", "functionality", "data", TEST_DATA_DIR)

# sample: "2021-07-31T22:44:17.7724489+08:00"
TIMESTAMP_PATTERN = r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+\+\d{2}:\d{2}$"
# sample: "-rw-rw-rw-"
REGULAR_FILE_MODE_PATTERN = r"^-((r|-)(w|-)(x|-)){3}$"
# sample: "drwxrwxrwx"
DIRECTORY_MODE_PATTERN = r"^d((r|-)(w|-)(x|-)){3}$"

LINUX_DIR_SIZE = 4096
NON_EXIST_FILE = "dummy"

class TestFile:

    # cases of get file meta

    def test_get_file_meta_text_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_PATH, "text.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/fileMeta?" + params
            conn.request("GET", url)
            response = conn.getresponse()
            
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that("text.txt").is_equal_to(data[0]["name"])
            assert_that(23).is_equal_to(data[0]["size"])
            assert_that(data[0]["lastModifiedTime"]).matches(TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data[0]["owner"]).is_true()
            assert_that(data[0]["mode"]).matches(REGULAR_FILE_MODE_PATTERN)
        finally:
            conn.close()

    def test_get_file_meta_unicode_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_PATH, "中文.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/fileMeta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that("中文.txt").is_equal_to(data[0]["name"])
            assert_that(20).is_equal_to(data[0]["size"])
            assert_that(data[0]["lastModifiedTime"]).matches(TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data[0]["owner"]).is_true()
            assert_that(data[0]["mode"]).matches(REGULAR_FILE_MODE_PATTERN)
        finally:
            conn.close()

    def test_get_file_meta_binary_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_PATH, "binary")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/fileMeta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that("binary").is_equal_to(data[0]["name"])
            assert_that(4).is_equal_to(data[0]["size"])
            assert_that(data[0]["lastModifiedTime"]).matches(TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data[0]["owner"]).is_true()
            assert_that(data[0]["mode"]).matches(REGULAR_FILE_MODE_PATTERN)
        finally:
            conn.close()

    def test_get_file_meta_directory(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":TEST_DATA_PATH})
            url = "/api/v1/fileMeta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that(TEST_DATA_DIR).is_equal_to(data[0]["name"])
            dir_size = 0
            if "Linux" == platform.system():
                dir_size = LINUX_DIR_SIZE
            assert_that(dir_size).is_equal_to(data[0]["size"])
            assert_that(data[0]["lastModifiedTime"]).matches(TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data[0]["owner"]).is_true()
            assert_that(data[0]["mode"]).matches(DIRECTORY_MODE_PATTERN)
        finally:
            conn.close()

    def test_get_file_meta_directory_list(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":TEST_DATA_PATH, "list":True})
            url = "/api/v1/fileMeta?" + params
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

    def test_get_file_meta_non_exist(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_PATH, NON_EXIST_FILE)
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/fileMeta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()

    def test_get_file_meta_empty_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_PATH, "empty")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/fileMeta?" + params
            conn.request("GET", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that("empty").is_equal_to(data[0]["name"])
            assert_that(0).is_equal_to(data[0]["size"])
            assert_that(data[0]["lastModifiedTime"]).matches(TIMESTAMP_PATTERN)
            if "Linux" == platform.system():
                assert_that(data[0]["owner"]).is_true()
            assert_that(data[0]["mode"]).matches(REGULAR_FILE_MODE_PATTERN)
        finally:
            conn.close()

    # cases of get file content

    def test_get_file_content_text_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_PATH, "text.txt")
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

    def test_get_file_content_unicode_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_PATH, "中文.txt")
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

    def test_get_file_content_binary_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_PATH, "binary")
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

    def test_get_file_content_empty_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_PATH, "empty")
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

    def test_get_file_content_not_exist(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_PATH, NON_EXIST_FILE)
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("GET", url)
            response = conn.getresponse()
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()
