import requests
import json
import sqlite3
import time
from urllib.parse import quote_plus
import hmac
import hashlib
import base64

'''
如果没有认证的请求，则每分钟最多10次请求(10 requests per minute)。所以每次爬取之后sleep 6秒
github api 参数参考：https://blog.csdn.net/Next_Second/article/details/78238328
'''

# 这里放钉钉机器人的token
token = "xxxxx"

def getTimestampAndSign():
    timestamp = (round(time.time() * 1000))
    secret = 'SEC1exxxxxxxx'
    secret_enc = secret.encode('utf-8')
    string_to_sign = '{}\n{}'.format(timestamp, secret)
    string_to_sign_enc = (string_to_sign).encode('utf-8')
    hmac_code = hmac.new(secret_enc, string_to_sign_enc, digestmod=hashlib.sha256).digest()
    sign = quote_plus(base64.b64encode(hmac_code))
    return sign,timestamp

def push_DingMes(content):

    sign,timestamp = getTimestampAndSign()
    baseUrl = "https://oapi.dingtalk.com/robot/send?access_token={}&timestamp={}&sign={}".format(token,timestamp,sign)

    # please set charset= utf-8
    HEADERS = {
        "Content-Type": "application/json ;charset=utf-8 "
    }

    stringBody ={
        "msgtype": "markdown",
        "markdown": {
            "title": "发现了新的CVE",
            "text": content
        },

        # "msgtype": "text",
        # "text": {
        #     "content": content
        # },
        "at": {
            # "atMobiles": ["13207488983"],
            # "isAtAll": True
        }
 }
    MessageBody = json.dumps(stringBody)
    result = requests.post(url=baseUrl, data=MessageBody, headers=HEADERS)
    # print(result)


def getall2020cve():
    nr = requests.get(url='https://api.github.com/search/repositories?q=CVE-2020&sort=updated&order=desc&per_page=1').text
    data = json.loads(nr)
    for i in range(int(data["total_count"]/50)+1):
        nr = requests.get(url='https://api.github.com/search/repositories?q=CVE-2020&sort=updated&order=desc&per_page=50&page={}'.format(i+1)).text
        data = json.loads(nr)
        for d in data['items']:
            if check(d["id"]):
                insert(d["id"],d["svn_url"],d["description"],d["updated_at"])
        time.sleep(6)

def getall2019cve():
    nr = requests.get(url='https://api.github.com/search/repositories?q=CVE-2019&sort=updated&order=desc&per_page=1').text
    data = json.loads(nr)
    for i in range(int(data["total_count"]/50)+1):
        nr = requests.get(url='https://api.github.com/search/repositories?q=CVE-2019&sort=updated&order=desc&per_page=50&page={}'.format(i+1)).text
        data = json.loads(nr)
        for d in data['items']:
            if check(d["id"]):
                insert(d["id"],d["svn_url"],d["description"],d["updated_at"])
        time.sleep(6)

def getallrce():
    nr = requests.get(url='https://api.github.com/search/repositories?q=rce poc&sort=updated&order=desc&per_page=1').text
    data = json.loads(nr)
    for i in range(int(data["total_count"]/50)+1):
        nr = requests.get(url='https://api.github.com/search/repositories?q=rce poc&sort=updated&order=desc&per_page=50&page={}'.format(i+1)).text
        data = json.loads(nr)
        for d in data['items']:
            if check(d["id"]):
                insert(d["id"],d["svn_url"],d["description"],d["updated_at"])
        time.sleep(6)

def get2020cve():
    nr = requests.get(url='https://api.github.com/search/repositories?q=CVE-2020&sort=updated&order=desc&per_page=50').text
    data = json.loads(nr)
    for d in data['items']:
        if check(d["id"]):
            insert(d["id"],d["svn_url"],d["description"],d["updated_at"])
            with open("log.txt", 'a+') as file_object:
                file_object.write(str(d["svn_url"])+"  "+str(d["description"])+"\n")
            content = """## Github 发现了新漏洞
url: {url}

描述: {description}

发现时间: {create_time}

请及时查看和处理
""".format(url=d["svn_url"], description=d["description"],create_time=time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
            push_DingMes(content)


def get2019cve():
    nr = requests.get(url='https://api.github.com/search/repositories?q=CVE-2019&sort=updated&order=desc&per_page=50').text
    data = json.loads(nr)
    for d in data['items']:
        if check(d["id"]):
            insert(d["id"],d["svn_url"],d["description"],d["updated_at"])
            with open("log.txt", 'a+') as file_object:
                file_object.write(str(d["svn_url"])+"  "+str(d["description"])+"\n")
            content = """## Github 发现了新漏洞
url: {url}

描述: {description}

发现时间: {create_time}

请及时查看和处理
""".format(url=d["svn_url"], description=d["description"],create_time=time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
            push_DingMes(content)

def getrce():
    nr = requests.get(url='https://api.github.com/search/repositories?q=rce poc&sort=updated&order=desc&per_page=50').text
    data = json.loads(nr)
    for d in data['items']:
        if check(d["id"]):
            insert(d["id"],d["svn_url"],d["description"],d["updated_at"])
            with open("log.txt", 'a+') as file_object:
                file_object.write(str(d["svn_url"])+"  "+str(d["description"])+"\n")
            content = """## Github 发现了新漏洞
url: {url}

描述: {description}

发现时间: {create_time}

请及时查看和处理
""".format(url=d["svn_url"], description=d["description"],create_time=time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()))
            push_DingMes(content)


def init():
    conn = sqlite3.connect('test.db')
    cursor = conn.cursor()
    cursor.execute('create table cvelist (id varchar(20) primary key, url varchar(60),description varchar(500),time varchar(40))')
    cursor.close()
    conn.commit()
    conn.close()

def insert(id,url,description,time):
    conn = sqlite3.connect('test.db')
    cursor = conn.cursor()
    print(url,description,time)
    cursor.execute('insert into cvelist (id, url, description, time) values (\'{}\', \'{}\',\'{}\', \'{}\')'.format(str(id),str(url),str(description).replace("'",""),str(time)))
    cursor.close()
    conn.commit()
    conn.close()

def check(id):
    conn = sqlite3.connect('test.db')
    cursor = conn.cursor()
    cursor.execute('select id from cvelist where id={}'.format(id))
    values = cursor.fetchall()
    if len(values):
        return False
    else:
        return True


if __name__ == "__main__":
    try:
        init()
    except:
        pass
    # getall2019cve()
    # getall2020cve()
    # getallrce()
    while(1):
        try:
            print(str(time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())) + "  进行监测扫描")
            get2019cve()
            get2020cve()
            getrce()
            time.sleep(600)
        except KeyboardInterrupt:
            exit(0)
        except:
            pass


