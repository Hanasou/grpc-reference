protoc -I src/ --go_out=plugins=grpc:./src src/greet/greetpb/greet.proto

protoc -I src/ --go_out=plugins=grpc:./src src/calculator/calculatorpb/calculator.proto

protoc -I src/ --go_out=plugins=grpc:./src src/blog/blogpb/blog.proto