import common
import http.client
import psutil
from assertpy import assert_that
from http import HTTPStatus

class TestGetDiskInfo:
    def test_ordinary(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/sys/disks/stats"
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(len(data)).is_equal_to(len(psutil.disk_partitions()))
        finally:
            conn.close()
