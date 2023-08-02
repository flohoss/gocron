// Code generated by swaggo/swag. DO NOT EDIT.

package docs

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
        "/commands": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "commands"
                ],
                "parameters": [
                    {
                        "description": "Command body",
                        "name": "command",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.CommandBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/jobs": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/database.Job"
                            }
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "parameters": [
                    {
                        "description": "Job",
                        "name": "job",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/database.Job"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/database.Job"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "parameters": [
                    {
                        "description": "job",
                        "name": "job",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/database.Job"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/database.Job"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/jobs/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Job ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/database.Job"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Job ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/jobs/{id}/runs": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Job ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/database.Run"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/echo.HTTPError"
                        }
                    }
                }
            }
        },
        "/system": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "system"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/system.Data"
                        }
                    }
                }
            }
        },
        "/system/logs": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "system"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/database.SystemLog"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.CommandBody": {
            "type": "object",
            "required": [
                "command"
            ],
            "properties": {
                "command": {
                    "type": "string"
                },
                "custom_command": {
                    "type": "string"
                },
                "job_id": {
                    "type": "integer"
                },
                "local_directory": {
                    "type": "string",
                    "example": "/"
                },
                "password_file_path": {
                    "type": "string",
                    "example": "/secrets/.resticpwd"
                },
                "restic_remote": {
                    "type": "string",
                    "example": "rclone:pcloud:Backups/gitea"
                }
            }
        },
        "database.Command": {
            "type": "object",
            "required": [
                "command",
                "sort_id",
                "type"
            ],
            "properties": {
                "command": {
                    "type": "string",
                    "example": "docker compose stop"
                },
                "file_output": {
                    "type": "string",
                    "example": ".dbBackup.sql"
                },
                "id": {
                    "type": "integer"
                },
                "job_id": {
                    "type": "integer"
                },
                "sort_id": {
                    "type": "integer"
                },
                "type": {
                    "type": "integer",
                    "enum": [
                        1,
                        2
                    ]
                }
            }
        },
        "database.CompressionType": {
            "type": "integer",
            "enum": [
                1,
                2,
                3
            ],
            "x-enum-varnames": [
                "Automatic",
                "Maximum",
                "NoCompression"
            ]
        },
        "database.Job": {
            "type": "object",
            "required": [
                "compression_type",
                "description",
                "local_directory",
                "password_file_path",
                "restic_remote",
                "retention_policy"
            ],
            "properties": {
                "compression_type": {
                    "enum": [
                        1,
                        2,
                        3
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/database.CompressionType"
                        }
                    ],
                    "example": 1
                },
                "description": {
                    "type": "string",
                    "example": "Gitea"
                },
                "id": {
                    "type": "integer"
                },
                "local_directory": {
                    "type": "string",
                    "example": "/opt/docker/gitea"
                },
                "password_file_path": {
                    "type": "string",
                    "example": "/secrets/.resticpwd"
                },
                "post_commands": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.Command"
                    }
                },
                "pre_commands": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.Command"
                    }
                },
                "restic_remote": {
                    "type": "string",
                    "example": "rclone:pcloud:Backups/gitea"
                },
                "retention_policy": {
                    "enum": [
                        1,
                        2,
                        3,
                        4,
                        5,
                        6,
                        7
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/database.RetentionPolicy"
                        }
                    ],
                    "example": 1
                },
                "routine_check": {
                    "type": "integer",
                    "maximum": 100,
                    "minimum": 0
                },
                "runs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.Run"
                    }
                },
                "status": {
                    "$ref": "#/definitions/database.LogSeverity"
                }
            }
        },
        "database.JobStats": {
            "type": "object",
            "required": [
                "check_runs",
                "custom_runs",
                "error_logs",
                "info_logs",
                "prune_runs",
                "restic_runs",
                "total_logs",
                "total_runs",
                "warning_logs"
            ],
            "properties": {
                "check_runs": {
                    "type": "integer"
                },
                "custom_runs": {
                    "type": "integer"
                },
                "error_logs": {
                    "type": "integer"
                },
                "info_logs": {
                    "type": "integer"
                },
                "prune_runs": {
                    "type": "integer"
                },
                "restic_runs": {
                    "type": "integer"
                },
                "total_logs": {
                    "type": "integer"
                },
                "total_runs": {
                    "type": "integer"
                },
                "warning_logs": {
                    "type": "integer"
                }
            }
        },
        "database.Log": {
            "type": "object",
            "required": [
                "log_severity",
                "log_type",
                "message"
            ],
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "log_severity": {
                    "$ref": "#/definitions/database.LogSeverity"
                },
                "log_type": {
                    "$ref": "#/definitions/database.LogType"
                },
                "message": {
                    "type": "string"
                },
                "run_id": {
                    "type": "integer"
                }
            }
        },
        "database.LogSeverity": {
            "type": "integer",
            "enum": [
                0,
                1,
                2,
                3
            ],
            "x-enum-varnames": [
                "LogNone",
                "LogInfo",
                "LogWarning",
                "LogError"
            ]
        },
        "database.LogType": {
            "type": "integer",
            "enum": [
                1,
                2,
                3,
                4,
                5
            ],
            "x-enum-varnames": [
                "LogGeneral",
                "LogRestic",
                "LogCustom",
                "LogPrune",
                "LogCheck"
            ]
        },
        "database.RetentionPolicy": {
            "type": "integer",
            "enum": [
                1,
                2,
                3,
                4,
                5,
                6,
                7
            ],
            "x-enum-varnames": [
                "KeepAll",
                "KeepDailyLast2",
                "KeepDailyLast7",
                "KeepDailyLast31",
                "KeepMostRecent7Daily",
                "KeepMostRecent31Daily",
                "KeepDailyFor5Years"
            ]
        },
        "database.Run": {
            "type": "object",
            "properties": {
                "end_time": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "job_id": {
                    "type": "integer"
                },
                "logs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.Log"
                    }
                },
                "start_time": {
                    "type": "integer"
                },
                "status": {
                    "$ref": "#/definitions/database.LogSeverity"
                }
            }
        },
        "database.SystemLog": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "log_severity": {
                    "$ref": "#/definitions/database.LogSeverity"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "echo.HTTPError": {
            "type": "object",
            "properties": {
                "message": {}
            }
        },
        "system.Configuration": {
            "type": "object",
            "required": [
                "hostname",
                "rclone_config_file"
            ],
            "properties": {
                "hostname": {
                    "type": "string"
                },
                "rclone_config_file": {
                    "type": "string"
                }
            }
        },
        "system.Data": {
            "type": "object",
            "required": [
                "configuration",
                "job_stats",
                "versions"
            ],
            "properties": {
                "configuration": {
                    "$ref": "#/definitions/system.Configuration"
                },
                "job_stats": {
                    "$ref": "#/definitions/database.JobStats"
                },
                "versions": {
                    "$ref": "#/definitions/system.Versions"
                }
            }
        },
        "system.Versions": {
            "type": "object",
            "required": [
                "compose",
                "docker",
                "go",
                "gobackup",
                "rclone",
                "restic"
            ],
            "properties": {
                "compose": {
                    "type": "string"
                },
                "docker": {
                    "type": "string"
                },
                "go": {
                    "type": "string"
                },
                "gobackup": {
                    "type": "string"
                },
                "rclone": {
                    "type": "string"
                },
                "restic": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
