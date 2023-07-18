 #!/bin/bash

 echo "###### Initializing application! ######"

docker compose down --volume
docker compose up --build
