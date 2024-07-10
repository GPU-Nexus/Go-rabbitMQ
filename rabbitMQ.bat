@echo off
REM Define variables
set REPO_URL=https://github.com/GPU-Nexus/rabbitMQWithSpringBoot.git
set DIR_NAME=rabbitMQWithSpringBoot
set CMD_START_DIR=C:\Users\HP

echo got to start directory
cd %CMD_START_DIR%

REM Open a new Command Prompt window and execute a command
echo make connection to rabbitmq for 15672 port
start cmd /k "gcloud compute ssh rabbitmq-1-stats-vm-0 --ssh-flag="-L" --ssh-flag="15672:localhost:15672" --project=cse12-427001 --zone=europe-west9-b"


REM for other sssh connection
echo make connection to rabbitmq for 5672 port
start cmd /k "gcloud compute ssh rabbitmq-1-stats-vm-0 --ssh-flag="-L" --ssh-flag="5672:localhost:5672" --project=cse12-427001 --zone=europe-west9-b"





