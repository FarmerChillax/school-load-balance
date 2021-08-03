# -*- coding: utf-8 -*-
'''
    :file: flashScan.py
    :author: -Farmer
    :url: https://blog.farmer233.top
    :date: 2021/08/03 14:10:59
'''

# -*- coding: utf-8 -*-
'''
    :file: scan.py
    :author: -Farmer
    :url: https://blog.farmer233.top
    :date: 2021/08/03 01:41:34
'''

from settings import SCAN_PROTOCOL, SCAN_FLAG, SCAN_RATE, SCAN_TIMEOUT
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


class PortScan(object):

    def __init__(self, host:str, all_ports:bool=False, rate:int=SCAN_RATE,
            protocol:str=SCAN_PROTOCOL, flag:str=SCAN_FLAG, timeout:int=SCAN_TIMEOUT) -> None:
        super().__init__()
        self.host = host
        self.rate = rate
        self.all_ports = all_ports
        self.protocol= protocol
        self.flag = flag
        self.timeout = timeout
        self.open_list = []

    async def async_port_check(self, semaphore:asyncio.Semaphore, ip_port):
        ip, port = ip_port
        url = f'{self.protocol}://{ip}:{port}'
        async with semaphore:
            print(f'scanning {ip}:{port}...')
            async with aiohttp.ClientSession(connector=aiohttp.TCPConnector(ssl=False)) as session:
                try:
                    async with session.get(url, timeout=self.timeout) as resp:
                        if resp.headers.get("server") == self.flag:
                            return (ip, port, True)
                except Exception as e:
                    # 记录错误日志？
                    pass
            return (ip, port, False)


    def callback(self, future):
        ip, port, status = future.result()
        if status:
            # 记录ip和port
            self.open_list.append(port)


    def async_port_scan(self):
        ports = [port for port in range(1, 65536)] if self.all_ports else [80, 443, 8080, 8000, 8888, 5000, 4000]
        url_port_list = [(self.host, port) for port in ports]

        # 限制并发量
        sem = asyncio.Semaphore(self.rate)
        # 任务队列
        loop = asyncio.get_event_loop()
        tasks = list()

        for ip_port in url_port_list:
            task = asyncio.ensure_future(self.async_port_check(semaphore=sem, ip_port=ip_port))
            task.add_done_callback(self.callback)
            tasks.append(task)

        loop.run_until_complete(asyncio.wait(tasks))

    def run(self):
        self.async_port_scan()


if __name__ == '__main__':
    host = 'farmer233.top'
    now = time.time

    start = now()
    scan = PortScan(host=host, all_ports=True, rate=20000, flag="nginx/1.18.0 (Ubuntu)")
    scan.run()
    end = now()
    print(f"start-end: {start}<->{end}, spend: {end-start}s\n")
    print(scan.open_list)
    print(len(scan.open_list))
