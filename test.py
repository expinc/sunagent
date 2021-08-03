#!/usr/bin/env python

import argparse
import os
import platform
import psutil
import subprocess

TEST_OUT_DIR = os.path.join("gen", "test")

def parse_args():
    parser = argparse.ArgumentParser(description='Run tests')
    parser.add_argument("-t", "--type", required=False, choices=["unit", "func"])
    parser.add_argument("-p", "--part", required=False)

    return parser.parse_args()

def kill_proc(name):
    for proc in psutil.process_iter(["name"]):
        if proc.info["name"] == name:
            proc.kill()

def unit_test():
    print("=====================")
    print("Starting unit test...")
    cover_profile = os.path.join(TEST_OUT_DIR, "cover-profile")
    test_log = os.path.join(TEST_OUT_DIR, "unit-test.log")
    cmd = "go test -v -coverpkg=./... ./... -coverprofile {} > {}".format(cover_profile, test_log)
    os.system(cmd)
    print("Unit test finished. See {} for test log".format(test_log))

    print("Generating coverage report...")
    coverage_report = os.path.join(TEST_OUT_DIR, "unit-test.html")
    cmd = "go tool cover -html {} -o {}".format(cover_profile, coverage_report)
    os.system(cmd)
    print("See {} for coverage report".format(coverage_report))

def func_test(part):
    print("=====================")
    print("Starting functionality test...")

    print("Making executable...")
    sys_type = platform.system()
    if "Linux" == sys_type:
        subprocess.check_call(["./make.sh"])
    elif "Windows" == sys_type:
        subprocess.check_call(["make.bat"])
    else:
        raise Exception("Unsupported platform: {}".format(sys_type))

    print("Starting executable...")
    sys_type = platform.system()
    if "Linux" == sys_type:
        subprocess.check_call("gen/sunagentd &", shell=True)
    elif "Windows" == sys_type:
        subprocess.check_call("start gen\sunagentd.exe", shell=True)

    print("Running test cases...")
    cmd = "python -m pytest"
    if part:
        cmd += " " + part
    try:
        subprocess.check_call(cmd, shell=True)
        print("Test succeeded")
    except subprocess.CalledProcessError as err:
        print("Test failed!")

    print("Stopping executable...")
    if "Linux" == sys_type:
        kill_proc("sunagentd")
    elif "Windows" == sys_type:
        kill_proc("sunagentd.exe")

if __name__ == "__main__":
    args = parse_args()
    os.makedirs(TEST_OUT_DIR, exist_ok=True)
    if "unit" == args.type:
        unit_test()
    elif "func" == args.type:
        func_test(args.part)
    else:
        unit_test()
        func_test()
