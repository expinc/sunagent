import common
import http.client
import os
import platform
import time
from assertpy import assert_that
from http import HTTPStatus

ARCANE_CONTENT = """
timeout: 60
spells:
  debian:
    args: echo hello
  rhel:
    args: echo hello
  suse:
    args: echo hello"""

OS_TYPE_LIST = ["linux", "windows"]

GRIMOIRE_FOLDER = os.path.join(os.getcwd(), "gen", "grimoires")

class TestGrimoireRemove:

    def _generate_arcane_name(self):
        return "myarcane" + str(int(time.time()))

    def _get_current_os_type(self):
        if "Windows" == platform.system():
            return "windows"
        else:
            return "linux"

    def test_default_os(self):
        conn = http.client.HTTPConnection(common.HOST, common.PORT)
        try:
            # Prepare arcane to remove
            arcane_name = self._generate_arcane_name()
            url = "/api/v1/grimoires/default/arcanes/" + arcane_name
            conn.request("PUT", url, ARCANE_CONTENT, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()
            common.assert_successful_response(response, HTTPStatus.OK)
            grimoire_file = os.path.join(GRIMOIRE_FOLDER, self._get_current_os_type() + ".yaml")
            with open(grimoire_file, "r") as f:
                content = f.read()
            assert_that(content).contains(arcane_name)

            # Send request
            conn.request("DELETE", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # Verify result
            common.assert_successful_response(response, HTTPStatus.OK)
            with open(grimoire_file, "r") as f:
                content = f.read()
            assert_that(content).does_not_contain(arcane_name)
        finally:
            conn.close()

    def test_all_os(self):
        conn = http.client.HTTPConnection(common.HOST, common.PORT)
        try:
            for os_type in OS_TYPE_LIST:
                # Prepare arcane to remove
                arcane_name = self._generate_arcane_name()
                url = "/api/v1/grimoires/" + os_type + "/arcanes/" + arcane_name
                conn.request("PUT", url, ARCANE_CONTENT, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
                response = conn.getresponse()
                common.assert_successful_response(response, HTTPStatus.OK)
                grimoire_file = os.path.join(GRIMOIRE_FOLDER, os_type + ".yaml")
                with open(grimoire_file, "r") as f:
                    content = f.read()
                assert_that(content).contains(arcane_name)

                # Send request
                conn.request("DELETE", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
                response = conn.getresponse()

                # Verify result
                common.assert_successful_response(response, HTTPStatus.OK)
                with open(grimoire_file, "r") as f:
                    content = f.read()
                assert_that(content).does_not_contain(arcane_name)
        finally:
            conn.close()

    def test_unsupported_os(self):
        conn = http.client.HTTPConnection(common.HOST, common.PORT)
        try:
            # Send request
            url = "/api/v1/grimoires/nonexist/arcanes/nonexist"
            conn.request("DELETE", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # Verify result
            common.assert_failed_response(response, HTTPStatus.INTERNAL_SERVER_ERROR)
        finally:
            conn.close()

    def test_nonexist_arcane(self):
        conn = http.client.HTTPConnection(common.HOST, common.PORT)
        try:
            # Send request
            url = "/api/v1/grimoires/default/arcanes/nonexist"
            conn.request("DELETE", url, headers={"Authorization": "Basic " + common.BASIC_AUTH_TOKEN})
            response = conn.getresponse()

            # Verify result
            common.assert_failed_response(response, HTTPStatus.NOT_FOUND)
        finally:
            conn.close()
