import common
import http.client
import platform
import psutil
import subprocess
import time
import urllib
from assertpy import assert_that
from http import HTTPStatus

class TestProcessTerminate:
    def test_terminate_by_pid(self):
        try:
            # launch process
            if "Linux" == platform.system():
                args = ["python", common.TEST_DUMMY_PROC, "60", "dummyproc"]
            else:
                args = ["timeout", "60"]
            dummy_proc = subprocess.Popen(args)
            # wait for process ready
            time.sleep(1)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/processes/" + str(dummy_proc.pid) + "/terminate"
            conn.request("POST", url)
            response = conn.getresponse()
            
            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that(data[0]).is_equal_to(dummy_proc.pid)

            # wait for process completed termination
            time.sleep(1)
            # verify process terminated:
            assert_that(dummy_proc.poll()).is_not_none()
        finally:
            conn.close()
            dummy_proc.kill()

    def test_terminate_by_pid_non_exist(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/processes/" + str(32768) + "/terminate"
            conn.request("POST", url)
            response = conn.getresponse()

            # verify response
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()

    def test_terminate_by_name(self):
        try:
            # launch process
            if "Linux" == platform.system():
                args = ["python", common.TEST_DUMMY_PROC, "60", "dummyproc"]
            else:
                args = ["timeout", "60"]
            dummy_proc = subprocess.Popen(args)
            # wait for process ready
            time.sleep(1)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            if "Linux" == platform.system():
                url = "/api/v1/processes/dummyproc/terminate"
            else:
                url = "/api/v1/processes/timeout.exe/terminate"
            conn.request("POST", url)
            response = conn.getresponse()

            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that(data[0]).is_equal_to(dummy_proc.pid)

            # wait for process completed termination
            time.sleep(1)
            # verify process terminated:
            assert_that(dummy_proc.poll()).is_not_none()
        finally:
            conn.close()
            dummy_proc.kill()

    def test_terminate_by_name_non_exist(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/processes/non-exist/terminate"
            conn.request("POST", url)
            response = conn.getresponse()

            # verify response
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()
