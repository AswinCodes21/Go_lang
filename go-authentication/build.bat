@echo off
setlocal enabledelayedexpansion

set APP_NAME=go-authentication
set DOCKER_COMPOSE=docker-compose
set GO=go

if "%1"=="" goto help

if "%1"=="build" (
    %GO% build -o %APP_NAME% ./cmd/main.go
    goto :eof
)

if "%1"=="run" (
    %GO% run ./cmd/main.go
    goto :eof
)

if "%1"=="docker-up" (
    %DOCKER_COMPOSE% up -d
    goto :eof
)

if "%1"=="docker-down" (
    %DOCKER_COMPOSE% down
    goto :eof
)

if "%1"=="docker-restart" (
    %DOCKER_COMPOSE% down
    %DOCKER_COMPOSE% up --build -d
    goto :eof
)

:help
echo Available commands:
echo   build         - Build the application
echo   run           - Run the application locally
echo   docker-up     - Start Docker containers
echo   docker-down   - Stop Docker containers
echo   docker-restart - Rebuild and restart Docker containers 