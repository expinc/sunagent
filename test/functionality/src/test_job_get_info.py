import common
import http.client
import time
import urllib
from assertpy import assert_that
from http import HTTPStatus

class TestJobGetInfo:

    def _create_exec_script_job(self, conn, duration_sec):
        scriptContent = "import time\nprint(\"sleeping\", flush=True)\ntime.sleep({})".format(duration_sec)
        params = urllib.parse.urlencode({"program":"python3", "async":"true"})
        url = "/api/v1/script/execute?" + params
        conn.request("POST", url, scriptContent, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
        response = conn.getresponse()
        return common.get_successful_response(response, HTTPStatus.OK)

    def _get_job(self, conn, id):
        url = "/api/v1/jobs/" + id
        conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
        return conn.getresponse()

    def test_get_successful_job(self):
        try:
            # create job
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            info = self._create_exec_script_job(conn, 1)

            # get job and verify status
            time.sleep(3)
            response = self._get_job(conn, info["id"])
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["id"]).is_equal_to(info["id"])
            assert_that(data["status"]).is_equal_to("SUCCESSFUL")
            assert_that(data["progress"]).is_equal_to(100)

        finally:
            conn.close()

    def test_get_failed_job(self):
        try:
            # create job
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            info = self._create_exec_script_job(conn, -1)

            # get job and verify status
            time.sleep(1)
            response = self._get_job(conn, info["id"])
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["id"]).is_equal_to(info["id"])
            assert_that(data["status"]).is_equal_to("FAILED")

        finally:
            conn.close()

    def test_get_canceled_job(self):
        try:
            # create job
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            info = self._create_exec_script_job(conn, 10)

            # cancel job
            url = "/api/v1/jobs/" + info["id"] + "/cancel"
            conn.request("POST", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            conn.getresponse().read()

            # get job and verify status
            time.sleep(1)
            response = self._get_job(conn, info["id"])
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["id"]).is_equal_to(info["id"])
            assert_that(data["status"]).is_equal_to("CANCELED")

        finally:
            conn.close()

    def test_get_executing_job(self):
        try:
            # create job
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            info = self._create_exec_script_job(conn, 10)

            # get job and verify status
            time.sleep(1)
            response = self._get_job(conn, info["id"])
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["id"]).is_equal_to(info["id"])
            assert_that(data["status"]).is_equal_to("EXECUTING")
            assert_that(data["progress"]).is_less_than(100)

        finally:
            conn.close()

    def test_get_nonexist_job(self):
        try:
            # get job and verify response
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            response = self._get_job(conn, "nonexist")
            data = common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()
