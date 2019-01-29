title AnonServer 
E:
cd E:\Programming\Go\anonymization
docker-compose down –v -t 3
docker rm -f -v anonymization_server
docker ps -a
docker-compose up -d --build 
docker logs -f  anonymization_server