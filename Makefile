PROJECT_NAME = githubby

CURRENT_DIR = $(shell pwd)
GO_CMD = go
GO_BUILD = $(GO_CMD) build

BINARY_EXT = .out
BINARY_NAME = $(PROJECT_NAME)$(BINARY_EXT)
BINARY_PATH = $(CURRENT_DIR)/$(BINARY_NAME)

LOGFILE_NAME = log.txt
LOGFILE_PATH = $(CURRENT_DIR)/$(LOGFILE_NAME)

install:
	go get github.com/teejays/clog

build:
	$(GO_BUILD) -o $(BINARY_NAME) main.go

run-dev:
	$(BINARY_PATH) --dir $(CURRENT_DIR) --dev --wait-max 30 --max 3

## Production
PRD_WAIT_MAX = 7200

RUN_CMD = $(BINARY_PATH) --dir $(CURRENT_DIR) --wait-max $(PRD_WAIT_MAX)

CRON_TIME = 11 03 * * *
CRON_JOB = "$(CRON_TIME) $(RUN_CMD) > $(LOGFILE_PATH) &"
CRON_TEMP_FILE = _cron.txt

schedule: setup-logfile gen-cronfile schedule-cronfile
	
schedule-cronfile:
	crontab $(CRON_TEMP_FILE)

gen-cronfile:
	echo $(CRON_JOB) > $(CRON_TEMP_FILE)

setup-logfile:
	echo > $(LOGFILE_PATH) && sudo chmod 777 $(LOGFILE_PATH)