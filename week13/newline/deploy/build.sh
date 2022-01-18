#!/bin/bash
# build micro services shell script

doBuild(){
  # build with parameter of tenant name, eg. bbcl
  echo "building binary ..."
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../src/srv/customer/customer-service_$1 ../src/srv/customer/*.go
  echo "customer service built successfully..."

  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../src/gateway/central/micro-gateway ../src/gateway/central/*.go
  echo "micro web built successfully..."

  if [ "" == "$2" ] || [ "dev" == "$2" ]; then
#    echo "build dev version"
    echo "use --config path:{/folder/configfile} syntax to build dev version"

    echo "changing dockerfile"
    cd ..
    pwd
    cp ./src/srv/customer/DockerfileDev ./src/srv/customer/Dockerfile
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
      sed -i "s/{tenant}/$1/g" ./src/srv/customer/Dockerfile
    elif [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' "s/{tenant}/$1/g" ./src/srv/customer/Dockerfile
    elif [[ "$OSTYPE" == "win32" ]]; then
      sed -i "s/{tenant}/$1/g" ./src/srv/customer/Dockerfile
    fi

    # change docker-compose file
    echo "change docker-compose file for $1"
    cp docker-compose.yaml docker-compose-$1.yaml
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
      sed -i "s/{tenant}/$1/g" docker-compose-$1.yaml
    elif [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' "s/{tenant}/$1/g" docker-compose-$1.yaml
    elif [[ "$OSTYPE" == "win32" ]]; then
      sed -i "s/{tenant}/$1/g" docker-compose-$1.yaml
    fi
#    docker-compose -f docker-compose-$1.yaml up --build -d

  else
    echo "build production version"
    # change dockerfile configuration with input parameter
    echo "changing Dockerfile ..."
    cd ..
    cp ./src/srv/customer/DockerfileTemplate ./src/srv/customer/Dockerfile

    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
      sed -i "s/{tenant}/$1/g" ./src/srv/customer/Dockerfile
    elif [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' "s/{tenant}/$1/g" ./src/srv/customer/Dockerfile
    elif [[ "$OSTYPE" == "win32" ]]; then
      sed -i "s/{tenant}/$1/g" ./src/srv/customer/Dockerfile
    fi

    # change docker-compose file
    echo "change docker-compose file for $1"
    cp docker-compose.yaml docker-compose-$1.yaml
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
      sed -i "s/{tenant}/$1/g" docker-compose-$1.yaml
    elif [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' "s/{tenant}/$1/g" docker-compose-$1.yaml
    elif [[ "$OSTYPE" == "win32" ]]; then
      sed -i "s/{tenant}/$1/g" docker-compose-$1.yaml
    fi
#    docker-compose -f docker-compose-$1.yaml up --build -d
  fi
}

# check input parameter whether null or not
if test -z "$1" ;then
  echo "you don't input tenant name, stop building ..."
else
  echo "you wanna build $1 $2"
  doBuild $1 $2
fi









