import common
import http.client
import os
import platform
import pytest
import urllib
from assertpy import assert_that
from http import HTTPStatus

class TestExecScript:
    def test_combined_output(self):
        try:
            # prepare script
            scriptPath = os.path.join(common.TEST_EXE_PATH, "output-script.py")
            with open(scriptPath, "rb") as f:
                scriptContent = f.read()

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"program":"python3"})
            url = "/api/v1/script/execute?" + params
            conn.request("POST", url, scriptContent, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            
            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["exitStatus"]).is_equal_to(0)
            assert_that(data["error"]).is_equal_to("")
            expected_output = "stdout 1\nstderr 1\nstdout 2\nstderr 2\n"
            if "Windows" == platform.system():
                expected_output = expected_output.replace("\n", "\r\n")
            assert_that(data["output"]).is_equal_to(expected_output)
        finally:
            conn.close()

    def test_separate_output(self):
        try:
            # prepare script
            scriptPath = os.path.join(common.TEST_EXE_PATH, "output-script.py")
            with open(scriptPath, "rb") as f:
                scriptContent = f.read()

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"program":"python3", "separateOutput":True})
            url = "/api/v1/script/execute?" + params
            conn.request("POST", url, scriptContent, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            
            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["exitStatus"]).is_equal_to(0)
            assert_that(data["error"]).is_equal_to("")
            expected_stdout = "stdout 1\nstdout 2\n"
            expected_stderr = "stderr 1\nstderr 2\n"
            if "Windows" == platform.system():
                expected_stdout = expected_stdout.replace("\n", "\r\n")
                expected_stderr = expected_stderr.replace("\n", "\r\n")
            assert_that(data["stdout"]).is_equal_to(expected_stdout)
            assert_that(data["stderr"]).is_equal_to(expected_stderr)
        finally:
            conn.close()

    def test_execute_fail(self):
        try:
            # prepare script
            scriptPath = os.path.join(common.TEST_EXE_PATH, "fail-script.py")
            with open(scriptPath, "rb") as f:
                scriptContent = f.read()

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"program":"python3"})
            url = "/api/v1/script/execute?" + params
            conn.request("POST", url, scriptContent, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            
            # verify response
            expected_output = "start script\nexit with 1\n"
            if "Windows" == platform.system():
                expected_output = expected_output.replace("\n", "\r\n")
            expected_data = {
                "output": expected_output,
                "exitStatus": 1,
                "error": "exit status 1"
            }
            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR, expected_data)
        finally:
            conn.close()

    def test_execute_timeout(self):
        try:
            # prepare script
            scriptContent = "import time\nprint(\"sleeping\", flush=True)\ntime.sleep(3)"

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"program":"python3", "waitSeconds":1})
            url = "/api/v1/script/execute?" + params
            conn.request("POST", url, scriptContent, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            
            # verify response
            expected_output = "sleeping\n"
            if "Windows" == platform.system():
                expected_output = expected_output.replace("\n", "\r\n")
            expected_data = {
                "output": expected_output,
                "exitStatus": 0,
                "error": "7 - Command execution timeout"
            }
            common.assert_failed_response(response, HTTPStatus.REQUEST_TIMEOUT, expected_data)
        finally:
            conn.close()

    @pytest.mark.skipif("Windows" != platform.system(), reason="only applicable to windows")
    def test_execute_cmd(self):
        try:
            # prepare script
            scriptContent = "python3 --version"

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"program":"cmd"})
            url = "/api/v1/script/execute?" + params
            conn.request("POST", url, scriptContent, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            
            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["exitStatus"]).is_equal_to(0)
            assert_that(data["error"]).is_equal_to("")
            expected_output_pattern = r"Python \d+.\d+.\d+\r\n"
            assert_that(data["output"]).matches(expected_output_pattern)
        finally:
            conn.close()

    @pytest.mark.skipif("Linux" != platform.system(), reason="only applicable to linux")
    def test_execute_sh(self):
        try:
            # prepare script
            scriptContent = "python3 --version"

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"program":"sh"})
            url = "/api/v1/script/execute?" + params
            conn.request("POST", url, scriptContent, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            
            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["exitStatus"]).is_equal_to(0)
            assert_that(data["error"]).is_equal_to("")
            expected_output_pattern = r"Python \d+.\d+.\d+\n"
            assert_that(data["output"]).matches(expected_output_pattern)
        finally:
            conn.close()
