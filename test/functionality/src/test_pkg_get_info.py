import common
import distro
import http.client
import platform
import pytest
from assertpy import assert_that
from http import HTTPStatus

@pytest.mark.skipif("Linux" != platform.system(), reason="Package management only implemented for linux")
class TestPackageGetInfo:
    @classmethod
    def setup_class(cls):
        if "ubuntu" == distro.id() or "debian" == distro.id():
            cls.must_have_pkg = "apt"
        # there may be "opensuse-leap" as distro.id()
        # so the distribution is determined by substring for opensuse
        elif "centos" == distro.id() or "opensuse" in distro.id():
            cls.must_have_pkg = "rpm"

    def test_ordinary(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/package/" + self.must_have_pkg
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["name"]).is_equal_to(self.must_have_pkg)
            assert_that(data["version"]).is_not_empty()
            assert_that(data["architecture"]).is_not_empty()
            assert_that(data["summary"]).is_not_empty()
        finally:
            conn.close()

    def test_not_installed(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/package/" + "notexist"
            conn.request("GET", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()
