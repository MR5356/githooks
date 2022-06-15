name=$1
ssh_url=$2
commitId=$3

git clone $ssh_url
cd $name
docker build -t $name:$commitId -f release/docker/Dockerfile .