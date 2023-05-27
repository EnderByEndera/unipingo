import requests




r = requests.post("https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal", data={
    "app_id": "cli_a4fa09d1f9fb5013",
    "app_secret": "Y3cuyQB716gJ6rnStNrbJbXVJWsRQiJb"
})

res = r.json()
print(res)
token = res["tenant_access_token"]
print("token:", token)

APPROVAL_CODE = "EAEEC8EC-99F4-4D9B-B7C4-00563DF880B7"# 审批表单的序号，可在表单编辑的URL的definitionCode中获取

def request_subscription():
    r = requests.post("https://www.feishu.cn/approval/openapi/v2/subscription/subscribe", data={
        "approval_code":  APPROVAL_CODE 
    }, headers={
        "Authorization": "Bearer "+token,
        #   "Content-Type": "application/json",
        "Connection": "keep-alive"})
    print("已经订阅事件", r.content)

request_subscription()
r = requests.post("https://www.feishu.cn/approval/openapi/v2/instance/create",
                  data={
                      "approval_code": "EAEEC8EC-99F4-4D9B-B7C4-00563DF880B7",
                      "user_id": "f4cg5853", # 在管理员后台-组织架构-成员与部门中，选择自己，可以定义user_id
                      "department_id": "5", # 开发部门的id，在管理员后台-组织架构-成员与部门中，选择开发部，即可找到。
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