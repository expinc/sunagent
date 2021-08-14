import common
import http.client
from assertpy import assert_that
from http import HTTPStatus

class TestGetMemStat:
    def test_ordinary(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/sys/mem/stats"
            conn.request("GET", url)
            response = conn.getresponse()
            common.assert_successful_response(response, HTTPStatus.OK)
        finally:
            conn.close()
