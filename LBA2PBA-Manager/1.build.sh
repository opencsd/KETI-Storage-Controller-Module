cd ./cmake/build
make clean
make -j
cd ../..

registry="ketidevit2"
image_name="lba2pba-manager"
version="v1.0"

# make image
docker build -t $image_name:$version . && \

# add tag
docker tag $image_name:$version $registry/$image_name:$version && \

# login
docker login && \

# push image
docker push $registry/$image_name:$version 
