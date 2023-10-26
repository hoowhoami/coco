import time

tokens = {}
error_tokens = {}


def set_token_to_cache(github_token, copilot_token):
    tokens[github_token] = copilot_token


def get_token_from_cache(github_token):
    if copilot_token := tokens[github_token]:
        if copilot_token["expires_at"] > int(time.time()) + 60:
            return copilot_token
    return None
