import common
import http.client
import time
import urllib
from assertpy import assert_that
from http import HTTPStatus

class TestJobListAll:

    def _create_exec_script_job(self, conn, duration_sec):
        scriptContent = "import time\nprint(\"sleeping\", flush=True)\ntime.sleep({})".format(duration_sec)
        params = urllib.parse.urlencode({"program":"python3", "async":"true"})
        url = "/api/v1/script/execute?" + params
        conn.request("POST", url, scriptContent, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
        response = conn.getresponse()
        return common.get_successful_response(response, HTTPStatus.OK)

    def test_ordinary(self):
        try:
            # create jobs
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            ids = []
            ids.append(self._create_exec_script_job(conn, 1)["id"])
            ids.append(self._create_exec_script_job(conn, -1)["id"])
            ids.append(self._create_exec_script_job(conn, 10)["id"])

            # list jobs
            time.sleep(3)
            url = "/api/v1/jobs"
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            res_ids = [info["id"] for info in data]
            for id in ids:
                assert_that(id in res_ids).is_true()
        finally:
            conn.close()
