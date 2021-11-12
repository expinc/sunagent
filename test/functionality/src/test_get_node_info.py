import common
import http.client
from assertpy import assert_that
from http import HTTPStatus

class TestGetNodeInfo:
    def test_ordinary(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/sys/info"
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()            
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["bootTime"]).matches(common.TIMESTAMP_PATTERN)
        finally:
            conn.close()
