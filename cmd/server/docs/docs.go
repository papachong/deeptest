// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "https://github.com/aaronchen2k/deeptest/issues",
            "email": "462626@qq.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/projects": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "project create",
                "parameters": [
                    {
                        "description": "Create project Request Object",
                        "name": "ProjectReq",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/serverDomain.ProjectReq"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "serverDomain.ProjectReq": {
            "type": "object"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "3.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "DeepTest服务端API文档",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
