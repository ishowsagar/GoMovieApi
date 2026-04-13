SHELL := /bin/sh
.DEFAULT_GOAL := help

APP_DIR := .
PKG ?= npm

.PHONY: help install dev build preview lint lint-fix clean-dist clean-modules reset \
		status branch add commit push

help:
	@echo "Frontend commands:"
	@echo "  make install              Install dependencies"
	@echo "  make dev                  Start Vite dev server"
	@echo "  make build                Build production bundle"
	@echo "  make preview              Preview production build"
	@echo "  make lint                 Run ESLint"
	@echo "  make lint-fix             Run ESLint with --fix"
	@echo "  make clean-dist           Remove dist folder"
	@echo "  make clean-modules        Remove node_modules"
	@echo "  make reset                Fresh install (remove node_modules + install)"
	@echo "  make status               Git status"
	@echo "  make branch name=feature/x   Create and switch branch"
	@echo "  make commit msg='your msg'   Add all + commit"
	@echo "  make push                 Push current branch with upstream"

install:
	cd $(APP_DIR) && $(PKG) install

dev:
	cd $(APP_DIR) && $(PKG) run dev

build:
	cd $(APP_DIR) && $(PKG) run build

preview:
	cd $(APP_DIR) && $(PKG) run preview

lint:
	cd $(APP_DIR) && $(PKG) run lint

lint-fix:
	cd $(APP_DIR) && npx eslint . --fix

clean-dist:
	cd $(APP_DIR) && rm -rf dist

clean-modules:
	cd $(APP_DIR) && rm -rf node_modules

reset: clean-modules install

status:
	git status -sb

branch:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make branch name=feature/your-branch"; \
		exit 1; \
	fi
	git checkout -b "$(name)"

add:
	git add .

commit: add
	@if [ -z "$(msg)" ]; then \
		echo "Usage: make commit msg='your message'"; \
		exit 1; \
	fi
	git commit -m "$(msg)"

push:
	@branch=$$(git rev-parse --abbrev-ref HEAD); \
	git push -u origin $$branch
