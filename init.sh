 #!/bin/bash

 echo "###### Initializing MEMBER-MANAGER! ######"

docker compose down --volumes
docker compose up --build
