# invertedIndexParallel
 CourseWork Parallel Computing
 
 How to build:
 1. clone repository   
 2. run 'docker build --pull --rm -f "Dockerfile" -t index:latest "."' in directory with repository
 3. run 'docker run --name index --mount type=bind,source=$SRC,target=/app/data -p 8000:8000 -p 8080:8080 index:latest -m $NUM_OF_MAPPERS -r $NUM_OF_REDUCERS', where $SRC is folder with files for index, $NUM_OF_MAPPERS is number of mapper threads, $NUM_OF_REDUCERS is number of reducer threads
 
