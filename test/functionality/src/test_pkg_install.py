import common
import distro
import http.client
import os
import platform
import pytest
import shutil
import urllib.parse
from assertpy import assert_that
from http import HTTPStatus

@pytest.mark.skipif("Linux" != platform.system(), reason="Package management only implemented for linux")
class TestPackageInstall:

    @classmethod
    def setup_class(cls):
        if "ubuntu" == distro.id() or "debian" == distro.id():
            cls.must_have_pkg = "apt"
        # there may be "opensuse-leap" as distro.id()
        # so the distribution is determined by substring for opensuse
        elif "centos" == distro.id() or "opensuse" in distro.id():
            cls.must_have_pkg = "rpm"

        shutil.rmtree(common.TEST_TMP_DIR, ignore_errors=True)
        os.makedirs(common.TEST_TMP_DIR, exist_ok=True)
        if not common.download_file(common.TEST_PKGS[distro.id()][2], common.TEST_TMP_DIR):
            raise Exception("Failed to download {}".format(common.TEST_PKGS[distro.id()][2]))

    @classmethod
    def teardown_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR)

    def test_ordinary(self):
        test_package = common.TEST_PKGS[distro.id()][0]
        common.remove_package(test_package)

        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/package/" + test_package
            conn.request("POST", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["name"]).is_equal_to(test_package)
            assert_that(data["version"]).is_not_empty()
            assert_that(data["architecture"]).is_not_empty()
            assert_that(data["summary"]).is_not_empty()

            installed = common.is_package_installed(test_package)
            assert_that(installed).is_equal_to(True)
        finally:
            conn.close()
            common.remove_package(test_package)

    def test_ordinary_by_file(self):
        test_package = common.TEST_PKGS[distro.id()][0]
        common.remove_package(test_package)

        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            file_path = os.path.join(common.TEST_TMP_DIR, common.TEST_PKGS[distro.id()][1])
            params = urllib.parse.urlencode({"path":file_path})
            url = "/api/v1/package?" + params
            conn.request("POST", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["name"]).is_equal_to(test_package)
            assert_that(data["version"]).is_not_empty()
            assert_that(data["architecture"]).is_not_empty()
            assert_that(data["summary"]).is_not_empty()

            installed = common.is_package_installed(test_package)
            assert_that(installed).is_equal_to(True)
        finally:
            conn.close()
            common.remove_package(test_package)

    def test_non_exist(self):
        test_package = "nonexist"
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/package/" + test_package
            conn.request("POST", url)
            response = conn.getresponse()
            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()

    def test_non_exist_by_file(self):
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            file_path = os.path.join(common.TEST_TMP_DIR, "nonexist")
            params = urllib.parse.urlencode({"path":file_path})
            url = "/api/v1/package?" + params
            conn.request("POST", url)
            response = conn.getresponse()
            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()

    def test_already_installed(self):
        test_package = self.must_have_pkg
        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/package/" + test_package
            conn.request("POST", url)
            response = conn.getresponse()
            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()

    def test_already_installed_by_file(self):
        test_package = common.TEST_PKGS[distro.id()][0]
        common.install_package(test_package)

        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            file_path = os.path.join(common.TEST_TMP_DIR, "nonexist")
            params = urllib.parse.urlencode({"path":file_path})
            url = "/api/v1/package?" + params
            conn.request("POST", url)
            response = conn.getresponse()
            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()
            common.remove_package(test_package)
