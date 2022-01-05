import common
import distro
import http.client
import platform
import pytest
import time
import urllib
from assertpy import assert_that
from http import HTTPStatus

class TestGrimoireCastArcane:

    @classmethod
    def setup_class(cls):
        if "ubuntu" == distro.id() or "debian" == distro.id():
            cls.must_have_pkg = "apt"
        # there may be "opensuse-leap" as distro.id()
        # so the distribution is determined by substring for opensuse
        elif "centos" == distro.id() or "opensuse" in distro.id():
            cls.must_have_pkg = "rpm"

    @pytest.mark.skipif("Linux" != platform.system(), reason="No arcane for windows")
    def test_ordinary(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/grimoires/default/arcanes/get-package/cast"
            conn.request("POST", url, self.must_have_pkg, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["exitStatus"]).is_equal_to(0)
            assert_that(data["error"]).is_equal_to("")
        finally:
            conn.close()

    @pytest.mark.skipif("Linux" != platform.system(), reason="No arcane for windows")
    def test_exec_fail(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/grimoires/default/arcanes/get-package/cast"
            conn.request("POST", url, "nonexist", headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify response
            data = common.get_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
            assert_that(data["exitStatus"]).is_equal_to(1)
            assert_that(data["error"]).is_equal_to("exit status 1")
        finally:
            conn.close()

    def test_non_exist_arcane(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/grimoires/default/arcanes/nonexist/cast"
            conn.request("POST", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify response
            data = common.get_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()

    def test_other_os(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            os_type = "windows"
            if "Windows" == platform.system():
                os_type = "linux"
            url = "/api/v1/grimoires/{}/arcanes/nonexist/cast".format(os_type)
            conn.request("POST", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify response
            data = common.get_failed_response(response, HTTPStatus.NOT_ACCEPTABLE)
        finally:
            conn.close()

    @pytest.mark.skipif("Linux" != platform.system(), reason="No arcane for windows")
    def test_async(self):
        try:
            # send request
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            params = urllib.parse.urlencode({"async":"true"})
            url = "/api/v1/grimoires/default/arcanes/get-package/cast?" + params
            conn.request("POST", url, self.must_have_pkg, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify response
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["name"]).is_equal_to("CastArcane")
            assert_that(data["status"]).is_equal_to("SPAWNED")

            # get result
            time.sleep(3)
            id = data["id"]
            url = "/api/v1/jobs/" + id
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # verify job status
            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["status"]).is_equal_to("SUCCESSFUL")

            # verify result
            expected_output_head = None
            if "ubuntu" == distro.id() or "debian" == distro.id():
                expected_output_head = "Package: apt\n"
            # there may be "opensuse-leap" as distro.id()
            # so the distribution is determined by substring for opensuse
            elif "centos" == distro.id() or "opensuse" in distro.id():
                expected_output_head = "Name        : rpm\n"
            assert_that(data["result"]["output"]).starts_with(expected_output_head)
        finally:
            conn.close()
