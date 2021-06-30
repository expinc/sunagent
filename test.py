#!/usr/bin/env python

import argparse
import os

TEST_OUT_DIR = os.path.join("gen", "test")

def parse_args():
    parser = argparse.ArgumentParser(description='Run tests')
    parser.add_argument("-t", "--type", required=False, choices=["unit"])

    return parser.parse_args()

def unit_test():
    cover_profile = os.path.join(TEST_OUT_DIR, "cover-profile")
    test_log = os.path.join(TEST_OUT_DIR, "unit-test.log")
    cmd = "go test -v -coverpkg=./... ./... -coverprofile {} > {}".format(cover_profile, test_log)
    os.system(cmd)

    test_render = os.path.join(TEST_OUT_DIR, "unit-test.html")
    cmd = "go tool cover -html {} -o {}".format(cover_profile, test_render)
    os.system(cmd)


if __name__ == "__main__":
    args = parse_args()
    os.makedirs(TEST_OUT_DIR, exist_ok=True)
    if "unit" == args.type:
        unit_test()
    else:
        unit_test()
