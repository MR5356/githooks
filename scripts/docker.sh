name=$1
ssh_url=$2
commitId=$3

echo "$PWD/docker.sh clone project start"
cd /tmp
git clone $ssh_url
cd $name
echo "$PWD/docker.sh clone project finish"

echo "$PWD/docker.sh build image start"
BUILD_FILE=release/docker/build.sh
if [ -f "$BUILD_FILE" ]; then
  /bin/bash $BUILD_FILE $commitId
else
  docker build -t $name:$commitId -f release/docker/Dockerfile .
fi
echo "$PWD/docker.sh build image finish"

echo "$PWD/docker.sh clean images start"
docker rmi `docker images|grep none|awk '{print $3 }'|xargs`
echo "$PWD/docker.sh clean images finish"

cd -
rm -rf /tmp/$name

echo "$PWD/docker.sh finish"