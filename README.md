# invertedIndexParallel
 CourseWork Parallel Computing
 
 How to build:
 1. clone repository   
 2. run 'docker build --pull --rm -f "Dockerfile" -t index:latest "."'
 3. docker run --name index --mount type=bind,source=$pwd\data,target=/app/data -p 8000:8000 -p 8080:8080 index:latest
 
