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
class TestPackageUpgrade:

    @classmethod
    def setup_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR, ignore_errors=True)
        os.makedirs(common.TEST_TMP_DIR, exist_ok=True)

        # download test package
        url = common.TEST_PKGS[distro.id()]["url"]
        if not common.download_file(url, common.TEST_TMP_DIR):
            raise Exception("Failed to download {}".format(url))
        # download test package of newer version
        url = common.TEST_PKGS[distro.id()]["newUrl"]
        if not common.download_file(url, common.TEST_TMP_DIR):
            raise Exception("Failed to download {}".format(url))

    @classmethod
    def teardown_class(cls):
        shutil.rmtree(common.TEST_TMP_DIR)

    def test_upgrade_non_installed(self):
        test_package = common.TEST_PKGS[distro.id()]["name"]
        common.remove_package(test_package)

        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            url = "/api/v1/package/" + test_package
            conn.request("PUT", url)
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

    def test_upgrade_ordinary(self):
        test_package = common.TEST_PKGS[distro.id()]["name"]
        common.remove_package(test_package)
        common.install_package(test_package)
        origin_version = common.get_package_version(test_package)

        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            file_path = os.path.join(common.TEST_TMP_DIR, common.TEST_PKGS[distro.id()]["newFile"])
            params = urllib.parse.urlencode({"path":file_path})
            url = "/api/v1/package?" + params
            conn.request("PUT", url)
            response = conn.getresponse()

            data = common.get_successful_response(response, HTTPStatus.OK)
            assert_that(data["name"]).is_equal_to(test_package)
            assert_that(data["version"]).is_not_empty()
            assert_that(data["architecture"]).is_not_empty()
            assert_that(data["summary"]).is_not_empty()

            installed = common.is_package_installed(test_package)
            assert_that(installed).is_equal_to(True)
            assert_that(common.get_package_version(test_package)).is_not_equal_to(origin_version)
        finally:
            conn.close()
            common.remove_package(test_package)

    def test_upgrade_with_earlier(self):
        test_package = common.TEST_PKGS[distro.id()]["name"]
        common.remove_package(test_package)
        file_path = os.path.join(common.TEST_TMP_DIR, common.TEST_PKGS[distro.id()]["newFile"])
        common.install_package(file_path)
        origin_version = common.get_package_version(test_package)

        try:
            conn = http.client.HTTPConnection(common.HOST, common.PORT)
            file_path = os.path.join(common.TEST_TMP_DIR, common.TEST_PKGS[distro.id()]["file"])
            params = urllib.parse.urlencode({"path":file_path})
            url = "/api/v1/package?" + params
            conn.request("PUT", url)
            response = conn.getresponse()

            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()
            common.remove_package(test_package)
