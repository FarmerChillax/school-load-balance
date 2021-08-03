# -*- coding: utf-8 -*-
'''
    :file: scan.py
    :author: -Farmer
    :url: https://blog.farmer233.top
    :date: 2021/08/03 01:41:34
'''
from os import stat
from sys import flags
import aiohttp
import time
from typing import List
import asyncio

from aiohttp.client_exceptions import ClientConnectionError, ClientError, ClientHttpProxyError, ServerDisconnectedError
from asyncio.exceptions import TimeoutError

EXCEPTIONS = (
    ClientConnectionError,
    ConnectionRefusedError,
    TimeoutError,
    ServerDisconnectedError,
    ClientError,
    ClientHttpProxyError,
    AssertionError,
    AttributeError,
)

class PortScan(object):

    def __init__(self, ip_list:List[str], all_ports:bool=False, rate:int=2000,
            protocol:str="http", flag:str="ZFSOFT.Inc", timeout:int=100) -> None:
        super().__init__()
        self.ip_list = ip_list
        self.rate = rate
        self.all_ports = all_ports
        self.protocol= protocol
        self.flag = flag
        self.timeout = timeout

    async def async_port_check(self, semaphore, ip_port):
        async with semaphore:
            ip, port = ip_port
            url = f'{self.protocol}://{ip}:{port}'
            conn = aiohttp.TCPConnector(ssl=False)
            async with aiohttp.ClientSession(connector=conn) as session:
                try:
                    # request = session.get(url=url)
                    # resp = await asyncio.wait_for(request, timeout=self.timeout)
                    # # if resp.headers.get("server") == self.flag:
                    # if resp.status == 200:
                    #     # 找到地址
                    #     # print(f"发现教务系统地址: {ip}:{port}")
                    #     return (ip, port, True)
                    # return (ip, port, False)
                    async with session.get(url=url, timeout=self.timeout) as resp:
                        # if resp.headers.get("server") == self.flag:
                        if resp.status == 200:
                            # 找到地址
                            return (ip, port, True)
                        return (ip, port, False)
                except Exception as e:
                    # 错误记录
                    # print(e)
                    return (ip, port, False)

    def callback(self, future):
        ip, port, status = future.result()
        if status:
            # 记录ip和port
            print(ip, port, "ok")



    def async_port_scan(self):
        ports = [port for port in range(1, 65536)] if self.all_ports else [80, 443, 8080, 8000, 8888, 5000, 4000]
        
        ip_port_list = [(ip, int(port)) for ip in self.ip_list for port in ports]
        # 限制并发量
        sem = asyncio.Semaphore(self.rate)
        loop = asyncio.get_event_loop()
        tasks = list()
        for ip_port in ip_port_list:
            task = asyncio.ensure_future(self.async_port_check(semaphore=sem, ip_port=ip_port))

            task.add_done_callback(self.callback)
            tasks.append(task)

        loop.run_until_complete(asyncio.wait(tasks))



if __name__ == '__main__':
    ip_list = ['192.168.2.123']
    now = time.time

    start = now()
    scan = PortScan(ip_list=ip_list, all_ports=True, rate=1024, flag="nginx/1.18.0 (Ubuntu)")
    scan.async_port_scan()
    end = now()
    print(f"start-end: {start}<->{end}")
