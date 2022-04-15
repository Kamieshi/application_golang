// Package app GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package app

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "tags": [
                    "Auth"
                ],
                "parameters": [
                    {
                        "description": "Login form",
                        "name": "formLogin",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.usPss"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/auth/logout": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Auth"
                ],
                "responses": {}
            }
        },
        "/auth/refresh": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Auth"
                ],
                "parameters": [
                    {
                        "description": "refresh",
                        "name": "refreshData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.rft"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/entity": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Entity"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Entity"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Entity"
                ],
                "parameters": [
                    {
                        "description": "Entity model",
                        "name": "DataEntity",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Entity"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/entity/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Entity"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of the mongo",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "desc",
                        "schema": {
                            "$ref": "#/definitions/models.Entity"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Entity"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of the mongo",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "Entity"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of the mongo",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {}
            }
        },
        "/user": {
            "post": {
                "tags": [
                    "User Control"
                ],
                "parameters": [
                    {
                        "description": "Create new User",
                        "name": "userData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.usPss"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "handlers.rft": {
            "type": "object",
            "properties": {
                "refresh": {
                    "type": "string"
                }
            }
        },
        "handlers.usPss": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.Entity": {
            "type": "object",
            "properties": {
                "id": {},
                "is_active": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "integer",
                    "maximum": 100,
                    "minimum": 1
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Golang Application Swagger",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
