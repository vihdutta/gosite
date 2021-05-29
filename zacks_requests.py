import requests
import time
import glob
import os
from lxml import html, etree
from bs4 import BeautifulSoup as bs
from openpyxl import Workbook, load_workbook


stocks = []
unfound_stocks = []

temp_xlsx_file = glob.glob(rf"{os.getcwd()}\temp-xlsx\*.xlsx")[0]

wb = load_workbook(temp_xlsx_file)

ws1 = wb.active
for cell in ws1['A']:
    stocks.append(cell)

def ZacksRequests(stocks):
    for stock in stocks:
        user_agent = {'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36',}
        r = requests.get(f'https://www.zacks.com/stock/quote/{stock.value}?q={stock.value}', headers=user_agent, verify=True)
        tree = html.fromstring(r.text)

        try:
            stock_data = tree.xpath('//*[@id="quote_ribbon_v2"]/div[2]/div[1]/p/text()')[0]
        except Exception as e:
            with open('analysis.txt', 'a') as f:
                f.write(f'{stock.value} : {e}\n')
            continue

        if stock_data.strip() == '':
            unfound_stocks.append(stock.value)
            continue

        with open('analysis.txt', 'a') as f:
            f.write(f'{stock.value} : {stock_data.strip()}\n')

    with open('analysis.txt', 'a') as f:
            f.write(f'\nNO ZACK RATING\n{unfound_stocks}')

ZacksRequests(stocks)
