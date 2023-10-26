import aiohttp

from cache.cache_token import get_token_from_cache, set_token_to_cache
from config import GET_TOKEN_URL


async def get_copilot_token(github_token):
    copilot_token = get_token_from_cache(github_token)
    if not copilot_token:
        headers = {
            "Authorization": f"token {github_token}"
        }
        async with aiohttp.ClientSession as session:
            async with session.get(GET_TOKEN_URL, headers=headers) as resp:
                if resp.status != 200:
                    return resp.status, await resp.text()
                copilot_token = await resp.json()
                # put into cache
                set_token_to_cache(github_token, copilot_token)
    return 200, copilot_token
