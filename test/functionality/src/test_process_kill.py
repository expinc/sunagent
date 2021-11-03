import common
import http.client
import platform
import psutil
import subprocess
import time
import urllib
from assertpy import assert_that
from http import HTTPStatus

SIG_KILL = 9
SIG_USR1 = 10
SIG_TERM = 15

class TestProcessKill:
    def test_kill_by_pid(self):
        try:
            # launch process
            if "Linux" == platform.system():
                args = ["python3", common.TEST_DUMMY_PROC, "60", "dummyproc"]
            else:
                args = ["timeout", "60"]
            dummy_proc = subprocess.Popen(args)
            # wait for process ready
            time.sleep(1)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/processes/" + str(dummy_proc.pid) + "/kill"
            conn.request("POST", url)
            response = conn.getresponse()
            
            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that(data[0]).is_equal_to(dummy_proc.pid)

            # verify process behavior:
            if "Linux" == platform.system():
                assert_that(dummy_proc.wait()).is_equal_to(-SIG_TERM)
            else:
                assert_that(psutil.pid_exists(dummy_proc.pid)).is_false()
        finally:
            conn.close()
            dummy_proc.kill()

    def test_kill_by_pid_non_exist(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/processes/" + str(32768) + "/kill"
            conn.request("POST", url)
            response = conn.getresponse()

            # verify response
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()

    def test_kill_by_pid_specify_signal(self):
        try:
            # launch process
            if "Linux" == platform.system():
                args = ["python3", common.TEST_DUMMY_PROC, "60", "dummyproc"]
            else:
                args = ["timeout", "60"]
            dummy_proc = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            # wait for process ready
            time.sleep(1)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"signal":SIG_USR1})
            url = "/api/v1/processes/" + str(dummy_proc.pid) + "/kill?" + params
            conn.request("POST", url)
            response = conn.getresponse()

            if "Linux" == platform.system():
                # verify response
                data = common.get_successful_response(response, HTTPStatus.OK)
                assert_that(1).is_equal_to(len(data))
                assert_that(data[0]).is_equal_to(dummy_proc.pid)

                # verify process behavior:
                assert_that(dummy_proc.wait()).is_equal_to(-SIG_USR1)
            else:
                common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()
            dummy_proc.kill()

    def test_kill_by_name(self):
        try:
            # launch process
            if "Linux" == platform.system():
                args = ["python3", common.TEST_DUMMY_PROC, "60", "dummyproc"]
            else:
                args = ["timeout", "60"]
            dummy_proc = subprocess.Popen(args)
            # wait for process ready
            time.sleep(1)

            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            if "Linux" == platform.system():
                url = "/api/v1/processes/dummyproc/kill"
            else:
                url = "/api/v1/processes/timeout.exe/kill"
            conn.request("POST", url)
            response = conn.getresponse()

            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            assert_that(data[0]).is_equal_to(dummy_proc.pid)

            # verify process behavior:
            if "Linux" == platform.system():
                assert_that(dummy_proc.wait()).is_equal_to(-SIG_TERM)
            else:
                assert_that(psutil.pid_exists(dummy_proc.pid)).is_false()
        finally:
            conn.close()
            dummy_proc.kill()

    def test_kill_by_name_non_exist(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/processes/non-exist/kill"
            conn.request("POST", url)
            response = conn.getresponse()

            # verify response
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()
