#!/usr/bin/env python

import setproctitle
import sys
import time

if __name__ == "__main__":
    if len(sys.argv) > 2:
        setproctitle.setproctitle(sys.argv[2])
    time.sleep(int(sys.argv[1]))
