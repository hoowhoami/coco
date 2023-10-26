from typing import Optional

from werkzeug.test import EnvironBuilder
from werkzeug.wrappers import Request


def fake_request(method: str, headers: dict, json: Optional[dict] = None) -> Request:
    """
    根据提供的json数据和headers，构造一个request对象
    """

    # 使用 EnvironBuilder 创建一个新的请求环境
    builder = EnvironBuilder(method=method, headers=headers, json=json)
    env = builder.get_environ()

    return Request(env)
