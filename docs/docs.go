// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/cron/resend": {
            "post": {
                "description": "api resend ` + "`" + `failed` + "`" + ` out messages",
                "tags": [
                    "Retry"
                ],
                "summary": "api resend failed out messages",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/cron/retry": {
            "post": {
                "description": "api resend ` + "`" + `wait retry` + "`" + ` in messages, message will change status to",
                "tags": [
                    "Retry"
                ],
                "summary": "api retry ` + "`" + `wait retry` + "`" + ` in messages",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/cron/retry_previous": {
            "post": {
                "description": "api resend ` + "`" + `wait retry previous` + "`" + ` in messages, just retry in messages",
                "tags": [
                    "Retry"
                ],
                "summary": "api retry ` + "`" + `wait retry previous` + "`" + ` in messages",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/queue/messages": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "api publish out message to amqp",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Out Messages"
                ],
                "summary": "publish message to amqp",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schema.OutMsgBodyParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        },
                        "headers": {
                            "Token": {
                                "type": "string",
                                "description": "qwerty"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/queue/routing_keys": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "api create routing key",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Routing Keys"
                ],
                "summary": "create routing key",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/schema.RoutingCreateParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        },
                        "headers": {
                            "Token": {
                                "type": "string",
                                "description": "qwerty"
                            }
                        }
                    }
                }
            }
        },
        "/api/v1/queue/routing_keys/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "api retrieve routing key",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Routing Keys"
                ],
                "summary": "api retrieve routing key",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Routing Key ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "schema.OutMsgBodyParam": {
            "type": "object",
            "required": [
                "payload",
                "routing_key"
            ],
            "properties": {
                "origin_code": {
                    "type": "string",
                    "example": "code"
                },
                "origin_model": {
                    "type": "string",
                    "example": "model"
                },
                "payload": {
                    "type": "object"
                },
                "routing_key": {
                    "type": "string",
                    "example": "routing.key"
                }
            }
        },
        "schema.RoutingCreateParam": {
            "type": "object",
            "required": [
                "api_method",
                "api_url",
                "group",
                "name",
                "value"
            ],
            "properties": {
                "api_method": {
                    "type": "string"
                },
                "api_url": {
                    "type": "string"
                },
                "group": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "value": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
