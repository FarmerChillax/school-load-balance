# -*- coding: utf-8 -*-
'''
    :file: scam.py
    :author: -Farmer
    :url: https://blog.farmer233.top
    :date: 2021/08/03 18:52:43
'''
from typing import List
# from .port import PortScan
from settings import SCAN_PROTOCOL, SCAN_FLAG, SCAN_RATE, SCAN_TIMEOUT


class Scan():

    def __init__(self, ip_list: List[str], scan_all_port: bool = False, rate:int=SCAN_RATE,
            falg: str = SCAN_FLAG, protocol: str = SCAN_PROTOCOL, timeout:int=SCAN_TIMEOUT) -> None:
        self.ip_list = ip_list
        self.scan_all_port = scan_all_port
        self.rate = rate
        self.falg = falg
        self.protocol = protocol
        self.timeout = timeout
        self.result = {}

    def run():
        pass

# Network segment
if __name__ == '__main__':
    segment = '192.168.2.{}'
    ip_list = [segment % i for i in range(255)]
    print(ip_list[:10])
    print(ip_list[-10:])
