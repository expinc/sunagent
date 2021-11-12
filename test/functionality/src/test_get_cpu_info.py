import common
import http.client
import psutil
from assertpy import assert_that
from http import HTTPStatus

class TestGetCpuInfo:
    def test_ordinary(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/sys/cpus/info"
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()            
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["count"]).is_equal_to(psutil.cpu_count())
        finally:
            conn.close()
