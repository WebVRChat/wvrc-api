import requests

def test(n, desc, method, url, data, fail=False):
    print("# test %d: %s" % (n, desc))
    response = getattr(requests, method.lower())(url, data=data, headers={
        "Content-Type": "application/x-www-form-urlencoded",
    })
    if response.status_code == 400 and not fail:
        print("# test %d: \033[91mfailed\033[0m" % n)
        print("# %s" % response.json()["reason"])
        print()
        return

    print("# test %d: \033[92mpassed\033[0m" % n)
    print("# %s" % response.text)
    print()

    try:
        return response.json()
    except:
        return response.text

test(
    1, 
    "creating user", 
    "POST", 
    "http://localhost:5050/api/register",
    {"username": "test1", "password": "test2"}
)

token = test(
    2,
    "login user",
    "POST",
    "http://localhost:5050/api/login",
    {"username": "test1", "password": "test2"}
)['token']

invite = test(
    3,
    "creating room",
    "POST",
    "http://localhost:5050/api/room/create",
    {"token": token, "name": "Olydri", "description": "pouloulou", "is_private": "false"}
)['invite']

test(
    4,
    "joining room as owner",
    "PUT",
    "http://localhost:5050/api/room/join",
    {"token": token, "invite": invite},
    True
)

requests.post('http://localhost:5050/api/register', data={"username": "test2", "password": "test2"}, headers={"Content-Type": "application/x-www-form-urlencoded"})
token2 = requests.post('http://localhost:5050/api/login', data={"username": "test2", "password": "test2"}, headers={"Content-Type": "application/x-www-form-urlencoded"}).json()['token']

test(
    5,
    "joining room as user",
    "PUT",
    "http://localhost:5050/api/room/join",
    {"token": token2, "invite": invite},
)

test(
    6,
    "joining room as user x2",
    "PUT",
    "http://localhost:5050/api/room/join",
    {"token": token2, "invite": invite},
    True
)

test(
    13,
    "listing room users",
    "POST",
    "http://localhost:5050/api/room/users",
    {"token": token2, "invite": invite},
    False
)

test(
    7,
    "leaving room as owner",
    "PUT",
    "http://localhost:5050/api/room/leave",
    {"token": token, "invite": invite},
    True
)

test(
    8,
    "deleting room as user",
    "PUT",
    "http://localhost:5050/api/room/delete",
    {"token": token2, "invite": invite},
    True
)

test(
    9,
    "leaving room as user",
    "PUT",
    "http://localhost:5050/api/room/leave",
    {"token": token2, "invite": invite},
)

test(
    10,
    "leaving room as user x2",
    "PUT",
    "http://localhost:5050/api/room/leave",
    {"token": token2, "invite": invite},
    True
)

test(
    11,
    "deleting room as owner",
    "PUT",
    "http://localhost:5050/api/room/delete",
    {"token": token, "invite": invite}
)

test(
    12,
    "deleting room as owner x2",
    "PUT",
    "http://localhost:5050/api/room/delete",
    {"token": token, "invite": invite},
    True
)