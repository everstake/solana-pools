// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/epoch": {
            "get": {
                "description": "get epoch",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "RestAPI",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/tools.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/v1.epoch"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "404": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "default": {
                        "description": "default response",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    }
                }
            }
        },
        "/pool-statistic": {
            "get": {
                "description": "get statistic by pool",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "RestAPI",
                "parameters": [
                    {
                        "type": "string",
                        "default": "2021-01-01T15:04:05Z",
                        "description": "first date for aggregation",
                        "name": "from",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "2021-12-01T15:04:05Z",
                        "description": "second date for aggregation",
                        "name": "to",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "everSOL",
                        "description": "pool name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "enum": [
                            "day",
                            "week",
                            "month",
                            "year"
                        ],
                        "type": "string",
                        "description": "aggregation",
                        "name": "aggregation",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/tools.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/v1.poolStatistic"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "404": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "default": {
                        "description": "default response",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    }
                }
            }
        },
        "/pool/{name}": {
            "get": {
                "description": "get pool",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "WebSocket",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Pool name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/tools.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/v1.PoolDetails"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "404": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "default": {
                        "description": "default response",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    }
                }
            }
        },
        "/pools": {
            "get": {
                "description": "get pools",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "RestAPI",
                "parameters": [
                    {
                        "type": "number",
                        "default": 1,
                        "description": "offset for aggregation",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "default": 10,
                        "description": "limit for aggregation",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "stake-pool name",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/tools.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/v1.poolMainPage"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "404": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "default": {
                        "description": "default response",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    }
                }
            }
        },
        "/pools-statistic": {
            "get": {
                "description": "get statistic",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "WebSocket",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/tools.ResponseData"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/v1.TotalPoolsStatistic"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "404": {
                        "description": "bad request",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    },
                    "default": {
                        "description": "default response",
                        "schema": {
                            "$ref": "#/definitions/tools.ResponseError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "tools.ResponseData": {
            "type": "object",
            "properties": {
                "data": {}
            }
        },
        "tools.ResponseError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "v1.PoolDetails": {
            "type": "object",
            "properties": {
                "active_stake": {
                    "type": "number"
                },
                "address": {
                    "type": "string"
                },
                "apy": {
                    "type": "number"
                },
                "avg_score": {
                    "type": "integer"
                },
                "avg_skipped_slots": {
                    "type": "number"
                },
                "delinquent": {
                    "type": "number"
                },
                "depossit_fee": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "rewards_fee": {
                    "type": "number"
                },
                "staking_accounts": {
                    "type": "integer"
                },
                "tokens_supply": {
                    "type": "number"
                },
                "unstake_liquidity": {
                    "type": "number"
                },
                "validators": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.Validator"
                    }
                },
                "withdrawal_fee": {
                    "type": "number"
                }
            }
        },
        "v1.TotalPoolsStatistic": {
            "type": "object",
            "properties": {
                "avg_performance_score": {
                    "type": "integer"
                },
                "max_performance_score": {
                    "type": "integer"
                },
                "min_performance_score": {
                    "type": "integer"
                },
                "network_apy": {
                    "type": "number"
                },
                "pools": {
                    "type": "integer"
                },
                "skipped_slot": {
                    "type": "number"
                },
                "total_active_stake": {
                    "type": "number"
                },
                "total_active_stake_pool": {
                    "type": "number"
                },
                "total_unstake_liquidity": {
                    "type": "number"
                },
                "total_validators": {
                    "type": "integer"
                },
                "usd": {
                    "type": "number"
                }
            }
        },
        "v1.Validator": {
            "type": "object",
            "properties": {
                "apy": {
                    "type": "number"
                },
                "data_center": {
                    "type": "string"
                },
                "fee": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "node_pk": {
                    "type": "string"
                },
                "pool_active_stake": {
                    "type": "number"
                },
                "score": {
                    "type": "integer"
                },
                "skipped_slots": {
                    "type": "number"
                },
                "total_active_stake": {
                    "type": "number"
                },
                "vote_pk": {
                    "type": "string"
                }
            }
        },
        "v1.epoch": {
            "type": "object",
            "properties": {
                "end_epoch": {
                    "type": "string"
                },
                "epoch": {
                    "type": "integer"
                },
                "progress": {
                    "type": "integer"
                },
                "slots_in_epoch": {
                    "type": "integer"
                },
                "sps": {
                    "type": "number"
                }
            }
        },
        "v1.poolMainPage": {
            "type": "object",
            "properties": {
                "active_stake": {
                    "type": "number"
                },
                "address": {
                    "type": "string"
                },
                "apy": {
                    "type": "number"
                },
                "avg_score": {
                    "type": "integer"
                },
                "avg_skipped_slots": {
                    "type": "number"
                },
                "delinquent": {
                    "type": "number"
                },
                "depossit_fee": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "rewards_fee": {
                    "type": "number"
                },
                "staking_accounts": {
                    "type": "integer"
                },
                "tokens_supply": {
                    "type": "number"
                },
                "unstake_liquidity": {
                    "type": "number"
                },
                "validators": {
                    "type": "integer"
                },
                "withdrawal_fee": {
                    "type": "number"
                }
            }
        },
        "v1.poolStatistic": {
            "type": "object",
            "properties": {
                "active_stake": {
                    "type": "number"
                },
                "apy": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string"
                },
                "number_of_validators": {
                    "type": "integer"
                },
                "unstacked_liquidity": {
                    "type": "number"
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
	BasePath:    "/v1",
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
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
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
	swag.Register("swagger", &s{})
}
