#!/usr/bin/env python

import common
import http.client
import os
import shutil
import urllib.parse
from assertpy import assert_that
from http import HTTPStatus

class TestFileDelete:

    @classmethod
    def setup_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR, ignore_errors=True)
        os.makedirs(common.TEST_TMP_DIR, exist_ok=True)

    @classmethod
    def teardown_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR)

    def test_text_file(self):
        try:
            # prepare original file
            originPath = os.path.join(common.TEST_DATA_PATH, "text.txt")
            shutil.copy(originPath, common.TEST_TMP_DIR)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_TMP_DIR, "text.txt")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("DELETE", url)
            response = conn.getresponse()
            common.assert_successful_response(response, HTTPStatus.OK)

            # verify file deleted
            assert_that(os.path.exists(path)).is_false()
        finally:
            conn.close()

    def test_empty_directory(self):
        try:
            # prepare directory
            dirPath = os.path.join(common.TEST_TMP_DIR, "dir", "subdir")
            os.makedirs(dirPath, exist_ok=True)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":dirPath})
            url = "/api/v1/file?" + params
            conn.request("DELETE", url)
            response = conn.getresponse()
            common.assert_successful_response(response, HTTPStatus.OK)

            # verify file deleted
            assert_that(os.path.exists(dirPath)).is_false()
        finally:
            conn.close()

    def test_non_empty_directory_not_recursive(self):
        try:
            # prepare directory
            dirPath = os.path.join(common.TEST_TMP_DIR, "dir", "subdir")
            os.makedirs(dirPath, exist_ok=True)
            filePath = os.path.join(common.TEST_DATA_PATH, "text.txt")
            shutil.copy(filePath, dirPath)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":dirPath})
            url = "/api/v1/file?" + params
            conn.request("DELETE", url)
            response = conn.getresponse()
            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()

    def test_non_empty_directory_recursive(self):
        try:
            # prepare directory
            dirPath = os.path.join(common.TEST_TMP_DIR, "dir", "subdir")
            os.makedirs(dirPath, exist_ok=True)
            filePath = os.path.join(common.TEST_DATA_PATH, "text.txt")
            shutil.copy(filePath, dirPath)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"path":dirPath, "recursive":True})
            url = "/api/v1/file?" + params
            conn.request("DELETE", url)
            response = conn.getresponse()
            common.assert_successful_response(response, HTTPStatus.OK)

            # verify directory deleted
            assert_that(os.path.exists(dirPath)).is_false()
        finally:
            conn.close()

    def test_non_exist_file(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(common.TEST_TMP_DIR, "non-exist")
            params = urllib.parse.urlencode({"path":path})
            url = "/api/v1/file?" + params
            conn.request("DELETE", url)
            response = conn.getresponse()
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()
