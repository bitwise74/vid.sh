BACKEND_DIR=backend
FRONTEND_DIR=frontend
BUILD_DIR=build

.PHONY: all build clean deploy lint

all: build

clean:
	rm -rf $(BUILD_DIR)/*

build-frontend:
	cd $(FRONTEND_DIR) && npm install && npm run build
	cp -r $(FRONTEND_DIR)/build $(BUILD_DIR)/frontend

build-backend:
	cd $(BACKEND_DIR) && go build -o ../$(BUILD_DIR)/backend/server ./

build: clean build-frontend build-backend