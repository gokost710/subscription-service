package swagger

import "github.com/swaggo/swag"

const docName = "swagger"

type spec struct{}

func Register() {
	swag.Register(docName, spec{})
}

func (spec) ReadDoc() string {
	return openAPISpec
}

const openAPISpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Subscription Service API",
    "description": "REST API for aggregating user online subscriptions.",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:8080"
    }
  ],
  "paths": {
    "/health": {
      "get": {
        "summary": "Health check",
        "responses": {
          "200": {
            "description": "Service is healthy",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/HealthResponse"
                }
              }
            }
          }
        }
      }
    },
    "/subscriptions": {
      "post": {
        "summary": "Create subscription",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SubscriptionRequest"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Subscription created",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SubscriptionResponse"
                }
              }
            }
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "500": {
            "$ref": "#/components/responses/InternalServerError"
          }
        }
      },
      "get": {
        "summary": "List subscriptions",
        "parameters": [
          {
            "name": "user_id",
            "in": "query",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "name": "service_name",
            "in": "query",
            "schema": {
              "type": "string"
            }
          },
          {
            "name": "limit",
            "in": "query",
            "schema": {
              "type": "integer",
              "default": 50,
              "maximum": 100
            }
          },
          {
            "name": "offset",
            "in": "query",
            "schema": {
              "type": "integer",
              "default": 0,
              "minimum": 0
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Subscription list",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SubscriptionListResponse"
                }
              }
            }
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "500": {
            "$ref": "#/components/responses/InternalServerError"
          }
        }
      }
    },
    "/subscriptions/{id}": {
      "get": {
        "summary": "Get subscription by ID",
        "parameters": [
          {
            "$ref": "#/components/parameters/SubscriptionID"
          }
        ],
        "responses": {
          "200": {
            "description": "Subscription found",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SubscriptionResponse"
                }
              }
            }
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "404": {
            "$ref": "#/components/responses/NotFound"
          },
          "500": {
            "$ref": "#/components/responses/InternalServerError"
          }
        }
      },
      "put": {
        "summary": "Update subscription",
        "parameters": [
          {
            "$ref": "#/components/parameters/SubscriptionID"
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SubscriptionRequest"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Subscription updated",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SubscriptionResponse"
                }
              }
            }
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "404": {
            "$ref": "#/components/responses/NotFound"
          },
          "500": {
            "$ref": "#/components/responses/InternalServerError"
          }
        }
      },
      "delete": {
        "summary": "Delete subscription",
        "parameters": [
          {
            "$ref": "#/components/parameters/SubscriptionID"
          }
        ],
        "responses": {
          "204": {
            "description": "Subscription deleted"
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "404": {
            "$ref": "#/components/responses/NotFound"
          },
          "500": {
            "$ref": "#/components/responses/InternalServerError"
          }
        }
      }
    },
    "/subscriptions/summary": {
      "get": {
        "summary": "Calculate total subscription price",
        "parameters": [
          {
            "name": "from",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "01-2025"
            }
          },
          {
            "name": "to",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "12-2025"
            }
          },
          {
            "name": "user_id",
            "in": "query",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          },
          {
            "name": "service_name",
            "in": "query",
            "schema": {
              "type": "string"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Total price calculated",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SummaryResponse"
                }
              }
            }
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "500": {
            "$ref": "#/components/responses/InternalServerError"
          }
        }
      }
    }
  },
  "components": {
    "parameters": {
      "SubscriptionID": {
        "name": "id",
        "in": "path",
        "required": true,
        "schema": {
          "type": "integer",
          "format": "int64",
          "minimum": 1
        }
      }
    },
    "responses": {
      "BadRequest": {
        "description": "Invalid request",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/ErrorResponse"
            }
          }
        }
      },
      "NotFound": {
        "description": "Subscription not found",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/ErrorResponse"
            }
          }
        }
      },
      "InternalServerError": {
        "description": "Internal server error",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/ErrorResponse"
            }
          }
        }
      }
    },
    "schemas": {
      "HealthResponse": {
        "type": "object",
        "properties": {
          "status": {
            "type": "string",
            "example": "ok"
          }
        }
      },
      "SubscriptionRequest": {
        "type": "object",
        "required": [
          "service_name",
          "price",
          "user_id",
          "start_date"
        ],
        "properties": {
          "service_name": {
            "type": "string",
            "example": "Yandex Plus"
          },
          "price": {
            "type": "integer",
            "minimum": 0,
            "example": 400
          },
          "user_id": {
            "type": "string",
            "format": "uuid",
            "example": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
          },
          "start_date": {
            "type": "string",
            "example": "07-2025"
          },
          "end_date": {
            "type": "string",
            "nullable": true,
            "example": "12-2025"
          }
        }
      },
      "SubscriptionResponse": {
        "allOf": [
          {
            "type": "object",
            "required": [
              "id"
            ],
            "properties": {
              "id": {
                "type": "integer",
                "format": "int64",
                "example": 1
              }
            }
          },
          {
            "$ref": "#/components/schemas/SubscriptionRequest"
          }
        ]
      },
      "SubscriptionListResponse": {
        "type": "object",
        "properties": {
          "items": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/SubscriptionResponse"
            }
          }
        }
      },
      "SummaryResponse": {
        "type": "object",
        "properties": {
          "total_price": {
            "type": "integer",
            "example": 1200
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "error": {
            "type": "string"
          }
        }
      }
    }
  }
}`
