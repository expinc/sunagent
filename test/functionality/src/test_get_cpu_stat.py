import common
import http.client
import psutil
import urllib.parse
from assertpy import assert_that
from http import HTTPStatus

class TestGetCpuStat:
    def test_all_cpus(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/sys/cpus/stats"
            conn.request("GET", url)
            response = conn.getresponse()            
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(len(data["usages"])).is_equal_to(1)
        finally:
            conn.close()

    def test_per_cpus(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"perCpu":True})
            url = "/api/v1/sys/cpus/stats?" + params
            conn.request("GET", url)
            response = conn.getresponse()            
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(len(data["usages"])).is_equal_to(psutil.cpu_count())
        finally:
            conn.close()
