import requests




r = requests.post("https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal", data={
    "app_id": "cli_a4c936cd91f8500d",
    "app_secret": "S7g6tDGloq9kxnDDuSkhRf2XJEOBC3NO"
})

res = r.json()
print(res)
token = res["tenant_access_token"]
print("token:", token)

def request_subscription():
    r = requests.post("https://www.feishu.cn/approval/openapi/v2/subscription/subscribe", data={
        "approval_code":  "0266C9C4-96ED-4185-9A72-D3997222835E"
    }, headers={
        "Authorization": "Bearer "+token,
        #   "Content-Type": "application/json",
        "Connection": "keep-alive"})
    print("已经订阅事件", r.content)

request_subscription()
r = requests.post("https://www.feishu.cn/approval/openapi/v2/instance/create",
                  data={
                      "approval_code": "0266C9C4-96ED-4185-9A72-D3997222835E",
                      "user_id": "6cg59473",
                      "department_id": "g71ca93cf249d1cd",
                      "node_approver_user_id_list": {

                      },
                      "form": "[{\"id\": \"reason\", \"value\": \"111111\", \"type\": \"input\"}]"
                  },
                  headers={
                      "Authorization": "Bearer "+token,
                      #   "Content-Type": "application/json",
                      "Connection": "keep-alive"}
                  )
print(r.content, r.json(), r.headers)
instance_code = r.json()["data"]["instance_code"]
print("instance_code", instance_code)
print("已经发送请求，编号为：", instance_code)