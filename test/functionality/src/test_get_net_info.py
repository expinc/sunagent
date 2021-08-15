import common
import http.client
import psutil
from assertpy import assert_that
from http import HTTPStatus

class TestGetNetInfo:
    def test_ordinary(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/sys/net/info"
            conn.request("GET", url)
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(len(data)).is_equal_to(len(psutil.net_if_stats()))
        finally:
            conn.close()
