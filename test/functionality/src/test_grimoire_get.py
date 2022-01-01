import common
import http.client
from assertpy import assert_that
from http import HTTPStatus

OS_TYPE_LIST = ["linux", "windows"]

class TestGrimoireGet:

    def test_default_os(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/grimoires/default"
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify result
            data = common.get_binary_response(response, HTTPStatus.OK)
            assert_that(data.decode()).starts_with("arcanes:")
        finally:
            conn.close()

    def test_specifying_os(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)

            for os_type in OS_TYPE_LIST:
                # send request
                url = "/api/v1/grimoires/" + os_type
                conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
                response = conn.getresponse()

                # verify result
                data = common.get_binary_response(response, HTTPStatus.OK)
                assert_that(data.decode()).starts_with("arcanes:")
        finally:
            conn.close()

    def test_unsupported_os(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/grimoires/nonexist"
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify result
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()
