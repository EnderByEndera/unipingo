
任务：
- 获取本地Python解释器的依赖，将相关依赖复制到临时文件夹中，并生成dockerfile。将临时文件夹上传到服务端，此时可在临时文件夹中执行构建Docker的流程。

login pgsql:
```sh
psql -h 127.0.0.1 -p 5432 -U melodie -d postgres
```