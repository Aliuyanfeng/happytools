'''
Author: LiuYanFeng
Date: 2026-03-12 09:44:34
LastEditors: LiuYanFeng
LastEditTime: 2026-03-13 11:39:38
FilePath: \happytools\example\download\main.py
Description: 像珍惜礼物一样珍惜今天

Copyright (c) 2026 by ${git_name_email}, All Rights Reserved. 
'''
import requests
import json
from urllib import parse

url = "https://0.0.0.0:8888/external/fileDownload"
# params = {"file_id": "019cdab4bace7ac89128d426c8828146"}
params = {"file_id": "069b671304d4788cad3a5daa0eb25cf3"}
# 添加 Range 头，从第 1MB 开始下载
headers = {
    # "Range": "bytes=1048576-"  # 从第 1MB 开始
}
response = requests.get(url, params=params, headers=headers, stream=True, verify=False)
# print(response)
if response.status_code == 200:
    # 从响应头获取文件名，没有则默认
    content_disposition = response.headers.get("Content-Disposition", "")
    file_name = "download_file.zip"
    if "filename*=" in content_disposition:
        # 优先取 filename*（RFC 5987 标准，支持中文）
        encoded_name = content_disposition.split("filename*=")[-1].strip()
        # 去掉 UTF-8'' 前缀再解码
        file_name = parse.unquote(encoded_name.split("''")[-1])
    elif "filename=" in content_disposition:
        # 降级取 filename
        file_name = parse.unquote(content_disposition.split("filename=")[-1].strip())

    with open(file_name, "wb") as f:
        for chunk in response.iter_content(chunk_size=1024 * 1024):  # 1MB
            if chunk:
                f.write(chunk)

    print(f"下载完成: {response.status_code}, {file_name}")
else:
    print(f"下载失败: {response.status_code}, {response.text}")