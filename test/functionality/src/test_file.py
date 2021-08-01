#!/usr/bin/env python

import common
import http.client
import os
import os.path
import platform
import re
from assertpy import assert_that
from http import HTTPStatus

TEST_DATA_DIR = os.path.join(os.getcwd(), "test", "functionality", "data", "dir")

# sample: "2021-07-31T22:44:17.7724489+08:00"
TIMESTAMP_PATTERN = r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+\+\d{2}:\d{2}$"
# sample: "-rw-rw-rw-"
REGULAR_FILE_MODE_PATTERN = r"^-((r|-)(w|-)(x|-)){3}$"

class TestFile:
    def test_get_file_meta_text_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            path = os.path.join(TEST_DATA_DIR, "text.txt")
            url = common.url_with_params("/api/v1/fileMeta", {"path":path})
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
