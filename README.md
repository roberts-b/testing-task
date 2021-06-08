# Testing Task

Server runs on port 8888
Navigate to root folder of project where Dockerfile is located

## Steps to build and run

1. Build image: 
```
docker build -t testing-image .
```

2. Run image:
```
docker run -p 8888:8888 testing-image
```