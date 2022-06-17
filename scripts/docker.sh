name=$1
ssh_url=$2
commitId=$3

echo "scripts/docker.sh clone project start"
cd /tmp
git clone $ssh_url
cd $name
echo "scripts/docker.sh clone project finish"

echo "scripts/docker.sh build image start"
BUILD_FILE=release/docker/build.sh
Dockerfile=release/docker/Dockerfile

if [ -f "$Dockerfile" ]; then
  if [ -f "$BUILD_FILE" ]; then
    /bin/bash $BUILD_FILE $commitId
  else
    docker build -t $name:$commitId -f release/docker/Dockerfile .
  fi
  echo "scripts/docker.sh build image finish"
else
  echo "scripts/docker.sh build image skip"
fi

echo "scripts/docker.sh clean images start"
docker rmi `docker images|grep none|awk '{print $3 }'|xargs`
echo "scripts/docker.sh clean images finish"

cd -
rm -rf /tmp/$name

echo "scripts/docker.sh finish"
