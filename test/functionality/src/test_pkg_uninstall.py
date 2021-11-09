import common
import distro
import http.client
import platform
import pytest
from assertpy import assert_that
from http import HTTPStatus

@pytest.mark.skipif("Linux" != platform.system(), reason="Package management only implemented for linux")
class TestPackageInstall:

    def test_ordinary(self):
        test_package = common.TEST_PKGS[distro.id()]["name"]
        common.install_package(test_package)

        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/package/" + test_package
            conn.request("DELETE", url)
            response = conn.getresponse()

            common.assert_successful_response(response, HTTPStatus.OK)

            installed = common.is_package_installed(test_package)
            assert_that(installed).is_equal_to(False)
        finally:
            conn.close()
            common.remove_package(test_package)

    def test_uninstall_nonexist(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/package/nonexist"
            conn.request("DELETE", url)
            response = conn.getresponse()

            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()
