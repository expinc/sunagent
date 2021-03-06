import common
import http.client
from http import HTTPStatus

class TestManage:
    def test_get_info(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            conn.request("GET", "/api/v1/info")
            response = conn.getresponse()
            
            expected_data = {"version": "2.0.0"}
            common.assert_successful_response(response, HTTPStatus.OK, expected_data)
        finally:
            conn.close()
