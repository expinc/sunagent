import common
import http.client
import time
import urllib
from assertpy import assert_that
from http import HTTPStatus

class TestJobCancel:

    def _create_exec_script_job(self, conn, duration_sec):
        scriptContent = "import time\nprint(\"sleeping\", flush=True)\ntime.sleep({})".format(duration_sec)
        params = urllib.parse.urlencode({"program":"python3", "async":"true"})
        url = "/api/v1/script/execute?" + params
        conn.request("POST", url, scriptContent, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
        response = conn.getresponse()
        return common.get_successful_response(response, HTTPStatus.OK)

    def _cancel_job(self, conn, id):
        url = "/api/v1/jobs/" + id + "/cancel"
        conn.request("POST", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
        return conn.getresponse()

    def _get_job(self, conn, id):
        url = "/api/v1/jobs/" + id
        conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
        return conn.getresponse()

    def test_ordinary(self):
        try:
            # create job
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            data = self._create_exec_script_job(conn, 10)
            id = data["id"]

            # cancel job
            response = self._cancel_job(conn, id)
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["id"]).is_equal_to(id)

            # verify status
            response = self._get_job(conn, id)
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["id"]).is_equal_to(id)
            assert_that(data["status"]).is_equal_to("CANCELED")
            assert_that(data["progress"]).is_less_than(100)
        finally:
            conn.close()

    def test_cancel_again(self):
        try:
            # create job
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            data = self._create_exec_script_job(conn, 10)
            id = data["id"]

            # cancel job
            response = self._cancel_job(conn, id)
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["id"]).is_equal_to(id)

            # cancel again
            response = self._cancel_job(conn, id)
            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR, None)
        finally:
            conn.close()

    def test_cancel_nonexist(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            response = self._cancel_job(conn, "nonexist")
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND, None)
        finally:
            conn.close()
