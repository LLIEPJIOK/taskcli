## About
Taskcli is a console application for creating and managing your tasks. You can add, update, remove, and view tasks. Additionally, there is a calendar that shows your activity for each day in the last six months:
- Days without color: 0 tasks created
- White days: 1-4 tasks created
- Blue days: 5-9 tasks created
- Green days: 10 or more tasks created.

## Getting started
To start application, you only need Docker running on your computer. 

To start working with application, follow these steps in console:
1. Clone the repository:
   ```bash
   git clone https://github.com/LLIEPJIOK/taskcli.git
   ```
2. Go to the project folder:
   ```bash
   cd taskcli
   ```
3. Start a docker container with the project:
   ```bash
   docker-compose up
   ```
4. Open another console where you will write commands. Each command has the following structure:
    ```bash
    docker exec -it taskcli-app-1 taskcli [COMMAND]
    ```
    You can view all available commands with:
    ```bash
    docker exec -it taskcli-app-1 taskcli info
    ```

To run test type:
   ```bash
   docker exec -it taskcli-app-1 go test -v ./[folder with tests]/*
   ```