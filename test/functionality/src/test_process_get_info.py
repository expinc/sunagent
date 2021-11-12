import common
import http.client
import platform
import subprocess
from assertpy import assert_that
from http import HTTPStatus

class TestGetProcInfo:

    dummy_proc1 = None
    dummy_proc2 = None

    @classmethod
    def setup_class(cls):
        args = ["python3", common.TEST_DUMMY_PROC, "60"]
        cls.dummy_proc1 = subprocess.Popen(args)
        cls.dummy_proc2 = subprocess.Popen(args)

    @classmethod
    def teardown_class(cls):
        cls.dummy_proc1.kill()
        cls.dummy_proc2.kill()

    def test_get_by_pid(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/processes/" + str(self.dummy_proc1.pid)
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(1).is_equal_to(len(data))
            info = data[0]
            assert_that(info["pid"]).is_equal_to(self.dummy_proc1.pid)
        finally:
            conn.close()

    def test_get_by_pid_non_exist(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/processes/32768"
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()            
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()

    def test_get_by_name(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            name = "python3.exe"
            if "Linux" == platform.system():
                name = "python3"
            url = "/api/v1/processes/" + name
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            
            data = common.get_successful_response(response, HTTPStatus.OK)
            print("data: {}".format(data))
            assert_that(2 + common.TEST_SCRIPT_PYTHON_PROC_COUNT).is_equal_to(len(data))
            pids = [info["pid"] for info in data]
            assert_that(pids).contains(self.dummy_proc1.pid, self.dummy_proc2.pid)
        finally:
            conn.close()

    def test_get_by_name_non_exist(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/processes/non-exist"
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()
