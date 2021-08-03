# -*- coding: utf-8 -*-
'''
    :file: scan.py
    :author: -Farmer
    :url: https://blog.farmer233.top
    :date: 2021/08/03 20:37:46
'''

from asyncio import tasks
from .settings import SCAN_PROTOCOL, SCAN_FLAG, SCAN_RATE, SCAN_SEGMENT, SCAN_TIMEOUT
import aiohttp
import time
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


class Scan(object):

    def __init__(self, all_ports:bool=False, rate:int=SCAN_RATE,
            protocol:str=SCAN_PROTOCOL, flag:str=SCAN_FLAG, timeout:int=SCAN_TIMEOUT) -> None:
        super().__init__()
        # self.segment = segment
        self.rate = rate
        self.all_ports = all_ports
        self.protocol= protocol
        self.flag = flag
        self.timeout = timeout
        self.common_port = [80, 443, 8080, 8000, 8888, 5000, 4000]
        self.open_list = {}



    async def async_port_check(self, semaphore:asyncio.Semaphore, ip_port):
        ip, port = ip_port
        url = f'{self.protocol}://{ip}:{port}'
        async with semaphore:
            print(f'scanning {ip}:{port}...; flag: {self.flag}')
            async with aiohttp.ClientSession(connector=aiohttp.TCPConnector(ssl=False)) as session:
                try:
                    async with session.get(url, timeout=self.timeout) as resp:
                        # if resp.status == 200:
                        #     return (ip, port, True)
                        if resp.headers.get("server") == self.flag:
                            return (ip, port, True)
                except Exception as e:
                    # 记录错误日志？
                    pass
            return (ip, port, False)


    def callback(self, future):
        """扫描回调函数
            用于记录扫描信息、写库操作等。
        Args:
            future ([type]): 扫描结果
        """
        ip, port, status = future.result()
        if status:
            # 记录ip和port
            try:
                if ip in self.open_list:
                    self.open_list[ip].append(port)
                else:
                    self.open_list[ip] = [port]
            except Exception as e:
                print(e)

    def async_scan(self, ip_port_list):
        # 限制并发量
        sem = asyncio.Semaphore(self.rate)
        # 任务队列
        loop = asyncio.get_event_loop()
        tasks = list()
        for ip_port in ip_port_list:
            task = asyncio.ensure_future(self.async_port_check(semaphore=sem, ip_port=ip_port))
            task.add_done_callback(self.callback)
            tasks.append(task)

        loop.run_until_complete(asyncio.wait(tasks))
        pass

    def async_port_scan(self, host:str):
        """异步端口扫描
                传入扫描目标主机
        Args:
            host (str): 目标主机
        """

        ports = [port for port in range(1, 65536)] if self.all_ports else self.common_port
        ip_port_list = [(host, port) for port in ports]

        self.async_scan(ip_port_list=ip_port_list)


    def async_segment_scan(self, segment:str=SCAN_SEGMENT):
        """网段扫描器
            目前只支持扫描ipv4地址段的最后一段，即: 172.16.1.x
        Args:
            segment (str, optional): 扫描网段. Defaults to SCAN_SEGMENT.
        """

        host_list = [segment.format(i) for i in range(1, 256)]
        ports = [port for port in range(1, 65536)] if self.all_ports else self.common_port
        ip_port_list = [(host, port) for host in host_list for port in ports]

        self.async_scan(ip_port_list)


    def run(self, host=None):
        """扫描器启动模式
            如果指定host，则只扫描目标主机。留空则扫描配置文件的网段。
        Args:
            host (str, optional): 扫描目标主机. Defaults to None.
        """
        scan_segment = True if host is None else False
        if not scan_segment:
            self.async_port_scan(host)
        else:
            self.async_segment_scan()




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
