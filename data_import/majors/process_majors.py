import pandas as pd
import json
import numpy as np

df = pd.read_csv("专业名录2023.csv", encoding="GB18030", dtype={"专业代码": str})
df = df.replace({np.nan:None})
new_rows = []
category = ""
for row in df.to_dict(orient="records"):
    if row["序号"].isdigit():
        new_rows.append(row)
    elif row["序号"].endswith("学"):
        category = row["序号"]    
        print(category)
    row["大类"] = category
with open("../majors.json", "w") as f:
	json.dump(new_rows, f, ensure_ascii=False, indent=4)