import common
import http.client
from assertpy import assert_that
from http import HTTPStatus

# sample: 2021-08-13T21:40:44+08:00
TIME_FORMAT = r"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\+\d{2}:\d{2}"

class TestGetNodeInfo:
    def test_ordinary(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/sys/info"
            conn.request("GET", url)
            response = conn.getresponse()            
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["bootTime"]).matches(TIME_FORMAT)
        finally:
            conn.close()
