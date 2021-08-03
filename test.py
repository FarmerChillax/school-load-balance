# -*- coding: utf-8 -*-
'''
    :file: test.py
    :author: -Farmer
    :url: https://blog.farmer233.top
    :date: 2021/08/03 20:17:35
'''
from scanner import Scan
import time

if __name__ == '__main__':
    host = 'farmer233.top'
    now = time.time

    start = now()
    scan = Scan(all_ports=False, rate=1024, flag="nginx/1.18.0 (Ubuntu)")
    scan.run()
    end = now()
    print(f"start-end: {start}<->{end}, spend: {end-start}s\n")
    print(scan.open_list)
    print(len(scan.open_list))
