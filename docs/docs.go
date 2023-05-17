// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms",
        "contact": {
            "name": "Songyue Chen",
            "email": "enderbybear@foxmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/questionbox/question/new": {
            "post": {
                "description": "通过提问表单，在数据库中新增一个问题",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "questionbox"
                ],
                "summary": "新增一个问题",
                "parameters": [
                    {
                        "description": "新增问题请求",
                        "name": "newQuestionReq",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.NewQuestionReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "新增问题响应“",
                        "schema": {
                            "$ref": "#/definitions/models.NewQuestionRes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/svcerror.SvcErr"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/svcerror.SvcErr"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.EntityWithName": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models.NewQuestionReq": {
            "type": "object",
            "properties": {
                "description": {
                    "description": "问题描述",
                    "type": "string"
                },
                "major": {
                    "description": "提问专业",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.EntityWithName"
                        }
                    ]
                },
                "questioner": {
                    "description": "提问人信息",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.PersonalInfo"
                        }
                    ]
                },
                "school": {
                    "description": "提问学校",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.EntityWithName"
                        }
                    ]
                },
                "title": {
                    "description": "问题标题",
                    "type": "string"
                }
            }
        },
        "models.NewQuestionRes": {
            "type": "object",
            "properties": {
                "docID": {
                    "description": "新增问题ID",
                    "type": "string"
                }
            }
        },
        "models.PersonalInfo": {
            "type": "object",
            "properties": {
                "CEEPlace": {
                    "description": "高考所在地",
                    "type": "string"
                },
                "age": {
                    "description": "年龄",
                    "type": "integer"
                },
                "gender": {
                    "description": "性别",
                    "type": "string"
                },
                "situation": {
                    "description": "具体情况",
                    "type": "string"
                },
                "subject": {
                    "description": "高考科目",
                    "type": "string"
                }
            }
        },
        "svcerror.SvcErr": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "msg": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "127.0.0.1",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "UniPingo Backend",
	Description:      "This is the backend for UniPingo application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
