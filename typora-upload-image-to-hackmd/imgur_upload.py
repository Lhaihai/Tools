import requests
import sys
import json
from base64 import b64encode

URL = "https://api.imgur.com/3/image"

CONTENT_TYPE_TABLE = {
    'jpg': 'image/jpeg',
    'jpeg': 'image/jpeg',
    'png': 'image/png',
    'bmp': 'application/x-bmp',
    'gif': 'image/gif'
}




def upload(client_id: str, imgPath: str):
    suffix = imgPath.split('.')[-1]
    if suffix.lower() in CONTENT_TYPE_TABLE:
        contentType = CONTENT_TYPE_TABLE[suffix.lower()]
    else:
        contentType = 'image/jpeg'

    data = {
        'image': b64encode(open(imgPath, 'rb').read()),
        'type': 'base64',
        'name': '1.jpg',
        'title': 'Picture no. 1'
    }

    headers = {
        'Authorization': 'Client-ID {}'.format(client_id)
    }
    try:
        response = requests.request("POST", URL, headers=headers, data=data)

        result = json.loads(response.content)
        # print(response.text.encode('utf8'))
        return result['data']['link']
    except Exception as e:
        return e.__str__()


if __name__ == '__main__':
    link = upload(sys.argv[1], sys.argv[2])
    print('Upload Success:')
    print(link)