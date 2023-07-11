 #!/bin/bash

 echo "###### Initializing application! ######"

 docker-compose down --volumes
 docker compose up --build
